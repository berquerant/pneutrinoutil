package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/berquerant/pneutrinoutil/cli/ctl"
	"github.com/goccy/go-yaml"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type SynthesizeParams struct {
	ScorePath      string            `json:"scorePath,omitempty" jsonschema:"Path to the input MusicXML score file (.musicxml)"`
	ScoreContent   string            `json:"scoreContent,omitempty" jsonschema:"Raw MusicXML content string if file path is not available"`
	Model          string            `json:"model,omitempty" jsonschema:"Primary singer model name (e.g. MERROW, KIRITAN; default: MERROW)"`
	SupportModel   string            `json:"supportModel,omitempty" jsonschema:"Secondary support model name for voice blending in NEUTRINO v3"`
	Transpose      int               `json:"transpose,omitempty" jsonschema:"Key transpose offset in semitones (e.g. -2, 0, 3)"`
	Thread         int               `json:"thread,omitempty" jsonschema:"Parallel threads for synthesis session (default: 4)"`
	Wait           bool              `json:"wait,omitempty" jsonschema:"If true, wait for completion up to timeoutSeconds. If false, returns immediately with requestId (recommended for API mode)"`
	TimeoutSeconds int               `json:"timeoutSeconds,omitempty" jsonschema:"Maximum seconds to wait if wait is true (default: 30)"`
	Options        *ExecutionOptions `json:"options,omitempty" jsonschema:"Optional execution mode overrides (mode: 'standalone' or 'api', neutrinoDir, serverUri)"`
}

type SynthesizeResult struct {
	Mode        string  `json:"mode"`
	Status      string  `json:"status"` // succeed | pending | still_running | failed
	RequestID   string  `json:"requestId,omitempty"`
	WavPath     string  `json:"wavPath,omitempty"`
	ResultDir   string  `json:"resultDir,omitempty"`
	Log         string  `json:"log,omitempty"`
	ElapsedSec  float64 `json:"elapsedSeconds"`
	Message     string  `json:"message,omitempty"`
	ErrorDetail string  `json:"errorDetail,omitempty"`
}

type ListModelsParams struct {
	Options *ExecutionOptions `json:"options,omitempty" jsonschema:"Optional execution mode overrides"`
}

type GetInfoParams struct {
	Options *ExecutionOptions `json:"options,omitempty" jsonschema:"Optional execution mode overrides"`
}

type GetConfigSkeletonParams struct {
	JSONFormat bool `json:"jsonFormat,omitempty" jsonschema:"Output skeleton configuration as JSON instead of YAML (default: false)"`
}

type CheckProcessParams struct {
	RequestID   string            `json:"requestId" jsonschema:"Process Request ID to check"`
	DownloadWav bool              `json:"downloadWav,omitempty" jsonschema:"If true and process succeeded, download WAV file to local work directory"`
	Options     *ExecutionOptions `json:"options,omitempty" jsonschema:"Optional execution mode overrides"`
}

type CheckProcessResult struct {
	RequestID string `json:"requestId"`
	Status    string `json:"status"` // pending | running | succeed | failed | not_found | unsupported
	WavPath   string `json:"wavPath,omitempty"`
	Log       string `json:"log,omitempty"`
	Message   string `json:"message,omitempty"`
}

func RegisterTools(server *mcp.Server, cfg *Config) {
	standaloneRunner := NewStandaloneRunner(cfg)
	apiRunner := NewAPIRunner(cfg)

	getRunner := func(opt *ExecutionOptions) Runner {
		mode, _, _, _ := cfg.ResolveEffectiveConfig(opt)
		if mode == ModeAPI {
			return apiRunner
		}
		return standaloneRunner
	}

	// 1. synthesize tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "synthesize",
		Description: "Synthesize singing voice WAV audio from a MusicXML score using NEUTRINO. Supports local direct execution ('standalone' mode) or via pneutrinoutil-server REST API ('api' mode). In API mode, returns immediately with a requestId unless wait=true is specified.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, params SynthesizeParams) (*mcp.CallToolResult, any, error) {
		runner := getRunner(params.Options)
		res, err := runner.Synthesize(ctx, &params)
		if err != nil {
			return nil, nil, fmt.Errorf("synthesis failed: %w", err)
		}

		resBytes, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			return nil, nil, err
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: string(resBytes)},
			},
		}, res, nil
	})

	// 2. list_models tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_models",
		Description: "List available NEUTRINO singer voice models (e.g. MERROW, KIRITAN, ZUNDAMON). In standalone mode, scans local model directory; in API mode, returns known available models.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, params ListModelsParams) (*mcp.CallToolResult, any, error) {
		runner := getRunner(params.Options)
		models, err := runner.ListModels(ctx, params.Options)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to list models: %w", err)
		}

		modelsBytes, err := json.MarshalIndent(models, "", "  ")
		if err != nil {
			return nil, nil, err
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: string(modelsBytes)},
			},
		}, models, nil
	})

	// 3. get_neutrino_info tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_neutrino_info",
		Description: "Get version and system information about the NEUTRINO synthesis engine and pneutrinoutil setup.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, params GetInfoParams) (*mcp.CallToolResult, any, error) {
		runner := getRunner(params.Options)
		inf, err := runner.GetInfo(ctx, params.Options)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get info: %w", err)
		}

		infBytes, err := json.MarshalIndent(inf, "", "  ")
		if err != nil {
			return nil, nil, err
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: string(infBytes)},
			},
		}, inf, nil
	})

	// 4. get_config_skeleton tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_config_skeleton",
		Description: "Get default synthesis configuration template (YAML or JSON) showing all supported parameters.",
	}, func(_ context.Context, _ *mcp.CallToolRequest, params GetConfigSkeletonParams) (*mcp.CallToolResult, any, error) {
		c, err := ctl.NewDefaultConfig()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create default config: %w", err)
		}

		var out []byte
		if params.JSONFormat {
			out, err = json.MarshalIndent(c, "", "  ")
		} else {
			out, err = yaml.Marshal(c)
		}
		if err != nil {
			return nil, nil, err
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: string(out)},
			},
		}, c, nil
	})

	// 5. check_process tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "check_process",
		Description: "Check the status of a background synthesis job submitted in API mode by its requestId. If finished with success and downloadWav=true, downloads the WAV file.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, params CheckProcessParams) (*mcp.CallToolResult, any, error) {
		runner := getRunner(params.Options)
		res, err := runner.CheckProcess(ctx, &params)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to check process: %w", err)
		}

		resBytes, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			return nil, nil, err
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: string(resBytes)},
			},
		}, res, nil
	})
}

func RegisterResources(server *mcp.Server, cfg *Config) {
	server.AddResource(&mcp.Resource{
		Name:        "info",
		MIMEType:    "application/json",
		URI:         "neutrino://info",
		Description: "Current system information and NEUTRINO version",
	}, func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		runner := NewStandaloneRunner(cfg)
		inf, err := runner.GetInfo(ctx, nil)
		if err != nil {
			inf = nil
		}
		b, _ := json.Marshal(inf)
		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:  req.Params.URI,
					Text: string(b),
				},
			},
		}, nil
	})

	server.AddResource(&mcp.Resource{
		Name:        "default_config",
		MIMEType:    "application/x-yaml",
		URI:         "neutrino://config/default",
		Description: "Default config.yml skeleton for NEUTRINO synthesis",
	}, func(_ context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		c, err := ctl.NewDefaultConfig()
		if err != nil {
			return nil, err
		}
		b, err := yaml.Marshal(c)
		if err != nil {
			return nil, err
		}
		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:  req.Params.URI,
					Text: string(b),
				},
			},
		}, nil
	})
}
