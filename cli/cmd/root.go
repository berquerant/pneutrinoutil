package cmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/berquerant/execx"
	"github.com/berquerant/pneutrinoutil/cli/ctl"
	"github.com/berquerant/pneutrinoutil/cli/task"
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
	rootCmd.Flags().String("play", "", "play command generated wav after running, wav file will be passed to 1st argument")
	rootCmd.Flags().String("hook", "", "command to be executed after running, result dir will be passed to 1st argument")
	rootCmd.Flags().Bool("list-tasks", false, "list task names")
	rootCmd.Flags().StringSlice("env", nil, "names of additional environment variables to allow reading; all allows everythings")
	rootCmd.Flags().StringP("shell", "s", "bash", "shell command to execute")

	var c ctl.Config
	if err := c.SetFlags(rootCmd.Flags()); err != nil {
		panic(err)
	}

	rootCmd.Flags().StringSliceP("include", "i", nil, "include task names")
	rootCmd.Flags().StringSliceP("exclude", "e", nil, "exclude task names")
}

var rootCmd = &cobra.Command{
	Use:   "pneutrinoutil [CONFIG_YML|CONFIG_JSON]",
	Short: `Generate .wav from .musicxml using NEUTRINO`,
	Long: `Generate .wav from .musicxml using NEUTRINO

e.g.
pneutrinoutil --neutrinoDir /path/to/NEUTRINO --workDir /path/to/install-result --score /path/to/some.musicxml`,
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

		var (
			dir        = newDir(cmd, now)
			play, _    = cmd.Flags().GetString("play")
			hook, _    = cmd.Flags().GetString("hook")
			include, _ = cmd.Flags().GetStringSlice("include")
			exclude, _ = cmd.Flags().GetStringSlice("exclude")
		)

		tasks := task.NewGenerator(dir, c, play, hook).ExecutableTasks()
		taskNames := make([]string, len(tasks.Tasks))
		for i, t := range tasks.Tasks {
			taskNames[i] = t.Name
		}

		if list, _ := cmd.Flags().GetBool("list-tasks"); list {
			fmt.Println(strings.Join(taskNames, "\n"))
			return nil
		}

		tasks.Entrypoint = prepareTaskEntrypoint(taskNames, include, exclude)

		environWhiteList, _ := cmd.Flags().GetStringSlice("env")
		tasks.Env.Merge(prepareAdditionalEnviron(environWhiteList))

		if dry, _ := cmd.Flags().GetBool("dry"); dry {
			slog.Info("generated script should be called on the dir", "dir", dir.NeutrinoDir())
			fmt.Println(tasks.String())
			return nil
		}

		shell, _ := cmd.Flags().GetString("shell")
		return tasks.IntoScript(shell).Runner(func(cmd *execx.Cmd) error {
			cmd.Dir = dir.NeutrinoDir()
			slog.Info("exec", "dir", cmd.Dir, "args", cmd.Args)
			return cmd.Exec()
		})
	},
}
