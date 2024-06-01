package version

import (
	"bytes"
	"time"
)

// Version is the version info for pig-hole.
type Version struct {
	Version    string
	Build      string
	BuildDate  string
	CommitHash string
}

var (
	// Version is Major, Minor, Patch, etc are set at build time
	version    string
	build      string
	buildDate  string
	commitHash string

	// AppVersion is pig-hole version.
	AppVersion Version
)

func init() {
	if version == "" {
		AppVersion.Version = "0.0.0"
	} else {
		AppVersion.Version = version
	}

	if build == "" {
		AppVersion.Build = "dev"
	} else {
		AppVersion.Build = build
	}

	if buildDate == "" {
		AppVersion.BuildDate = time.Now().Format(time.RFC822)
	} else {
		AppVersion.BuildDate = buildDate
	}

	if commitHash == "" {
		AppVersion.CommitHash = "none"
	} else {
		AppVersion.CommitHash = commitHash
	}
}

func (v Version) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(v.Version)
	buffer.WriteString("-" + v.Build)
	return buffer.String()
}
