package pworker

import (
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/berquerant/pneutrinoutil/cli/ctl"
	"github.com/berquerant/pneutrinoutil/pkg/logx"
	"github.com/berquerant/pneutrinoutil/pkg/pathx"
	"github.com/berquerant/pneutrinoutil/server/alog"
	"gopkg.in/yaml.v3"
)

var (
	ErrNotDir = errors.New("NotDir")
)

func LoadResultElement(entry os.DirEntry) (*pathx.ResultElement, error) {
	if !entry.IsDir() {
		return nil, ErrNotDir
	}
	re, err := pathx.ParseResultElement(entry.Name())
	alog.L().Debug("load result element", slog.String("name", entry.Name()), logx.Err(err))
	if err != nil {
		return nil, err
	}
	return re, nil
}

func LoadResultConfig(dir string, elem *pathx.ResultElement) (*ctl.Config, error) {
	name := filepath.Join(dir, elem.String(), "config.yml")
	f, err := os.Open(name)
	alog.L().Debug("load result config", slog.String("name", name), logx.Err(err))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var c ctl.Config
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, err
	}

	return &c, nil
}
