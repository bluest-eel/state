package common

import (
	"time"

	"github.com/bluest-eel/common/util"
)

// Versioning data
var (
	version    string
	buildDate  string
	gitCommit  string
	gitBranch  string
	gitSummary string
)

// VersionData stuff for things
func VersionData() *util.Version {
	buildTime := time.Now()
	return &util.Version{
		Semantic:   version,
		BuildDate:  buildDate,
		BuildTime:  &buildTime,
		GitCommit:  gitCommit,
		GitBranch:  gitBranch,
		GitSummary: gitSummary,
	}
}

// BuildString ...
func BuildString() string {
	return util.BuildString(VersionData())
}

// VersionString ...
func VersionString() string {
	return util.VersionString(VersionData())
}

// VersionedBuildString ...
func VersionedBuildString() string {
	return util.VersionedBuildString(VersionData())
}
