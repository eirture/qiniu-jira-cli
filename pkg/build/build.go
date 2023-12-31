package build

import "runtime"

var (
	Version   = "unknown"
	Commit    = "unknown"
	GoVersion = runtime.Version()
)
