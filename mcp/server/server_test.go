package server_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	mcpserver "github.com/berquerant/pneutrinoutil/mcp/server"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
)

func TestMCPServer_ToolsAndResources(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Setup a mock REST API server to simulate pneutrinoutil-server
	mockAPIServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/version":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"version":"0.9.0","revision":"mock123"}`))
		case "/proc":
			if r.Method == http.MethodPost {
				w.Header().Set("X-Request-Id", "mock_proc_12345")
				w.WriteHeader(http.StatusAccepted)
				_, _ = w.Write([]byte(`{"ok":true,"data":"accepted"}`))
			}
		case "/proc/mock_proc_12345/detail":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"ok": true,
				"data": {
					"requestId": "mock_proc_12345",
					"status": "succeed",
					"title": "sample"
				}
			}`))
		case "/proc/mock_proc_12345/wav":
			w.Header().Set("Content-Type", "audio/wav")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("RIFFmockWAVdata"))
		case "/proc/mock_proc_12345/log":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("mock execution log"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer mockAPIServer.Close()

	// In-memory MCP transport
	clientTransport, serverTransport := mcp.NewInMemoryTransports()

	// Create dummy neutrino dir
	tmpDir := t.TempDir()
	modelDir := filepath.Join(tmpDir, "model", "TEST_SINGER")
	assert.Nil(t, os.MkdirAll(modelDir, 0755))
	assert.Nil(t, os.WriteFile(filepath.Join(modelDir, "info.toml"), []byte("name = 'Test Singer'\n"), 0644))

	binDir := filepath.Join(tmpDir, "bin")
	assert.Nil(t, os.MkdirAll(binDir, 0755))
	dummyBin := filepath.Join(binDir, "neutrino")
	assert.Nil(t, os.WriteFile(dummyBin, []byte("#!/bin/sh\necho 'NEUTRINO - v3.0.0'\n"), 0755))

	// Config
	cfg := &mcpserver.Config{
		Mode:        mcpserver.ModeStandalone,
		NeutrinoDir: tmpDir,
		WorkDir:     filepath.Join(tmpDir, "work"),
		ServerURI:   mockAPIServer.URL,
		HTTPClient:  mockAPIServer.Client(),
	}

	srv := mcpserver.New(cfg)
	serverSession, err := srv.Connect(ctx, serverTransport, nil)
	assert.Nil(t, err)
	defer serverSession.Close()

	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "1.0"}, nil)
	clientSession, err := client.Connect(ctx, clientTransport, nil)
	assert.Nil(t, err)
	defer clientSession.Close()

	t.Run("list tools", func(t *testing.T) {
		tools, err := clientSession.ListTools(ctx, nil)
		assert.Nil(t, err)
		assert.NotNil(t, tools)

		toolNames := make([]string, len(tools.Tools))
		for i, tool := range tools.Tools {
			toolNames[i] = tool.Name
		}
		assert.Contains(t, toolNames, "synthesize")
		assert.Contains(t, toolNames, "list_models")
		assert.Contains(t, toolNames, "get_neutrino_info")
		assert.Contains(t, toolNames, "get_config_skeleton")
		assert.Contains(t, toolNames, "check_process")
	})

	t.Run("call get_config_skeleton", func(t *testing.T) {
		res, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
			Name: "get_config_skeleton",
			Arguments: map[string]any{
				"jsonFormat": true,
			},
		})
		assert.Nil(t, err)
		assert.NotEmpty(t, res.Content)
		text := res.Content[0].(*mcp.TextContent).Text
		assert.Contains(t, text, "model")
	})

	t.Run("call list_models in standalone mode", func(t *testing.T) {
		res, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
			Name: "list_models",
		})
		assert.Nil(t, err)
		assert.NotEmpty(t, res.Content)
		text := res.Content[0].(*mcp.TextContent).Text
		assert.Contains(t, text, "TEST_SINGER")
	})

	t.Run("call list_models with mode=api override", func(t *testing.T) {
		apiMode := "api"
		res, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
			Name: "list_models",
			Arguments: map[string]any{
				"options": map[string]any{
					"mode": apiMode,
				},
			},
		})
		assert.Nil(t, err)
		assert.NotEmpty(t, res.Content)
		text := res.Content[0].(*mcp.TextContent).Text
		assert.Contains(t, text, "MERROW")
	})

	t.Run("call synthesize in api mode with wait=false", func(t *testing.T) {
		apiMode := "api"
		res, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
			Name: "synthesize",
			Arguments: map[string]any{
				"scoreContent": "<score-partwise/>",
				"model":        "MERROW",
				"wait":         false,
				"options": map[string]any{
					"mode": apiMode,
				},
			},
		})
		assert.Nil(t, err)
		assert.NotEmpty(t, res.Content)
		text := res.Content[0].(*mcp.TextContent).Text

		var out mcpserver.SynthesizeResult
		assert.Nil(t, json.Unmarshal([]byte(text), &out))
		assert.Equal(t, "pending", out.Status)
		assert.Equal(t, "mock_proc_12345", out.RequestID)
	})

	t.Run("call check_process in api mode", func(t *testing.T) {
		apiMode := "api"
		res, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
			Name: "check_process",
			Arguments: map[string]any{
				"requestId":   "mock_proc_12345",
				"downloadWav": true,
				"options": map[string]any{
					"mode": apiMode,
				},
			},
		})
		assert.Nil(t, err)
		assert.NotEmpty(t, res.Content)
		text := res.Content[0].(*mcp.TextContent).Text

		var out mcpserver.CheckProcessResult
		assert.Nil(t, json.Unmarshal([]byte(text), &out))
		assert.Equal(t, "succeed", out.Status)
		assert.FileExists(t, out.WavPath)
		assert.Equal(t, "mock execution log", out.Log)
	})

	t.Run("read resource neutrino://info", func(t *testing.T) {
		res, err := clientSession.ReadResource(ctx, &mcp.ReadResourceParams{
			URI: "neutrino://info",
		})
		assert.Nil(t, err)
		assert.NotEmpty(t, res.Contents)
	})
}
