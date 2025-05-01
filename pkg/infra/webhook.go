package infra

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Webhooker interface {
	Webhook(ctx context.Context, v any) error
}

var _ Webhooker = &Webhook{}

func NewWebhook(endpoint string, timeout time.Duration) *Webhook {
	return &Webhook{
		endpoint: endpoint,
		timeout:  timeout,
	}
}

type Webhook struct {
	endpoint string
	timeout  time.Duration
}

func (w *Webhook) Webhook(ctx context.Context, v any) error {
	if w == nil || w.endpoint == "" {
		return nil
	}

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, w.endpoint, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Set("content-type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	return nil
}
