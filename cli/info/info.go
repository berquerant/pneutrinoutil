package info

import (
	"context"

	"github.com/berquerant/pneutrinoutil/pkg/version"
)

type Info struct {
	Version  string   `json:"version"`
	Revision string   `json:"revision"`
	Neutrino Neutrino `json:"neutrino"`
}

type Neutrino struct {
	Version string  `json:"version"`
	Models  []Model `json:"models"`
}

type Model struct {
	ID   string `json:"id"`
	Data any    `json:"data,omitempty"`
}

func NewBuilder(neutrinoDir string) *Builder {
	return &Builder{
		neutrinoDir: neutrinoDir,
	}
}

type Builder struct {
	neutrinoDir string
}

func (b Builder) Build(ctx context.Context) (*Info, error) {
	models, err := getModels(b.neutrinoDir)
	if err != nil {
		return nil, err
	}
	neutrinoVersion, err := getNeutrinoVersion(ctx, b.neutrinoDir)
	if err != nil {
		return nil, err
	}
	return &Info{
		Version:  version.Version,
		Revision: version.Revision,
		Neutrino: Neutrino{
			Version: neutrinoVersion,
			Models:  models,
		},
	}, nil
}
