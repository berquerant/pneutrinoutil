package ctl

import (
	"github.com/berquerant/structconfig"
	"github.com/spf13/pflag"
)

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
