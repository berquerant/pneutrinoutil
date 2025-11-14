package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"strings"
	"time"

	"github.com/berquerant/pneutrinoutil/pkg/logx"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.Flags().String("server", os.Getenv("SERVER_URI"), "pneutrinoutil-server URI")
	rootCmd.Flags().StringP("content", "c", "dummy musicxml by gendata", "musicxml content; @filename is available")
}

var rootCmd = &cobra.Command{
	Use:   "gendata",
	Short: "generate artifacts",
	Long:  `generate artifacts. read basenames from stdin`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt)
		defer stop()

		logger := logx.NewLTSVLogger(os.Stderr, slog.LevelDebug)
		filenameOrContent, _ := cmd.Flags().GetString("content")
		content, err := readContent(filenameOrContent)
		if err != nil {
			return err
		}

		var (
			ids     = map[string]string{}
			url, _  = cmd.Flags().GetString("server")
			scanner = bufio.NewScanner(os.Stdin)
		)
		for scanner.Scan() {
			basename := scanner.Text()
			rid, err := post(ctx, url, basename, content)
			if err != nil {
				return err
			}
			logger.Info("post", slog.String("basename", basename), slog.String("rid", rid))
			ids[rid] = basename
		}

		for rid, basename := range ids {
			lg := logger.With(slog.String("rid", rid), slog.String("basename", basename))
			lg.Info("wait")
			if err := wait(ctx, url, rid); err != nil {
				lg.Error("wait", slog.String("err", err.Error()))
				continue
			}
			fmt.Printf("%s %s\n", rid, basename)
		}
		return nil
	},
}

func readContent(filenameOrContent string) (string, error) {
	if !strings.HasPrefix(filenameOrContent, "@") {
		return filenameOrContent, nil
	}
	b, err := os.ReadFile(filenameOrContent[1:])
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func wait(ctx context.Context, url, rid string) error {
	type responseData struct {
		OK   bool `json:"ok"`
		Data struct {
			Status string `json:"status"`
		} `json:"data"`
	}
	for {
		time.Sleep(time.Millisecond * 300)
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url+"/proc/"+rid+"/detail", nil)
			if err != nil {
				return err
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				continue
			}
			b, err := io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			if err != nil {
				continue
			}
			var data responseData
			if err := json.Unmarshal(b, &data); err != nil {
				continue
			}
			if !data.OK {
				continue
			}
			if slices.Contains([]string{"succeed", "failed"}, data.Data.Status) {
				return nil
			}
		}
	}
}

func post(ctx context.Context, url, basename, content string) (string, error) {
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	fw, err := w.CreateFormFile("score", basename+".musicxml")
	if err != nil {
		return "", err
	}
	if _, err := fmt.Fprint(fw, content); err != nil {
		return "", err
	}
	if err := w.Close(); err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url+"/proc", &body)
	if err != nil {
		return "", err
	}
	req.Header.Set("content-type", w.FormDataContentType())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if status := resp.StatusCode; status != http.StatusAccepted {
		return "", fmt.Errorf("status: want %d got %d", http.StatusAccepted, status)
	}
	return resp.Header.Get("x-request-id"), nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
