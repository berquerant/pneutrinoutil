package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	cli "github.com/berquerant/pneutrinoutil/cli/cmd"
	"github.com/berquerant/pneutrinoutil/pkg/logx"
	"github.com/berquerant/pneutrinoutil/pkg/pathx"
	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func init() {
	cli.InitFlags(rootCmd)
	rootCmd.Flags().Bool("fail", false, "If true, exit with 1")
}

var rootCmd = &cobra.Command{
	Use:   "mockcli",
	Short: "mock pneutrinoutil CLI for testing",
	RunE: func(cmd *cobra.Command, args []string) error {
		now := time.Now()
		logger := logx.NewLTSVLogger(os.Stdout, slog.LevelDebug)
		if v, _ := cmd.Flags().GetBool("fail"); v {
			logger.Error("Should fail due to --fail")
			os.Exit(1)
		}
		// collect flag values
		flags := map[string]string{}
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			flags[f.Name] = f.Value.String()
		})
		logger.Info("flags", logx.JSON("values", flags))

		// parse config
		c, err := cli.NewConfig(cmd, args)
		if err != nil {
			return err
		}
		logger.Info("config", logx.JSON("values", c))

		var (
			open = func(name, path string) (*os.File, error) {
				logger.Info("open "+name, slog.String("path", path))
				f, err := os.Open(path)
				if err != nil {
					return nil, fmt.Errorf("%w: open %s, %s", err, name, path)
				}
				return f, nil
			}
			create = func(name, path string) (*os.File, error) {
				logger.Info("create "+name, slog.String("path", path))
				f, err := os.Create(path)
				if err != nil {
					return nil, fmt.Errorf("%w: create %s, %s", err, name, path)
				}
				return f, nil
			}
		)

		// open score file
		score, err := open("score", c.Score)
		if err != nil {
			return err
		}
		defer func() { _ = score.Close() }()

		// create result directory
		workDir, _ := cmd.Flags().GetString("workDir")
		re := pathx.NewResultElement(c.Basename(), now, now.Unix(), 9101)
		resultDir := filepath.Join(workDir, "result", re.String())
		logger.Info("mkdir", slog.String("path", resultDir))
		if err := pathx.EnsureDir(resultDir); err != nil {
			return err
		}

		// create config file in result directory
		configFileName := filepath.Join(resultDir, "config.yml")
		configFile, err := create("config", configFileName)
		if err != nil {
			return err
		}
		defer func() { _ = configFile.Close() }()
		configYaml, err := yaml.Marshal(c)
		if err != nil {
			return err
		}
		if _, err := configFile.Write(configYaml); err != nil {
			return err
		}

		// copy score to result directory
		resultScorePath := filepath.Join(resultDir, c.Basename()+".musicxml")
		resultScore, err := create("musicxml", resultScorePath)
		if err != nil {
			return err
		}
		defer func() { _ = resultScore.Close() }()
		if _, err := io.Copy(resultScore, score); err != nil {
			return err
		}

		// create empty wav files
		wavPath := filepath.Join(resultDir, c.Basename()+".wav")
		wav, err := create("wav", wavPath)
		if err != nil {
			return err
		}
		_ = wav.Close()
		worldWavPath := filepath.Join(resultDir, c.Basename()+"_world.wav")
		worldWav, err := create("world_wav", worldWavPath)
		if err != nil {
			return err
		}
		_ = worldWav.Close()

		return nil
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
