package ctl

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/berquerant/execx"
	"github.com/berquerant/pneutrinoutil/cli/info"
	"github.com/berquerant/pneutrinoutil/pkg/pathx"
	"github.com/berquerant/structconfig"
	"github.com/spf13/pflag"
)

// NewDefaultConfig returns a new Config.
// Default values are from "default" struct tags.
func NewDefaultConfig() (*Config, error) {
	sc := structconfig.New[Config]()
	var c Config
	if err := sc.FromDefault(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

type Config struct {
	Description string `json:"desc" yaml:"desc" name:"desc" usage:"description of config"`
	// Project settings
	Score      string `json:"score" yaml:"score" name:"score" usage:"score file, required"`
	NumThreads int    `json:"thread" yaml:"thread" name:"thread" usage:"number of parallel in session" default:"4"`
	// NEUTRINO
	ModelDir        string `json:"model" yaml:"model" name:"model" usage:"singer" default:"MERROW"`
	SupportModelDir string `json:"supportModel" yaml:"supportModel" name:"supportModel" usage:"support singer"`
	Transpose       int    `json:"transpose" yaml:"transpose" name:"transpose" usage:"change the key and estimate" default:"0"`
	// Info
	NeutrinoVersion  string `json:"neutrinoVersion" yaml:"neutrinoVersion"`
	ModelData        any    `json:"modelData" yaml:"modelData"`
	SupportModelData any    `json:"supportModelData" yaml:"supportModelData"`
}

func (c Config) Basename() string { return pathx.Basename(c.Score) }

func (c Config) envMap() map[string]any {
	return map[string]any{
		"NumThreads":      c.NumThreads,
		"ModelDir":        c.ModelDir,
		"SupportModelDir": c.SupportModelDir,
		"BASENAME":        c.Basename(),
		"Transpose":       c.Transpose,
	}
}

func (c Config) Env() execx.Env {
	e := execx.NewEnv()
	for k, v := range c.envMap() {
		e.Set(k, fmt.Sprint(v))
	}
	return e
}

// ApplyFlagValues sets flag values to this.
func (c *Config) ApplyFlagValues(fs *pflag.FlagSet) error {
	sc := structconfig.New[Config]()
	return sc.FromFlags(c, fs)
}

// SetFlags sets command-line flags.
// Flag name is from "name" struct tag.
// Flag usage is from "usage" struct tag.
// Flag default value is from "default" struct tag.
func (c Config) SetFlags(fs *pflag.FlagSet) error {
	sc := structconfig.New[Config]()
	return sc.SetFlags(fs)
}

func (c *Config) SetInfo(ctx context.Context, neutrinoDir string) error {
	dir, err := filepath.Abs(neutrinoDir)
	if err != nil {
		return err
	}

	neutrinoVersion, err := info.GetNeutrinoVersion(ctx, dir)
	if err != nil {
		return err
	}
	c.NeutrinoVersion = neutrinoVersion

	model, err := info.ReadModelInfo(filepath.Join(dir, "model", c.ModelDir))
	if err != nil {
		return err
	}
	c.ModelData = model

	if x := c.SupportModelDir; x != "" {
		model, err := info.ReadModelInfo(filepath.Join(dir, "model", x))
		if err != nil {
			return err
		}
		c.SupportModelData = model
	}
	return nil
}
