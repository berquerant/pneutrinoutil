package ctl

import (
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/berquerant/execx"
	"github.com/berquerant/pneutrinoutil/pkg/reflectx"
)

// NewDefaultConfig returns a new Config.
// Default values are from "default" struct tags.
func NewDefaultConfig() (*Config, error) {
	var c Config
	s, err := reflectx.NewFields(c)
	if err != nil {
		return nil, err
	}
	// get value of pointer to mutate
	v := reflect.ValueOf(&c)

	for _, f := range s {
		dv, ok := f.Tag.Lookup("default")
		if !ok {
			continue
		}
		newErr := func() error {
			return fmt.Errorf("unexpected default value of Config: %s, got %s", f.Name, dv)
		}
		fv := v.Elem().FieldByName(f.Name)

		// set default values from "default" tag
		switch f.Type.Kind() {
		case reflect.String:
			fv.SetString(dv)
		case reflect.Int:
			i, err := strconv.Atoi(dv)
			if err != nil {
				return nil, errors.Join(err, newErr())
			}
			fv.SetInt(int64(i))
		case reflect.Float32:
			f, err := strconv.ParseFloat(dv, 32)
			if err != nil {
				return nil, errors.Join(err, newErr())
			}
			fv.SetFloat(f)
		default:
			return nil, newErr()
		}
	}

	return &c, nil
}

type Config struct {
	Description string `yaml:"desc" usage:"description of config"`
	// musicXML_to_label
	// Suffix string `yaml:"suffix"`
	// Project settings
	Score         string `yaml:"score" usage:"score file, required"`
	NumThreads    int    `yaml:"thread" usage:"number of processor to use" default:"4"`
	InferenceMode int    `yaml:"inference" usage:"quality, processing speed: 2 (elements), 3 (standard) or 4 (advanced)" default:"3"`
	// NEUTRINO
	ModelDir   string `yaml:"model" usage:"singer" default:"MERROW"`
	StyleShift int    `yaml:"styleShift" usage:"change the key and estimate to change the style of singing"`
	RandomSeed int    `yaml:"randomSeed" usage:"random seed" default:"1234"`
	// NSF
	PitchShiftNsf float32 `yaml:"pitchShiftNsf" usage:"change pitch via NSF"`
	// WORLD
	PitchShiftWorld    float32 `yaml:"pitchShiftWorld" usage:"change pitch via WORLD"`
	FormantShift       float32 `yaml:"formantShift" usage:"change voice quality" default:"1"`
	SmoothPitch        float32 `yaml:"smoothPitch" usage:"[0, 100]%"`
	SmoothFormant      float32 `yaml:"smoothFormant" usage:"[0, 100]%"`
	EnhanceBreathiness float32 `yaml:"enhanceBreathiness" usage:"[0, 100]%"`
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
	} {
		e.Set(x.k, fmt.Sprint(x.v))
	}
	return e
}
