package ctl

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/berquerant/execx"
	"github.com/berquerant/pneutrinoutil/cli/info"
	"github.com/berquerant/pneutrinoutil/pkg/pathx"
	"github.com/berquerant/structconfig"
	"github.com/spf13/pflag"
)

// NewDefaultConfig returns a new Config.
// Default values are from "default" struct tags.
func NewDefaultConfig() (*Config, error) {
	sc := structconfig.New[Config]()
	var c Config
	if err := sc.FromDefault(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

type Config struct {
	Description string `json:"desc" yaml:"desc" name:"desc" usage:"description of config"`
	NumParallel int    `json:"parallel" yaml:"parallel" name:"parallel" usage:"number of parallel" default:"1"`
	// musicXML_to_label
	// Suffix string `yaml:"suffix"`
	// Project settings
	Score         string `json:"score" yaml:"score" name:"score" usage:"score file, required"`
	NumThreads    int    `json:"thread" yaml:"thread" name:"thread" usage:"number of parallel in session" default:"4"`
	InferenceMode int    `json:"inference" yaml:"inference" name:"inference" usage:"quality, processing speed: 2 (elements), 3 (standard) or 4 (advanced)" default:"3"`
	// NEUTRINO
	ModelDir   string `json:"model" yaml:"model" name:"model" usage:"singer" default:"MERROW"`
	StyleShift int    `json:"styleShift" yaml:"styleShift" name:"styleShift" usage:"change the key and estimate to change the style of singing"`
	RandomSeed int    `json:"randomSeed" yaml:"randomSeed" name:"randomSeed" usage:"random seed" default:"1234"`
	// NSF
	PitchShiftNsf float32 `json:"pitchShiftNsf" yaml:"pitchShiftNsf" name:"pitchShiftNsf" usage:"change pitch via NSF"`
	// WORLD
	PitchShiftWorld    float32 `json:"pitchShiftWorld" yaml:"pitchShiftWorld" name:"pitchShiftWorld" usage:"change pitch via WORLD"`
	FormantShift       float32 `json:"formantShift" yaml:"formantShift" name:"formantShift" usage:"change voice quality" default:"1.0"`
	SmoothPitch        float32 `json:"smoothPitch" yaml:"smoothPitch" name:"smoothPitch" usage:"[0, 100]%"`
	SmoothFormant      float32 `json:"smoothFormant" yaml:"smoothFormant" name:"smoothFormant" usage:"[0, 100]%"`
	EnhanceBreathiness float32 `json:"enhanceBreathiness" yaml:"enhanceBreathiness" name:"enhanceBreathiness" usage:"[0, 100]%"`
	// Info
	NeutrinoVersion string `json:"neutrinoVersion" yaml:"neutrinoVersion"`
	ModelData       any    `json:"modelData" yaml:"modelData"`
}

func (c Config) Basename() string { return pathx.Basename(c.Score) }

func (c Config) nsfModel() string {
	switch c.InferenceMode {
	case 4:
		return "va"
	case 3:
		return "vs"
	case 2:
		return "ve"
	default:
		return "ve"
	}
}

func (c Config) samplingFreq() int {
	switch c.InferenceMode {
	case 4, 3:
		return 48
	case 2:
		return 24
	default:
		return 24
	}
}

func (c Config) envMap() map[string]any {
	return map[string]any{
		"NumThreads":         c.NumThreads,
		"InferenceMode":      c.InferenceMode,
		"ModelDir":           c.ModelDir,
		"StyleShift":         c.StyleShift,
		"PitchShiftNsf":      c.PitchShiftNsf,
		"PitchShiftWorld":    c.PitchShiftWorld,
		"FormantShift":       c.FormantShift,
		"SmoothPitch":        c.SmoothPitch,
		"SmoothFormant":      c.SmoothFormant,
		"EnhanceBreathiness": c.EnhanceBreathiness,
		"NsfModel":           c.nsfModel(),
		"SamplingFreq":       c.samplingFreq(),
		"BASENAME":           c.Basename(),
		"RandomSeed":         c.RandomSeed,
		"NumParallel":        c.NumParallel,
	}
}

func (c Config) Env() execx.Env {
	e := execx.NewEnv()
	for k, v := range c.envMap() {
		e.Set(k, fmt.Sprint(v))
	}
	return e
}

// ApplyFlagValues sets flag values to this.
func (c *Config) ApplyFlagValues(fs *pflag.FlagSet) error {
	sc := structconfig.New[Config]()
	return sc.FromFlags(c, fs)
}

// SetFlags sets command-line flags.
// Flag name is from "name" struct tag.
// Flag usage is from "usage" struct tag.
// Flag default value is from "default" struct tag.
func (c Config) SetFlags(fs *pflag.FlagSet) error {
	sc := structconfig.New[Config]()
	return sc.SetFlags(fs)
}

func (c *Config) SetInfo(ctx context.Context, neutrinoDir string) error {
	dir, err := filepath.Abs(neutrinoDir)
	if err != nil {
		return err
	}

	neutrinoVersion, err := info.GetNeutrinoVersion(ctx, dir)
	if err != nil {
		return err
	}
	c.NeutrinoVersion = neutrinoVersion

	model, err := info.ReadModelInfo(filepath.Join(dir, "model", c.ModelDir))
	if err != nil {
		return err
	}
	c.ModelData = model
	return nil
}
