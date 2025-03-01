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

func InitFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().Bool("debug", false, "enable debug")
	cmd.PersistentFlags().StringP("workDir", "w", ".", "working directory")
	cmd.PersistentFlags().StringP("neutrinoDir", "n", "./dist/NEUTRINO", "NEUTRINO directory")
	cmd.Flags().Bool("dry", false, "dryrun")
	cmd.Flags().String("play", "", "play command generated wav after running, wav file will be passed to 1st argument")
	cmd.Flags().String("hook", "", "command to be executed after running, result dir will be passed to 1st argument")
	cmd.Flags().Bool("list-tasks", false, "list task names")
	cmd.Flags().StringSlice("env", nil, "names of additional environment variables to allow reading; all allows everythings")
	cmd.Flags().StringP("shell", "s", "bash", "shell command to execute")

	var c ctl.Config
	if err := c.SetFlags(cmd.Flags()); err != nil {
		panic(err)
	}

	cmd.Flags().StringSliceP("include", "i", nil, "include task names")
	cmd.Flags().StringSliceP("exclude", "e", nil, "exclude task names")
}

func init() {
	InitFlags(rootCmd)
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
		c, err := NewConfig(cmd, args)
		if err != nil {
			return err
		}

		var (
			dir        = NewDir(cmd, now)
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
