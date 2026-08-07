package version_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/berquerant/pneutrinoutil/pkg/version"
)

func TestWrite(t *testing.T) {
	var buf bytes.Buffer
	version.Write(&buf)
	output := buf.String()

	expectedPrefixes := []string{
		"Version:",
		"Revision:",
		"BuildDate:",
		"GoVersion:",
	}

	for _, prefix := range expectedPrefixes {
		if !strings.Contains(output, prefix) {
			t.Errorf("expected output to contain %q, got:\n%s", prefix, output)
		}
	}
}
