package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/berquerant/pneutrinoutil/cli/info"
	"github.com/berquerant/pneutrinoutil/pkg/pathx"
)

type APIRunner struct {
	cfg *Config
}

func NewAPIRunner(cfg *Config) *APIRunner {
	return &APIRunner{cfg: cfg}
}

func (r *APIRunner) GetInfo(ctx context.Context, opt *ExecutionOptions) (*info.Info, error) {
	_, _, _, serverURI := r.cfg.ResolveEffectiveConfig(opt)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/version", serverURI), nil)
	if err != nil {
		return nil, err
	}

	resp, err := r.cfg.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to contact server at %s: %w", serverURI, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	var ver struct {
		Version  string `json:"version"`
		Revision string `json:"revision"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&ver); err != nil {
		return nil, err
	}

	return &info.Info{
		Version:  ver.Version,
		Revision: ver.Revision,
	}, nil
}

func (r *APIRunner) ListModels(_ context.Context, _ *ExecutionOptions) ([]info.Model, error) {
	return []info.Model{
		{ID: "MERROW"},
		{ID: "KIRITAN"},
		{ID: "ZUNDAMON"},
		{ID: "ITAKO"},
	}, nil
}

func (r *APIRunner) Synthesize(ctx context.Context, params *SynthesizeParams) (*SynthesizeResult, error) {
	startTime := time.Now()
	_, _, _, serverURI := r.cfg.ResolveEffectiveConfig(params.Options)

	requestID, err := r.submitScoreJob(ctx, serverURI, params)
	if err != nil {
		return nil, err
	}

	if !params.Wait {
		return &SynthesizeResult{
			Mode:       string(ModeAPI),
			Status:     "pending",
			RequestID:  requestID,
			ElapsedSec: time.Since(startTime).Seconds(),
			Message:    fmt.Sprintf("Job submitted successfully (requestId: %s). Use tool 'check_process' to check status or download WAV.", requestID),
		}, nil
	}

	return r.pollJobCompletion(ctx, params, requestID, startTime)
}

func (r *APIRunner) submitScoreJob(ctx context.Context, serverURI string, params *SynthesizeParams) (string, error) {
	body, contentType, err := buildScoreMultipartForm(params)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/proc", serverURI), body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", contentType)

	resp, err := r.cfg.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to %s: %w", serverURI, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusAccepted {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(respBody))
	}

	requestID := resp.Header.Get("X-Request-Id")
	if requestID == "" {
		return "", fmt.Errorf("server did not return X-Request-Id header")
	}
	return requestID, nil
}

func buildScoreMultipartForm(params *SynthesizeParams) (*bytes.Buffer, string, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("score", "score.musicxml")
	if err != nil {
		return nil, "", fmt.Errorf("failed to create multipart form: %w", err)
	}

	if err := writeScorePayload(part, params); err != nil {
		return nil, "", err
	}

	if params.Model != "" {
		_ = writer.WriteField("model", params.Model)
	}
	if params.SupportModel != "" {
		_ = writer.WriteField("supportModel", params.SupportModel)
	}
	if params.Transpose != 0 {
		_ = writer.WriteField("transpose", strconv.Itoa(params.Transpose))
	}
	if err := writer.Close(); err != nil {
		return nil, "", fmt.Errorf("failed to close multipart writer: %w", err)
	}
	return &body, writer.FormDataContentType(), nil
}

func writeScorePayload(w io.Writer, params *SynthesizeParams) error {
	if params.ScorePath != "" {
		f, err := os.Open(params.ScorePath)
		if err != nil {
			return fmt.Errorf("failed to open score file %s: %w", params.ScorePath, err)
		}
		defer func() { _ = f.Close() }()
		if _, err := io.Copy(w, f); err != nil {
			return fmt.Errorf("failed to copy score data: %w", err)
		}
		return nil
	}
	if params.ScoreContent != "" {
		if _, err := io.WriteString(w, params.ScoreContent); err != nil {
			return fmt.Errorf("failed to write score content: %w", err)
		}
		return nil
	}
	return fmt.Errorf("either scorePath or scoreContent must be provided")
}

func (r *APIRunner) pollJobCompletion(ctx context.Context, params *SynthesizeParams, requestID string, startTime time.Time) (*SynthesizeResult, error) {
	timeoutSec := 30
	if params.TimeoutSeconds > 0 {
		timeoutSec = params.TimeoutSeconds
	}

	pollTicker := time.NewTicker(1500 * time.Millisecond)
	defer pollTicker.Stop()
	timeoutTimer := time.NewTimer(time.Duration(timeoutSec) * time.Second)
	defer timeoutTimer.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-timeoutTimer.C:
			return &SynthesizeResult{
				Mode:       string(ModeAPI),
				Status:     "still_running",
				RequestID:  requestID,
				ElapsedSec: time.Since(startTime).Seconds(),
				Message:    fmt.Sprintf("Process is still running after %d seconds. Use tool 'check_process' with requestId '%s' to check status later.", timeoutSec, requestID),
			}, nil
		case <-pollTicker.C:
			detail, err := r.CheckProcess(ctx, &CheckProcessParams{
				RequestID:   requestID,
				DownloadWav: true,
				Options:     params.Options,
			})
			if err != nil {
				continue
			}
			if detail.Status == "succeed" {
				return &SynthesizeResult{
					Mode:       string(ModeAPI),
					Status:     "succeed",
					RequestID:  requestID,
					WavPath:    detail.WavPath,
					Log:        detail.Log,
					ElapsedSec: time.Since(startTime).Seconds(),
				}, nil
			}
			if detail.Status == "failed" {
				return &SynthesizeResult{
					Mode:        string(ModeAPI),
					Status:      "failed",
					RequestID:   requestID,
					Log:         detail.Log,
					ErrorDetail: "Synthesis process failed on worker",
					ElapsedSec:  time.Since(startTime).Seconds(),
				}, nil
			}
		}
	}
}

func (r *APIRunner) CheckProcess(ctx context.Context, params *CheckProcessParams) (*CheckProcessResult, error) {
	_, _, workDir, serverURI := r.cfg.ResolveEffectiveConfig(params.Options)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/proc/%s/detail", serverURI, params.RequestID), nil)
	if err != nil {
		return nil, err
	}

	resp, err := r.cfg.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch process detail: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return &CheckProcessResult{
			RequestID: params.RequestID,
			Status:    "not_found",
			Message:   fmt.Sprintf("Process %s not found on server", params.RequestID),
		}, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %d for detail", resp.StatusCode)
	}

	var detail struct {
		OK   bool `json:"ok"`
		Data struct {
			RequestID   string  `json:"requestId"`
			Status      string  `json:"status"`
			Title       string  `json:"title"`
			Command     *string `json:"command"`
			CreatedAt   string  `json:"createdAt"`
			StartedAt   *string `json:"startedAt"`
			CompletedAt *string `json:"completedAt"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&detail); err != nil {
		return nil, err
	}

	status := detail.Data.Status
	result := &CheckProcessResult{
		RequestID: params.RequestID,
		Status:    status,
	}

	if status == "succeed" || status == "failed" {
		logReq, _ := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/proc/%s/log", serverURI, params.RequestID), nil)
		if logResp, err := r.cfg.HTTPClient.Do(logReq); err == nil {
			if logResp.StatusCode == http.StatusOK {
				logBytes, _ := io.ReadAll(logResp.Body)
				result.Log = string(logBytes)
			}
			_ = logResp.Body.Close()
		}
	}

	if params.DownloadWav && status == "succeed" {
		wavReq, _ := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/proc/%s/wav", serverURI, params.RequestID), nil)
		if wavResp, err := r.cfg.HTTPClient.Do(wavReq); err == nil {
			if wavResp.StatusCode == http.StatusOK {
				destDir := filepath.Join(workDir, "downloads", params.RequestID)
				_ = pathx.EnsureDir(destDir)
				destFile := filepath.Join(destDir, fmt.Sprintf("%s.wav", params.RequestID))
				if outF, err := os.Create(destFile); err == nil {
					_, _ = io.Copy(outF, wavResp.Body)
					_ = outF.Close()
					result.WavPath = destFile
				}
			}
			_ = wavResp.Body.Close()
		}
	}

	return result, nil
}
