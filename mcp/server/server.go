package server

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/berquerant/pneutrinoutil/pkg/version"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// New creates and configures the MCP server instance.
func New(cfg *Config) *mcp.Server {
	impl := &mcp.Implementation{
		Name:    "pneutrinoutil",
		Version: version.Version,
	}

	opts := &mcp.ServerOptions{
		Instructions: "pneutrinoutil MCP Server provides singing voice synthesis capabilities via STUDIO NEUTRINO. You can inspect available singer models with list_models, query engine info with get_neutrino_info, and synthesize WAV audio from MusicXML scores with synthesize. Supports both standalone direct execution and remote REST API mode.",
	}

	server := mcp.NewServer(impl, opts)
	RegisterTools(server, cfg)
	RegisterResources(server, cfg)

	return server
}

// RunStdio starts the MCP server over standard I/O.
// All application logs are directed to logDest (typically os.Stderr) to ensure
// standard output is reserved strictly for JSON-RPC framing.
func RunStdio(ctx context.Context, server *mcp.Server, logDest io.Writer) error {
	t := &mcp.LoggingTransport{
		Transport: &mcp.StdioTransport{},
		Writer:    logDest,
	}

	slog.Info("starting pneutrinoutil MCP server over stdio", slog.String("version", version.Version))
	if err := server.Run(ctx, t); err != nil {
		return fmt.Errorf("mcp server error: %w", err)
	}

	return nil
}
