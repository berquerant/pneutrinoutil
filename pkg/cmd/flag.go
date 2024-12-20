package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/berquerant/pneutrinoutil/pkg/ctl"
	"github.com/berquerant/pneutrinoutil/pkg/task"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func newDir(cmd *cobra.Command) *task.Dir {
	workDir, _ := cmd.Flags().GetString("workDir")
	neutrinoDir, _ := cmd.Flags().GetString("neutrinoDir")
	pwd := os.Getenv("PWD")
	return task.NewDir(workDir, neutrinoDir, pwd)
}

func newConfig(cmd *cobra.Command, args []string) (*ctl.Config, error) {
	if len(args) > 1 {
		return nil, fmt.Errorf("%w: require 0 or 1 argument for config file", ErrArgument)
	}

	c, err := ctl.NewDefaultConfig()
	if err != nil {
		return nil, err
	}
	// override by file
	if len(args) == 1 {
		if err := func() error {
			f, err := os.Open(args[0])
			if err != nil {
				return err
			}
			defer f.Close()
			b, err := io.ReadAll(f)
			if err != nil {
				return err
			}
			if err := yaml.Unmarshal(b, c); err != nil {
				return err
			}
			return nil
		}(); err != nil {
			return nil, err
		}
	}
	// override by flags
	if err := c.ApplyFlagValues(cmd.Flags()); err != nil {
		return nil, err
	}
	return c, nil
}
