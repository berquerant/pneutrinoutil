package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/berquerant/pneutrinoutil/cli/ctl"
	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(skeletonCmd)
	skeletonCmd.Flags().Bool("json", false, "as json")
}

var skeletonCmd = &cobra.Command{
	Use:   "skeleton",
	Short: "Dump default config.yml",
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := ctl.NewDefaultConfig()
		if err != nil {
			return err
		}

		var b []byte
		if x, _ := cmd.Flags().GetBool("json"); x {
			b, err = json.Marshal(c)
		} else {
			b, err = yaml.Marshal(c)
		}
		if err != nil {
			return err
		}
		fmt.Printf("%s", b)
		return nil
	},
}
