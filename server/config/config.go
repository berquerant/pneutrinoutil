package config

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/berquerant/pneutrinoutil/pkg/pathx"
	"github.com/berquerant/pneutrinoutil/server/alog"
	"github.com/berquerant/structconfig"
	"github.com/labstack/gommon/log"
	"github.com/spf13/pflag"
)

type Config struct {
	Host                       string `name:"host" usage:"server host"`
	Port                       uint   `name:"port" short:"p" default:"9101" usage:"server port"`
	Pneutrinoutil              string `name:"pneutrinoutil" short:"x" default:"./dist/pneutrinoutil" usage:"pneutrinoutil executable"`
	WorkDir                    string `name:"workDir" short:"w" usage:"working directory; $HOME/.pneutrinoutil-server or .pneutrinoutil-server if no $HOME"`
	NeutrinoDir                string `name:"neutrinoDir" short:"n" default:"./dist/NEUTRINO" usage:"NEUTRINO directory"`
	ShutdownPeriodSeconds      int    `name:"shutdownPeriodSeconds" default:"10" usage:"duration the server needs to shut down gracefully"`
	ProcessTimeoutSeconds      int    `name:"processTimeoutSeconds" default:"1200" usage:"duration pneutrinoutil timeout"`
	AccessLogFile              string `name:"accessLogFile" default:"-" usage:"access log file; - means stderr"`
	AccessLogWriter            io.Writer
	Shell                      string `name:"shell" short:"s" default:"bash" usage:"shell command to execute"`
	Concurrency                int    `name:"concurrency" short:"c" default:"1" usage:"pneutrinoutil process concurrency"`
	Debug                      bool   `name:"debug" usage:"enable debug logs"`
	NotificationCommand        string `name:"notify" usage:"command to notify process completed; the server invokes the command like: COMMAND REQUEST_ID STATUS; STATUS=0 means success; REQUEST_ID and STATUS are available as environment variables"`
	NotificationTimeoutSeconds int    `name:"notifyTimeoutSeconds" default:"10" usage:"duration notification timeout"`
	NotificationLogFile        string `name:"notifyLogFile" default:"" usage:"notification log file; - means stderr; empty means /dev/null"`
	NotificationLogWriter      io.Writer
}

func (c Config) Addr() string                 { return fmt.Sprintf("%s:%d", c.Host, c.Port) }
func (c Config) TmpDir() string               { return filepath.Join(c.WorkDir, "tmp") }
func (c Config) UploadDir() string            { return filepath.Join(c.TmpDir(), "upload") }
func (c Config) PneutrinoutilWorkDir() string { return filepath.Join(c.WorkDir, "pneutrinoutil") }
func (c Config) LogDir() string               { return filepath.Join(c.TmpDir(), "log") }

func (c Config) ShutdownPeriod() time.Duration {
	return time.Duration(c.ShutdownPeriodSeconds) * time.Second
}
func (c Config) ProcessTimeout() time.Duration {
	return time.Duration(c.ProcessTimeoutSeconds) * time.Second
}
func (c Config) NotificationTimeout() time.Duration {
	return time.Duration(c.NotificationTimeoutSeconds) * time.Second
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
	if c.WorkDir == "" {
		c.WorkDir = filepath.Join(pathx.UserHomeDirOr("."), ".pneutrinoutil-server")
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
	// open notification log
	{
		w, err := c.prepareNotificationLogWriter()
		if err != nil {
			return err
		}
		c.NotificationLogWriter = w
	}
	// ensure directories
	for _, dir := range []string{
		c.WorkDir,
		c.TmpDir(),
		c.UploadDir(),
		c.PneutrinoutilWorkDir(),
		c.LogDir(),
	} {
		alog.L().Debug("ensure dir", slog.String("path", dir))
		if err := pathx.EnsureDir(dir); err != nil {
			return err
		}
	}
	return nil
}

const (
	AccessLogStderr       = "-"
	NotificationLogStderr = "-"
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

func (c Config) prepareNotificationLogWriter() (io.Writer, error) {
	switch c.NotificationLogFile {
	case "":
		return io.Discard, nil
	case NotificationLogStderr:
		return os.Stderr, nil
	default:
		f, err := os.OpenFile(c.NotificationLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("%w: open notification log file", err)
		}
		return f, nil
	}
}

func (c Config) LogAttr() []any {
	attrs := []any{}
	for k, v := range c.intoMap() {
		attrs = append(attrs, slog.Any(k, v))
	}
	return attrs
}

func (c Config) intoMap() map[string]any {
	return map[string]any{
		"host":                  c.Host,
		"port":                  c.Port,
		"pneutrinoutil":         c.Pneutrinoutil,
		"workDir":               c.WorkDir,
		"neutrinoDir":           c.NeutrinoDir,
		"shutdownPeriodSeconds": c.ShutdownPeriodSeconds,
		"processTimeoutSeconds": c.ProcessTimeoutSeconds,
		"accessLogFile":         c.AccessLogFile,
		"shell":                 c.Shell,
		"concurrency":           c.Concurrency,
		"debug":                 c.Debug,
		"notify":                c.NotificationCommand,
		"notifyTimeoutSeconds":  c.NotificationTimeoutSeconds,
		"notifyLogFile":         c.NotificationLogFile,
	}
}

func (c Config) Validate() error {
	if _, err := exec.LookPath(c.Pneutrinoutil); err != nil {
		return fmt.Errorf("%w: look path %s", err, c.Pneutrinoutil)
	}
	if x := c.NotificationCommand; x != "" {
		if _, err := exec.LookPath(x); err != nil {
			return fmt.Errorf("%w: look path %s", err, x)
		}
	}
	return nil
}

func (c *Config) Close() error {
	var err error
	if x, ok := c.AccessLogWriter.(io.Closer); ok {
		err = errors.Join(err, x.Close())
	}
	if x, ok := c.NotificationLogWriter.(io.Closer); ok {
		err = errors.Join(err, x.Close())
	}
	return err
}
