package version

import (
	"fmt"
	"runtime"
)

var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

// GetVersionInfo returns a formatted string with full version information.
func GetVersionInfo() string {
	return fmt.Sprintf("Version: %s\nGit Commit: %s\nBuild Date: %s\nGo Version: %s\nOS/Arch: %s/%s",
		Version, Commit, BuildDate, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
