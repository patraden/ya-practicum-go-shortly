package version

import (
	"github.com/rs/zerolog"
)

var (
	// buildVersion is the application build version, example: v1.0.0.
	buildVersion = "N/A"

	// buildDate is the application build date, example: 01.02.2025.
	buildDate = "N/A"

	// buildCommit is the application build commit, example: abcd1234.
	buildCommit = "N/A"
)

// Version is provides application build version details.
type Version struct {
	BuildVersion string
	BuildDate    string
	BuildCommit  string
	log          *zerolog.Logger
}

// NewVersion creates a new version instance.
func NewVersion(log *zerolog.Logger) *Version {
	return &Version{
		BuildVersion: buildVersion,
		BuildDate:    buildDate,
		BuildCommit:  buildCommit,
		log:          log,
	}
}

// BuildVersion returns the application build version.
func (v *Version) Version() string {
	return v.BuildVersion
}

// BuildDate returns the application build date.
func (v *Version) Date() string {
	return v.BuildDate
}

// BuildCommit returns the application build commit.
func (v *Version) Commit() string {
	return v.BuildCommit
}

// Log prints build info to stdout.
func (v *Version) Log() {
	v.log.Info().Msgf("Build version: %s", v.Version())
	v.log.Info().Msgf("Build date: %s", v.Date())
	v.log.Info().Msgf("Build commit: %s", v.Commit())
}
