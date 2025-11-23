package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	uploadMaxSizeBytes = 1 << 20 // 1 MiB
)

type ErrorResponse struct {
	OK    bool   `json:"ok"` // false
	Error string `json:"error"`
}

func NewErrorResponse(msg string) *ErrorResponse {
	return &ErrorResponse{
		OK:    false,
		Error: msg,
	}
}

type SuccessResponse[T any] struct {
	OK   bool `json:"ok"` // true
	Data T    `json:"data"`
}

func NewSuccessResponse[T any](data T) *SuccessResponse[T] {
	return &SuccessResponse[T]{
		OK:   true,
		Data: data,
	}
}

// Error sends a JSON error response with status and message.
func Error(c echo.Context, status int, msg string) error {
	return c.JSON(status, NewErrorResponse(
		msg,
	))
}

// Success sends a JSON successful response with status and data.
func Success[T any](c echo.Context, status int, data T) error {
	return c.JSON(status, NewSuccessResponse(
		data,
	))
}

type StatusError struct {
	Status int // http status
	Err    error
	Msg    string // response body
}

var _ error = &StatusError{}

func (e StatusError) Error() string { return e.Err.Error() }
func (e StatusError) Unwrap() error { return e.Err }

func (e StatusError) AppendMessageToErr(msg string) *StatusError {
	return &StatusError{
		Status: e.Status,
		Err:    fmt.Errorf("%w: %s", e.Err, msg),
		Msg:    e.Msg,
	}
}

func NewStatusError(status int, err error, msg string) *StatusError {
	return &StatusError{
		Status: status,
		Err:    err,
		Msg:    msg,
	}
}

// Respond sends a JSON error response.
func (e StatusError) Respond(c echo.Context) error {
	return Error(c, e.Status, e.Msg)
}

type ReadFromFileResult struct {
	Blob []byte
	Name string
}

// ReadFormFile reads a form file and writes content to uploadDir.
func ReadFormFile(c echo.Context, name string, maxBytes int64) (*ReadFromFileResult, *StatusError) {
	fh, err := c.FormFile(name)
	if err != nil {
		return nil, NewStatusError(
			http.StatusBadRequest,
			fmt.Errorf("%w: read multipart form file: %s", err, name),
			fmt.Sprintf("failed to read form file: %s", name),
		)
	}
	if fh.Size > maxBytes {
		return nil, NewStatusError(
			http.StatusRequestEntityTooLarge,
			fmt.Errorf("RequestEntityTooLarge: %s", fh.Filename),
			fmt.Sprintf("file is too big: %s", fh.Filename),
		)
	}

	src, err := fh.Open()
	if err != nil {
		return nil, NewStatusError(
			http.StatusInternalServerError,
			fmt.Errorf("%w: open %s", err, fh.Filename),
			fmt.Sprintf("failed to open form file: %s", fh.Filename),
		)
	}
	defer func() { _ = src.Close() }()

	blob, err := io.ReadAll(src)
	if err != nil {
		return nil, NewStatusError(
			http.StatusInternalServerError,
			fmt.Errorf("%w: read form file from %s", err, fh.Filename),
			fmt.Sprintf("failed to read form file: %s", fh.Filename),
		)
	}
	return &ReadFromFileResult{
		Blob: blob,
		Name: fh.Filename,
	}, nil
}

// - timestamp 1136189045
// - RFC3339 2006-01-02T15:04:05Z07:00
type CustomTime time.Time

var _ echo.BindUnmarshaler = new(CustomTime)

func (t *CustomTime) String() string {
	return time.Time(*t).String()
}

func (t *CustomTime) Timestamp() int64 {
	return time.Time(*t).Unix()
}

func (t *CustomTime) UnmarshalParam(param string) error {
	v, err := t.parseRFC3339(param)
	if err == nil {
		*t = CustomTime(v)
		return nil
	}
	tv, tErr := t.parseTimestamp(param)
	if tErr == nil {
		*t = CustomTime(tv)
		return nil
	}
	return errors.Join(err, tErr)
}

func (*CustomTime) parseTimestamp(param string) (time.Time, error) {
	var t time.Time
	ts, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return t, err
	}
	return time.Unix(ts, 0), nil
}

func (*CustomTime) parseRFC3339(param string) (time.Time, error) {
	return time.Parse(time.RFC3339, param)
}
