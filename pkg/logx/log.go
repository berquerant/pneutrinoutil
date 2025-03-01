package logx

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"
)

func NewTextLogger(w io.Writer, level slog.Leveler) *slog.Logger {
	handler := slog.NewTextHandler(w, &slog.HandlerOptions{
		Level: level,
	})
	return slog.New(handler)
}

func NewLTSVLogger(w io.Writer, level slog.Leveler) *slog.Logger {
	handler := NewLTSVHandler(w, level)
	return slog.New(handler)
}

type LTSVHandler struct {
	w     io.Writer
	level slog.Leveler
}

func NewLTSVHandler(w io.Writer, level slog.Leveler) *LTSVHandler {
	return &LTSVHandler{
		w:     w,
		level: level,
	}
}

func (h *LTSVHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level.Level() >= h.level.Level()
}
func (h *LTSVHandler) WithAttrs(_ []slog.Attr) slog.Handler { return h }
func (h *LTSVHandler) WithGroup(_ string) slog.Handler      { return h }

func (h *LTSVHandler) Handle(_ context.Context, r slog.Record) error {
	var (
		attrs = make([]string, r.NumAttrs()+3)
		i     int
		add   = func(k, v string) {
			attrs[i] = fmt.Sprintf("%s=%s", k, v)
			i++
		}
	)

	add("time", r.Time.Format(time.RFC3339Nano))
	add("level", r.Level.String())
	add("msg", r.Message)
	r.Attrs(func(attr slog.Attr) bool {
		add(attr.Key, attr.Value.String())
		return true
	})
	_, err := fmt.Fprintln(h.w, strings.Join(attrs, "\t"))
	return err
}
