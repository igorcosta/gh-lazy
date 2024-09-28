package pkg // version.g

import (
	"fmt"
	"runtime"
)

// Version is the current version of the extension.
// This field is automatically updated by the CI/CD pipeline.
const Version = "v0.1.0"

// BuildDate is the date when the binary was built.
var BuildDate = "unknown"

// GitCommit is the git commit hash of the build.
var GitCommit = "unknown"

// GetVersionInfo returns a formatted string with full version information.
func GetVersionInfo() string {
	return fmt.Sprintf("Version: %s\nGit Commit: %s\nBuild Date: %s\nGo Version: %s\nOS/Arch: %s/%s",
		Version, GitCommit, BuildDate, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
