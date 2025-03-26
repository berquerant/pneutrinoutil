package version

import (
	"fmt"
	"io"
)

var (
	Version  = "unknown"
	Revision = "unknown"
)

func Write(w io.Writer) {
	fmt.Fprintf(w, "Version: %s\n", Version)
	fmt.Fprintf(w, "Revision: %s\n", Revision)
}
