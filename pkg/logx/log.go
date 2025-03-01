package logx

import (
	"encoding/json"
	"io"
	"log/slog"
)

func NewTextLogger(w io.Writer, level slog.Leveler) *slog.Logger {
	handler := slog.NewTextHandler(w, &slog.HandlerOptions{
		Level: level,
	})
	return slog.New(handler)
}

func Err(err error) slog.Attr {
	return slog.Any("err", err)
}

func Jsonify(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}

func JSON(k string, v any) slog.Attr {
	return slog.String(k, string(Jsonify(v)))
}
