package version

import (
	"fmt"
	"io"
	"runtime"
)

var (
	Version   = "unknown"
	Revision  = "unknown"
	BuildDate = "unknown"
	GoVersion = runtime.Version()
)

func Write(w io.Writer) {
	_, _ = fmt.Fprintf(w, "Version: %s\n", Version)
	_, _ = fmt.Fprintf(w, "Revision: %s\n", Revision)
	_, _ = fmt.Fprintf(w, "BuildDate: %s\n", BuildDate)
	_, _ = fmt.Fprintf(w, "GoVersion: %s\n", GoVersion)
}
