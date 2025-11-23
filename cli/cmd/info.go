package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/berquerant/pneutrinoutil/cli/info"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(infoCmd)
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Print system-wide information",
	RunE: func(cmd *cobra.Command, _ []string) error {
		neutrinoDir, _ := cmd.Flags().GetString("neutrinoDir")
		b := info.NewBuilder(neutrinoDir)
		r, err := b.Build(cmd.Context())
		if err != nil {
			return err
		}
		out, err := json.Marshal(r)
		if err != nil {
			return err
		}
		fmt.Printf("%s", out)
		return nil
	},
}
