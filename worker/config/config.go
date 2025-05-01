package config

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/berquerant/pneutrinoutil/pkg/alog"
	"github.com/berquerant/pneutrinoutil/pkg/infra"
	"github.com/berquerant/pneutrinoutil/pkg/pathx"
	"github.com/berquerant/structconfig"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hibiken/asynq"
	"github.com/spf13/pflag"
)

func New(fs *pflag.FlagSet) (*Config, error) {
	c, err := structconfig.NewConfigWithMerge(
		structconfig.New[Config](),
		structconfig.NewMerger[Config](),
		fs,
	)
	if err != nil {
		return nil, err
	}
	if c.WorkDir == "" {
		c.WorkDir = filepath.Join(pathx.UserHomeDirOr("."), ".pneutrinoutil-worker/workspace")
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
	for _, dir := range []string{
		c.StorageDir,
		c.WorkDir,
	} {
		alog.L().Debug("ensure dir", slog.String("path", dir))
		if err := pathx.EnsureDir(dir); err != nil {
			return err
		}
	}
	return nil
}

type Config struct {
	RedisDSN                    string `name:"redisDSN" usage:"format: redis://HOST:PORT/DB"`
	MysqlDSN                    string `name:"mysqlDSN" usage:"format: USER:PASS@tcp(HOST:PORT)/DB"`
	MysqlConnMaxLifetimeSeconds int    `name:"mysqlConnMaxLifetimeSeconds" default:"300" usage:"max amount of time a connection may be reused"`
	MysqlMaxIdleConns           int    `name:"mysqlMaxIdleConns" default:"3" usage:"maximum number of connections in the idle connection pool"`
	MysqlMaxOpenConns           int    `name:"mysqlMaxOpenConns" default:"3" usage:"maximum number of open connections to the database"`
	// StorageAddr           string `name:"storageAddr" usage:"object storage addr"`
	StorageDir            string `name:"storageDir" usage:"local storage directory; $HOME/.pneutrinoutil-worker/storage or .pneutrinoutil-worker/storage if no $HOME"`
	StorageBucket         string `name:"storageBucket" default:"pneutrinoutil-worker" usage:"storage bucket"`
	StoragePath           string `name:"storagePath" usage:"storage base path"`
	Pneutrinoutil         string `name:"pneutrinoutil" short:"x" default:"./dist/pneutrinoutil" usage:"pneutrinoutil executable"`
	WorkDir               string `name:"workDir" short:"w" usage:"working directory; $HOME/.pneutrinoutil-worker/workspace or .pneutrinoutil-worker/workspace if no $HOME"`
	NeutrinoDir           string `name:"neutrinoDir" short:"n" default:"./dist/NEUTRINO" usage:"NEUTRINO directory"`
	ShutdownPeriodSeconds int    `name:"shutdownPeriodSeconds" default:"10" usage:"duration the server needs to shut down gracefully"`
	Shell                 string `name:"shell" short:"s" default:"bash" usage:"shell command to execute"`
	Concurrency           int    `name:"concurrency" short:"c" default:"1" usage:"pneutrinoutil process concurrency"`
	Debug                 bool   `name:"debug" usage:"enable debug logs"`
	Webhook               string `name:"webhook" usage:"webhook endpoint to notify task completion"`
	WebhookTimeoutSeconds int    `name:"webhookTimeoutSeconds" default:"10" usage:"duration webhook timeout"`
}

func (c Config) NewWebhook() *infra.Webhook {
	if c.Webhook == "" {
		return nil
	}
	return infra.NewWebhook(c.Webhook, time.Duration(c.WebhookTimeoutSeconds)*time.Second)
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

func (c Config) SLogLevel() slog.Leveler {
	if c.Debug {
		return slog.LevelDebug
	}
	return slog.LevelInfo
}

func (c Config) AsynqLogLevel() asynq.LogLevel {
	// if c.Debug {
	// 	return asynq.DebugLevel
	// }
	return asynq.InfoLevel
}

func (c Config) Validate() error {
	if _, err := exec.LookPath(c.Pneutrinoutil); err != nil {
		return fmt.Errorf("%w: look path %s", err, c.Pneutrinoutil)
	}
	return nil
}

func (Config) Close() error { return nil }
