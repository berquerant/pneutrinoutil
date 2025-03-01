package logx

import (
	"encoding/json"
	"log/slog"
)

func Err(err error) slog.Attr {
	if err != nil {
		return slog.String("err", err.Error())
	}
	return slog.String("err", "nil")
}

func Jsonify(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}

func JSON(k string, v any) slog.Attr {
	return slog.String(k, string(Jsonify(v)))
}
