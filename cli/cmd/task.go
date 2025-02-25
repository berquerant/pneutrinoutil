package cmd

import (
	"slices"

	"github.com/berquerant/execx"
	"github.com/berquerant/pneutrinoutil/pkg/set"
)

func prepareTaskEntrypoint(taskNames, include, exclude []string) []string {
	includeSet := func() set.Set[string] {
		if len(include) > 0 {
			return set.New(include)
		}
		return set.New(taskNames)
	}()
	selectedTaskSet := includeSet.Diff(set.New(exclude))

	r := []string{"set -ex"}
	for _, name := range taskNames {
		if selectedTaskSet.In(name) {
			r = append(r, name)
		}
	}
	return r
}

func prepareAdditionalEnviron(whiteList []string) execx.Env {
	environ := execx.EnvFromEnviron()
	if slices.Contains(whiteList, "all") {
		return environ
	}

	r := execx.NewEnv()
	for _, key := range whiteList {
		if value, ok := environ.Get(key); ok {
			r.Set(key, value)
		}
	}
	return r
}
