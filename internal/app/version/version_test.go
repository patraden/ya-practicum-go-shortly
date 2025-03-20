package version_test

import (
	"testing"

	"github.com/rs/zerolog"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/logger"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/version"
)

func TestVersionMethods(t *testing.T) {
	t.Parallel()

	log := logger.NewLogger(zerolog.Disabled).GetLogger()

	tests := []struct {
		name        string
		version     *version.Version
		wantVersion string
		wantDate    string
		wantCommit  string
	}{
		{
			name:        "Default Version",
			version:     version.NewVersion(log),
			wantVersion: "N/A",
			wantDate:    "N/A",
			wantCommit:  "N/A",
		},
		{
			name: "Custom Version",
			version: &version.Version{
				BuildVersion: "v1.2.3",
				BuildDate:    "02.03.2025",
				BuildCommit:  "abcd1234",
			},
			wantVersion: "v1.2.3",
			wantDate:    "02.03.2025",
			wantCommit:  "abcd1234",
		},
	}

	for _, ttest := range tests {
		t.Run(ttest.name, func(t *testing.T) {
			t.Parallel()

			if got := ttest.version.Version(); got != ttest.wantVersion {
				t.Errorf("Version() = %q, want %q", got, ttest.wantVersion)
			}

			if got := ttest.version.Date(); got != ttest.wantDate {
				t.Errorf("Date() = %q, want %q", got, ttest.wantDate)
			}

			if got := ttest.version.Commit(); got != ttest.wantCommit {
				t.Errorf("Commit() = %q, want %q", got, ttest.wantCommit)
			}

			if ttest.name == "Default Version" {
				ttest.version.Log()
			}
		})
	}
}
