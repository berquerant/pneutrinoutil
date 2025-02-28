package pworker

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/berquerant/pneutrinoutil/cli/ctl"
	"github.com/berquerant/pneutrinoutil/pkg/pathx"
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
	if err != nil {
		return nil, err
	}
	return re, nil
}

func LoadResultConfig(dir string, elem *pathx.ResultElement) (*ctl.Config, error) {
	f, err := os.Open(filepath.Join(dir, elem.String(), "config.yml"))
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
