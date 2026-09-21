package server

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/berquerant/execx"
	"github.com/berquerant/pneutrinoutil/cli/ctl"
	"github.com/berquerant/pneutrinoutil/cli/info"
	"github.com/berquerant/pneutrinoutil/cli/task"
	"github.com/berquerant/pneutrinoutil/pkg/pathx"
)

type StandaloneRunner struct {
	cfg *Config
}

func NewStandaloneRunner(cfg *Config) *StandaloneRunner {
	return &StandaloneRunner{cfg: cfg}
}

func (r *StandaloneRunner) GetInfo(ctx context.Context, opt *ExecutionOptions) (*info.Info, error) {
	_, neutrinoDir, _, _ := r.cfg.ResolveEffectiveConfig(opt)
	b := info.NewBuilder(neutrinoDir)
	return b.Build(ctx)
}

func (r *StandaloneRunner) ListModels(ctx context.Context, opt *ExecutionOptions) ([]info.Model, error) {
	inf, err := r.GetInfo(ctx, opt)
	if err != nil {
		return nil, err
	}
	return inf.Neutrino.Models, nil
}

func (r *StandaloneRunner) Synthesize(ctx context.Context, params *SynthesizeParams) (*SynthesizeResult, error) {
	startTime := time.Now()
	_, neutrinoDir, workDir, _ := r.cfg.ResolveEffectiveConfig(params.Options)

	scorePath := params.ScorePath
	if scorePath == "" {
		if params.ScoreContent == "" {
			return nil, fmt.Errorf("either scorePath or scoreContent must be provided")
		}
		tmpScoreDir := filepath.Join(workDir, "scores")
		if err := pathx.EnsureDir(tmpScoreDir); err != nil {
			return nil, fmt.Errorf("failed to create score dir: %w", err)
		}
		tmpScoreFile := filepath.Join(tmpScoreDir, fmt.Sprintf("score_%d.musicxml", time.Now().UnixNano()))
		if err := os.WriteFile(tmpScoreFile, []byte(params.ScoreContent), 0644); err != nil {
			return nil, fmt.Errorf("failed to write score file: %w", err)
		}
		scorePath = tmpScoreFile
	}

	absScorePath, err := filepath.Abs(scorePath)
	if err != nil {
		return nil, fmt.Errorf("invalid scorePath: %w", err)
	}

	c, err := ctl.NewDefaultConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create default config: %w", err)
	}

	c.Score = absScorePath
	if params.Model != "" {
		c.ModelDir = params.Model
	}
	if params.SupportModel != "" {
		c.SupportModelDir = params.SupportModel
	}
	if params.Transpose != 0 {
		c.Transpose = params.Transpose
	}
	if params.Thread > 0 {
		c.NumThreads = params.Thread
	}

	dir := task.NewDir(workDir, neutrinoDir, neutrinoDir, startTime)
	if err := c.SetInfo(ctx, dir.NeutrinoDir()); err != nil {
		return nil, fmt.Errorf("failed to detect neutrino engine info: %w", err)
	}

	tasks, err := task.NewGenerator(dir, c, "", "").ExecutableTasks()
	if err != nil {
		return nil, fmt.Errorf("failed to generate executable tasks: %w", err)
	}

	taskNames := make([]string, len(tasks.Tasks))
	for i, t := range tasks.Tasks {
		taskNames[i] = t.Name
	}
	tasks.Entrypoint = taskNames
	tasks.Env.Set("PWD", dir.NeutrinoDir())

	var logBuf bytes.Buffer
	shell := "bash"
	runErr := tasks.IntoScript(shell).Runner(func(cmd *execx.Cmd) error {
		cmd.Dir = dir.NeutrinoDir()
		cmd.Stdout = &logBuf
		cmd.Stderr = &logBuf
		slog.Debug("executing standalone task", slog.String("dir", cmd.Dir), slog.Any("args", cmd.Args))
		return cmd.Exec()
	})

	logOutput := logBuf.String()
	if runErr != nil {
		return &SynthesizeResult{
			Mode:        string(ModeStandalone),
			Status:      "failed",
			Log:         logOutput,
			ElapsedSec:  time.Since(startTime).Seconds(),
			ErrorDetail: runErr.Error(),
		}, nil
	}

	resultDestDir := dir.ResultDestDir()
	basename := pathx.Basename(absScorePath)
	wavFile := filepath.Join(resultDestDir, fmt.Sprintf("%s.wav", basename))

	return &SynthesizeResult{
		Mode:       string(ModeStandalone),
		Status:     "succeed",
		WavPath:    wavFile,
		ResultDir:  resultDestDir,
		Log:        logOutput,
		ElapsedSec: time.Since(startTime).Seconds(),
	}, nil
}

func (r *StandaloneRunner) CheckProcess(_ context.Context, _ *CheckProcessParams) (*CheckProcessResult, error) {
	return &CheckProcessResult{
		Status:  "unsupported",
		Message: "check_process is only applicable for 'api' mode",
	}, nil
}
