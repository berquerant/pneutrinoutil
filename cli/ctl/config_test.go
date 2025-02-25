package ctl_test

import (
	"testing"

	"github.com/berquerant/pneutrinoutil/cli/ctl"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		got, err := ctl.NewDefaultConfig()
		if !assert.Nil(t, err) {
			t.Logf("%v", err)
			return
		}
		want := &ctl.Config{
			NumParallel:   1,
			NumThreads:    4,
			InferenceMode: 3,
			ModelDir:      "MERROW",
			RandomSeed:    1234,
			FormantShift:  1,
		}
		assert.Equal(t, want, got)
	})
}
