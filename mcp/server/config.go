package server

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/berquerant/pneutrinoutil/cli/info"
	"github.com/berquerant/pneutrinoutil/pkg/pathx"
)

// Mode represents the execution mode for synthesis or information querying.
type Mode string

const (
	ModeStandalone Mode = "standalone"
	ModeAPI        Mode = "api"
)

// Config holds the runtime configuration for the MCP server.
type Config struct {
	Mode        Mode
	NeutrinoDir string
	WorkDir     string
	ServerURI   string
	HTTPClient  *http.Client
}

func NewDefaultConfig() *Config {
	workDir := filepath.Join(pathx.UserHomeDirOr("."), ".pneutrinoutil")
	serverURI := os.Getenv("SERVER_URI")
	if serverURI == "" {
		serverURI = "http://127.0.0.1:9101/v1"
	}

	return &Config{
		Mode:        ModeStandalone,
		NeutrinoDir: "./dist/NEUTRINO",
		WorkDir:     workDir,
		ServerURI:   serverURI,
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// ExecutionOptions provides optional overrides for a specific tool call.
type ExecutionOptions struct {
	Mode        *string `json:"mode,omitempty" jsonschema:"Execution mode: 'standalone' (local direct NEUTRINO) or 'api' (pneutrinoutil-server REST API)"`
	NeutrinoDir *string `json:"neutrinoDir,omitempty" jsonschema:"Override NEUTRINO engine directory (standalone mode)"`
	WorkDir     *string `json:"workDir,omitempty" jsonschema:"Override working directory for intermediate files"`
	ServerURI   *string `json:"serverUri,omitempty" jsonschema:"Override REST API server URI (api mode)"`
}

// ResolveEffectiveConfig resolves the effective configuration for an operation.
func (c *Config) ResolveEffectiveConfig(opt *ExecutionOptions) (Mode, string, string, string) {
	mode := c.Mode
	neutrinoDir := c.NeutrinoDir
	workDir := c.WorkDir
	serverURI := c.ServerURI

	if opt != nil {
		if opt.Mode != nil && *opt.Mode != "" {
			mode = Mode(*opt.Mode)
		}
		if opt.NeutrinoDir != nil && *opt.NeutrinoDir != "" {
			neutrinoDir = *opt.NeutrinoDir
		}
		if opt.WorkDir != nil && *opt.WorkDir != "" {
			workDir = *opt.WorkDir
		}
		if opt.ServerURI != nil && *opt.ServerURI != "" {
			serverURI = *opt.ServerURI
		}
	}

	return mode, neutrinoDir, workDir, serverURI
}

// Runner is the common interface for running synthesis operations.
type Runner interface {
	Synthesize(ctx context.Context, params *SynthesizeParams) (*SynthesizeResult, error)
	GetInfo(ctx context.Context, opt *ExecutionOptions) (*info.Info, error)
	ListModels(ctx context.Context, opt *ExecutionOptions) ([]info.Model, error)
	CheckProcess(ctx context.Context, params *CheckProcessParams) (*CheckProcessResult, error)
}
