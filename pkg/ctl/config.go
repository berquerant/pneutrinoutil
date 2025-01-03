package ctl

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/berquerant/execx"
	"github.com/berquerant/structconfig"
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
	Description string `yaml:"desc" name:"desc" usage:"description of config"`
	NumParallel int    `yaml:"parallel" name:"parallel" usage:"number of parallel" default:"1"`
	// musicXML_to_label
	// Suffix string `yaml:"suffix"`
	// Project settings
	Score         string `yaml:"score" name:"score" usage:"score file, required"`
	NumThreads    int    `yaml:"thread" name:"thread" usage:"number of parallel in session" default:"4"`
	InferenceMode int    `yaml:"inference" name:"inference" usage:"quality, processing speed: 2 (elements), 3 (standard) or 4 (advanced)" default:"3"`
	// NEUTRINO
	ModelDir   string `yaml:"model" name:"model" usage:"singer" default:"MERROW"`
	StyleShift int    `yaml:"styleShift" name:"styleShift" usage:"change the key and estimate to change the style of singing"`
	RandomSeed int    `yaml:"randomSeed" name:"randomSeed" usage:"random seed" default:"1234"`
	// NSF
	PitchShiftNsf float32 `yaml:"pitchShiftNsf" name:"pitchShiftNsf" usage:"change pitch via NSF"`
	// WORLD
	PitchShiftWorld    float32 `yaml:"pitchShiftWorld" name:"pitchShiftWorld" usage:"change pitch via WORLD"`
	FormantShift       float32 `yaml:"formantShift" name:"formantShift" usage:"change voice quality" default:"1"`
	SmoothPitch        float32 `yaml:"smoothPitch" name:"smoothPitch" usage:"[0, 100]%"`
	SmoothFormant      float32 `yaml:"smoothFormant" name:"smoothFormant" usage:"[0, 100]%"`
	EnhanceBreathiness float32 `yaml:"enhanceBreathiness" name:"enhanceBreathiness" usage:"[0, 100]%"`
}

func (c Config) basename() string {
	ss := strings.Split(filepath.Base(c.Score), ".")
	return strings.Join(ss[:len(ss)-1], ".")
}

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

func (c Config) Env() execx.Env {
	e := execx.NewEnv()
	for _, x := range []struct {
		k string
		v any
	}{
		{k: "NumThreads", v: c.NumThreads},
		{k: "InferenceMode", v: c.InferenceMode},
		{k: "ModelDir", v: c.ModelDir},
		{k: "StyleShift", v: c.StyleShift},
		{k: "PitchShiftNsf", v: c.PitchShiftNsf},
		{k: "PitchShiftWorld", v: c.PitchShiftWorld},
		{k: "FormantShift", v: c.FormantShift},
		{k: "SmoothPitch", v: c.SmoothPitch},
		{k: "SmoothFormant", v: c.SmoothFormant},
		{k: "EnhanceBreathiness", v: c.EnhanceBreathiness},
		{k: "NsfModel", v: c.nsfModel()},
		{k: "SamplingFreq", v: c.samplingFreq()},
		{k: "BASENAME", v: c.basename()},
		{k: "RandomSeed", v: c.RandomSeed},
		{k: "NumParallel", v: c.NumParallel},
	} {
		e.Set(x.k, fmt.Sprint(x.v))
	}
	return e
}
