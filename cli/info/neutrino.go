package info

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var ErrNeutrinoVersion = errors.New("NeutrinoVersion")

func GetNeutrinoVersion(ctx context.Context, dir string) (string, error) {
	x, err := filepath.Abs(dir)
	if err != nil {
		return "", errors.Join(ErrNeutrinoVersion, err)
	}
	bin := filepath.Join(x, "bin", "NEUTRINO")
	if _, err := os.Stat(bin); err != nil {
		bin = filepath.Join(x, "bin", "neutrino")
	}
	dyld := filepath.Join(x, "bin")

	cmd := exec.CommandContext(ctx, bin)
	cmd.Env = append(cmd.Env, "DYLD_LIBRARY_PATH="+dyld)
	output, _ := cmd.Output()
	ss := bytes.SplitN(output, []byte("\n"), 2)
	if len(ss) != 2 {
		return "", fmt.Errorf("%w: failed to get output", ErrNeutrinoVersion)
	}
	s := string(ss[0])
	vs := strings.SplitN(s, "-", 2)
	if len(vs) != 2 {
		return "", fmt.Errorf("%w: failed to get output", ErrNeutrinoVersion)
	}
	return strings.TrimSpace(vs[1]), nil
}
