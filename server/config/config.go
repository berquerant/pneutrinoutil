package config

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/berquerant/pneutrinoutil/pkg/alog"
	"github.com/berquerant/pneutrinoutil/pkg/pathx"
	"github.com/berquerant/structconfig"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hibiken/asynq"
	"github.com/labstack/gommon/log"
	"github.com/spf13/pflag"
)

type Config struct {
	Host                        string `name:"host" usage:"server host"`
	Port                        uint   `name:"port" short:"p" default:"9101" usage:"server port"`
	ShutdownPeriodSeconds       int    `name:"shutdownPeriodSeconds" default:"10" usage:"duration the server needs to shut down gracefully"`
	ProcessTimeoutSeconds       int    `name:"processTimeoutSeconds" default:"1200" usage:"duration pneutrinoutil timeout"`
	AccessLogFile               string `name:"accessLogFile" default:"-" usage:"access log file; - means stderr"`
	AccessLogWriter             io.Writer
	Debug                       bool   `name:"debug" usage:"enable debug logs"`
	RedisDSN                    string `name:"redisDSN" usage:"format: redis://HOST:PORT/DB"`
	MysqlDSN                    string `name:"mysqlDSN" usage:"format: USER:PASS@tcp(HOST:PORT)/DB"`
	MysqlConnMaxLifetimeSeconds int    `name:"mysqlConnMaxLifetimeSeconds" default:"300" usage:"max amount of time a connection may be reused"`
	MysqlMaxIdleConns           int    `name:"mysqlMaxIdleConns" default:"3" usage:"maximum number of connections in the idle connection pool"`
	MysqlMaxOpenConns           int    `name:"mysqlMaxOpenConns" default:"3" usage:"maximum number of open connections to the database"`
	// StorageAddr           string `name:"storageAddr" usage:"object storage addr"`
	StorageDir    string `name:"storageDir" usage:"local storage directory; $HOME/.pneutrinoutil-worker/storage or .pneutrinoutil-worker/storage if no $HOME"`
	StorageBucket string `name:"storageBucket" default:"pneutrinoutil-worker" usage:"storage bucket"`
	StoragePath   string `name:"storagePath" usage:"storage base path"`
}

func (c Config) Addr() string { return fmt.Sprintf("%s:%d", c.Host, c.Port) }

func (c Config) NewAsynqClient() (*asynq.Client, error) {
	opt, err := asynq.ParseRedisURI(c.RedisDSN)
	if err != nil {
		return nil, err
	}
	return asynq.NewClient(opt), nil
}

func (c Config) NewSQL(ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("mysql", c.MysqlDSN)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(time.Duration(c.MysqlConnMaxLifetimeSeconds) * time.Second)
	db.SetMaxIdleConns(c.MysqlMaxIdleConns)
	db.SetMaxOpenConns(c.MysqlMaxOpenConns)

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

func (c Config) ShutdownPeriod() time.Duration {
	return time.Duration(c.ShutdownPeriodSeconds) * time.Second
}
func (c Config) ProcessTimeout() time.Duration {
	return time.Duration(c.ProcessTimeoutSeconds) * time.Second
}

func (c Config) EchoLogLevel() log.Lvl {
	if c.Debug {
		return log.DEBUG
	}
	return log.INFO
}
func (c Config) SLogLevel() slog.Leveler {
	if c.Debug {
		return slog.LevelDebug
	}
	return slog.LevelInfo
}

func New(fs *pflag.FlagSet) (*Config, error) {
	c, err := structconfig.NewConfigWithMerge(
		structconfig.New[Config](),
		structconfig.NewMerger[Config](),
		fs,
	)
	if err != nil {
		return nil, err
	}
	if c.StorageDir == "" {
		c.StorageDir = filepath.Join(pathx.UserHomeDirOr("."), ".pneutrinoutil-worker/storage")
	}
	return c, nil
}

func (c *Config) Init() error {
	if err := c.Validate(); err != nil {
		return err
	}
	// open access log
	{
		w, err := c.prepareAccessLogWriter()
		if err != nil {
			return err
		}
		c.AccessLogWriter = w
	}
	// ensure directories
	for _, dir := range []string{
		c.StorageDir,
	} {
		alog.L().Debug("ensure dir", slog.String("path", dir))
		if err := pathx.EnsureDir(dir); err != nil {
			return err
		}
	}
	return nil
}

const (
	AccessLogStderr = "-"
)

func (c Config) prepareAccessLogWriter() (io.Writer, error) {
	if c.AccessLogFile == AccessLogStderr {
		return os.Stderr, nil
	}

	f, err := os.OpenFile(c.AccessLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("%w: open access log file", err)
	}
	return f, nil
}

func (c Config) Validate() error {
	return nil
}

func (c *Config) Close() error {
	var err error
	if x, ok := c.AccessLogWriter.(io.Closer); ok {
		err = errors.Join(err, x.Close())
	}
	return err
}
