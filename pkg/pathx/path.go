package pathx

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

//go:generate go tool goconfig -field "Mode os.FileMode|Truncate bool" -option -output path_config_generated.go

const (
	defaultDirPermission  = 0755
	defaultFilePermission = 0644
)

func EnsureDir(path string, opt ...ConfigOption) error {
	slog.Debug("EnsureDir", slog.String("path", path))

	c := NewConfigBuilder().
		Mode(defaultDirPermission).
		Build()
	c.Apply(opt...)
	return os.MkdirAll(path, c.Mode.Get())
}

func EnsureFile(path string, opt ...ConfigOption) error {
	slog.Debug("EnsureFile", slog.String("path", path))

	c := NewConfigBuilder().
		Mode(defaultFilePermission).
		Truncate(false).
		Build()
	c.Apply(opt...)
	if err := EnsureDir(filepath.Dir(path)); err != nil {
		return err
	}
	if c.Truncate.Get() {
		if err := createOrTruncateFile(path); err != nil {
			return err
		}
	} else if err := createFileIfNotExists(path); err != nil {
		return err
	}
	return os.Chmod(path, c.Mode.Get())
}

func createOrTruncateFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	return f.Close()
}

func createFileIfNotExists(path string) error {
	_, err := os.Stat(path)
	switch {
	case err == nil:
		return nil
	case os.IsNotExist(err):
		return createOrTruncateFile(path)
	default:
		return err
	}
}

type ExistType int

const (
	Eunknown ExistType = iota
	EnotExist
	Edir
	Efile
)

func Exist(path string) ExistType {
	slog.Debug("Exist", slog.String("path", path))

	info, err := os.Stat(path)
	switch {
	case err == nil:
		if info.IsDir() {
			return Edir
		}
		return Efile
	case os.IsNotExist(err):
		return EnotExist
	default:
		return Eunknown
	}
}

func Basename(path string) string {
	b := filepath.Base(path)
	return strings.TrimSuffix(b, filepath.Ext(b))
}
