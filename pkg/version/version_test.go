package version_test

import (
	"bytes"
	"testing"

	"github.com/berquerant/pneutrinoutil/pkg/version"
	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {
	var buf bytes.Buffer
	version.Write(&buf)
	output := buf.String()

	for _, tc := range []struct {
		prefix string
	}{
		{prefix: "Version:"},
		{prefix: "Revision:"},
		{prefix: "BuildDate:"},
		{prefix: "GoVersion:"},
	} {
		t.Run(tc.prefix, func(t *testing.T) {
			assert.Contains(t, output, tc.prefix)
		})
	}
}
