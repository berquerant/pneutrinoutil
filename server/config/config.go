package config

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/berquerant/pneutrinoutil/pkg/pathx"
	"github.com/berquerant/structconfig"
	"github.com/spf13/pflag"
)

type Config struct {
	Host                  string    `name:"host" usage:"server host"`
	Port                  uint      `name:"port" short:"p" default:"9101" usage:"server port"`
	Pneutrinoutil         string    `name:"pneutrinoutil" short:"x" default:"pneutrinoutil" usage:"pneutrinoutil executable"`
	WorkDir               string    `name:"workDir" short:"w" default:"." usage:"working directory"`
	NeutrinoDir           string    `name:"neutrinoDir" short:"n" default:"./dist/NEUTRINO" usage:"NEUTRINO directory"`
	ShutdownPeriodSeconds int       `name:"shutdownPeriodSeconds" default:"10" usage:"duration the server needs to shut down gracefully"`
	ProcessTimeoutSeconds int       `name:"processTimeoutSeconds" default:"1200" usage:"duration pneutrinoutil timeout"`
	AccessLogFile         string    `name:"accessLogFile" default:"-" usage:"access log file; - means stdout"`
	AccessLogWriter       io.Writer `name:"-"`
	Shell                 string    `name:"shell" short:"s" default:"bash" usage:"shell command to execute"`
	Concurrency           int       `name:"concurrency" short:"c" default:"1" usage:"pneutrinoutil process concurrency"`
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

func New(fs *pflag.FlagSet) (*Config, error) {
	c, err := structconfig.NewConfigWithMerge(
		structconfig.New[Config](),
		structconfig.NewMerger[Config](),
		fs,
	)
	if err != nil {
		return nil, err
	}
	// open access log
	w, err := c.prepareAccessLogWriter()
	if err != nil {
		return nil, err
	}
	c.AccessLogWriter = w
	// ensure directories
	for _, dir := range []string{
		c.WorkDir,
		c.TmpDir(),
		c.UploadDir(),
		c.PneutrinoutilWorkDir(),
		c.LogDir(),
	} {
		if err := pathx.EnsureDir(dir); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (c Config) prepareAccessLogWriter() (io.Writer, error) {
	if c.AccessLogFile == "-" {
		return os.Stdout, nil
	}

	f, err := os.OpenFile(c.AccessLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("%w: open access log file", err)
	}
	return f, nil
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
	}
}

func (c Config) Validate() error {
	if _, err := exec.LookPath(c.Pneutrinoutil); err != nil {
		return fmt.Errorf("%w: look path %s", err, c.Pneutrinoutil)
	}
	return nil
}

func (c *Config) Close() error {
	if x, ok := c.AccessLogWriter.(*os.File); ok {
		return x.Close()
	}
	return nil
}
