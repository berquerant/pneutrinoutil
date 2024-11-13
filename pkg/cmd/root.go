package cmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/berquerant/pneutrinoutil/pkg/ctl"
	"github.com/berquerant/pneutrinoutil/pkg/logx"
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

		play, _ := cmd.Flags().GetBool("play")
		tasks := newTasks(dir, c, now, play)
		if dry, _ := cmd.Flags().GetBool("dry"); dry {
			ss := make([]string, len(tasks))
			for i, t := range tasks {
				ss[i] = t.String()
			}
			fmt.Println(strings.Join(ss, "\n"))
			return nil
		}

		type stat struct {
			title   string
			elapsed time.Duration
		}
		stats := []stat{}
		for _, t := range tasks {
			if err := t.Run(cmd.Context()); err != nil {
				return err
			}
			stats = append(stats, stat{
				title:   t.Title(),
				elapsed: t.Elapsed(),
			})
		}
		for i, s := range stats {
			slog.Info(
				"Stat",
				slog.String("title", s.title),
				slog.Int("index", i),
				slog.Duration("elapsed", s.elapsed),
				slog.Float64("elapsedSeconds", s.elapsed.Seconds()),
			)
		}
		return nil
	},
}
