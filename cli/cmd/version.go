package cmd

import (
	"os"

	"github.com/berquerant/pneutrinoutil/pkg/version"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print pneutrinoutil version",
	Run: func(_ *cobra.Command, _ []string) {
		version.Write(os.Stdout)
	},
}
