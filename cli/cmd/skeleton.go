package cmd

import (
	"fmt"

	"github.com/berquerant/pneutrinoutil/cli/ctl"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func init() {
	rootCmd.AddCommand(skeletonCmd)
}

var skeletonCmd = &cobra.Command{
	Use:   "skeleton",
	Short: "Dump default config.yml",
	RunE: func(_ *cobra.Command, _ []string) error {
		c, err := ctl.NewDefaultConfig()
		if err != nil {
			return err
		}
		b, err := yaml.Marshal(c)
		if err != nil {
			return err
		}
		fmt.Printf("%s", b)
		return nil
	},
}
