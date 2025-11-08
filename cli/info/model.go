package info

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/berquerant/pneutrinoutil/pkg/logx"
)

var ErrModels = errors.New("Models")

func getModels(neutrinoDir string) ([]Model, error) {
	x, err := filepath.Abs(neutrinoDir)
	if err != nil {
		return nil, errors.Join(ErrModels, err)
	}
	dir := filepath.Join(x, "model")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, errors.Join(ErrModels, err)
	}

	var models []Model
	for _, x := range entries {
		if !x.IsDir() {
			continue
		}
		m := Model{
			ID: x.Name(),
		}
		data, err := readModelInfo(filepath.Join(dir, x.Name()))
		if err == nil {
			m.Data = data
		} else {
			slog.Warn("failed to read model info", slog.String("model", m.ID), logx.Err(err))
		}
		models = append(models, m)
	}
	return models, nil
}

func readModelInfo(modelDir string) (map[string]any, error) {
	var m map[string]any
	if _, err := toml.DecodeFile(filepath.Join(modelDir, "info.toml"), &m); err != nil {
		return nil, err
	}
	return m, nil
}
