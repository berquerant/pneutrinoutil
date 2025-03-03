package logx

import (
	"encoding/json"
	"fmt"
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

func Array[T any](k string, v ...T) slog.Attr {
	ss := make([]string, len(v))
	for i, x := range v {
		ss[i] = fmt.Sprint(x)
	}
	return slog.Any(k, ss)
}
