package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	mcpserver "github.com/berquerant/pneutrinoutil/mcp/server"
	"github.com/berquerant/pneutrinoutil/pkg/logx"
	"github.com/berquerant/pneutrinoutil/pkg/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pneutrinoutil-mcp",
	Short: "Model Context Protocol (MCP) server for pneutrinoutil and NEUTRINO",
	Long: `pneutrinoutil-mcp runs an MCP server over stdio for AI assistants (such as Claude Desktop,
Antigravity, Cursor, etc.). It exposes tools and resources for NEUTRINO singing voice synthesis,
supporting both local standalone execution and remote pneutrinoutil-server REST API execution.`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		if v, _ := cmd.Flags().GetBool("version"); v {
			version.Write(os.Stdout)
			return nil
		}

		// Configure logging: ALWAYS to stderr to prevent corrupting stdio JSON-RPC protocol
		debugEnabled, _ := cmd.Flags().GetBool("debug")
		logLevel := slog.LevelInfo
		if debugEnabled {
			logLevel = slog.LevelDebug
		}
		logger := logx.NewTextLogger(os.Stderr, logLevel)
		slog.SetDefault(logger)

		cfg := mcpserver.NewDefaultConfig()

		modeStr, _ := cmd.Flags().GetString("mode")
		if modeStr != "" {
			cfg.Mode = mcpserver.Mode(modeStr)
		}
		neutrinoDir, _ := cmd.Flags().GetString("neutrinoDir")
		if neutrinoDir != "" {
			cfg.NeutrinoDir = neutrinoDir
		}
		workDir, _ := cmd.Flags().GetString("workDir")
		if workDir != "" {
			cfg.WorkDir = workDir
		}
		serverURI, _ := cmd.Flags().GetString("server")
		if serverURI != "" {
			cfg.ServerURI = serverURI
		}

		ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		srv := mcpserver.New(cfg)
		return mcpserver.RunStdio(ctx, srv, os.Stderr)
	},
}

func init() {
	defaultCfg := mcpserver.NewDefaultConfig()

	rootCmd.Flags().Bool("version", false, "print version")
	rootCmd.Flags().Bool("debug", false, "enable debug logging to stderr")
	rootCmd.Flags().StringP("mode", "m", string(defaultCfg.Mode), "default execution mode: 'standalone' (local direct) or 'api' (REST API)")
	rootCmd.Flags().StringP("neutrinoDir", "n", defaultCfg.NeutrinoDir, "default NEUTRINO directory for standalone mode")
	rootCmd.Flags().StringP("workDir", "w", defaultCfg.WorkDir, "working directory for intermediate files")
	rootCmd.Flags().StringP("server", "s", defaultCfg.ServerURI, "default pneutrinoutil-server URI for api mode")
}

func main() {
	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
