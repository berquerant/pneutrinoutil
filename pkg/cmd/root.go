package cmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/berquerant/pneutrinoutil/pkg/ctl"
	"github.com/berquerant/pneutrinoutil/pkg/logx"
	"github.com/berquerant/pneutrinoutil/pkg/task"
	"github.com/spf13/cobra"
)

func Main(ctx context.Context) {
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		slog.Error("main", logx.Err(err))
	}
}

var (
	ErrArgument = errors.New("Argument")
)

func init() {
	rootCmd.PersistentFlags().Bool("debug", false, "enable debug")
	rootCmd.PersistentFlags().StringP("workDir", "w", ".", "working directory")
	rootCmd.PersistentFlags().StringP("neutrinoDir", "n", "./dist/NEUTRINO", "NEUTRINO directory")
	rootCmd.Flags().Bool("dry", false, "dryrun")
	rootCmd.Flags().Bool("play", false, "play generated wav after running")

	var c ctl.Config
	if err := c.SetFlags(rootCmd.Flags()); err != nil {
		panic(err)
	}

	rootCmd.Flags().StringSliceP("include", "i", nil, "include task names")
	rootCmd.Flags().StringSliceP("exclude", "e", nil, "exclude task names")
}

var rootCmd = &cobra.Command{
	Use:   "pneutrinoutil [CONFIG_YML]",
	Short: `Generate .wav from .musicxml using NEUTRINO`,
	PersistentPreRun: func(cmd *cobra.Command, _ []string) {
		debugEnabled, _ := cmd.Flags().GetBool("debug")
		logLevel := slog.LevelInfo
		if debugEnabled {
			logLevel = slog.LevelDebug
		}
		logx.Setup(os.Stderr, logLevel)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		now := time.Now()
		c, err := newConfig(cmd, args)
		if err != nil {
			return err
		}
		dir := newDir(cmd)

		var (
			play, _    = cmd.Flags().GetBool("play")
			include, _ = cmd.Flags().GetStringSlice("include")
			exclude, _ = cmd.Flags().GetStringSlice("exclude")
		)

		runner := newTaskRunner(dir, c, now, play)

		if dry, _ := cmd.Flags().GetBool("dry"); dry {
			fmt.Println(runner.String())
			return nil
		}

		err = runner.Run(cmd.Context(), task.WithInclude(include), task.WithExclude(exclude))
		for i, s := range runner.Stats() {
			slog.Info(
				"Stat",
				slog.String("title", s.Title()),
				slog.Int("index", i),
				slog.Duration("elapsed", s.Elapsed()),
				slog.Float64("elapsedSeconds", s.Elapsed().Seconds()),
			)
		}
		return err
	},
}
