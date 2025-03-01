package logx

import (
	"encoding/json"
	"log/slog"
)

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
