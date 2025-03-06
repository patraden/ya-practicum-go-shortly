package version

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
}

// NewVersion creates a new version instance.
func NewVersion() *Version {
	return &Version{
		BuildVersion: buildVersion,
		BuildDate:    buildDate,
		BuildCommit:  buildCommit,
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
