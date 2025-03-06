package version_test

import (
	"testing"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/version"
)

func TestVersionMethods(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		version     *version.Version
		wantVersion string
		wantDate    string
		wantCommit  string
	}{
		{
			name:        "Default Version",
			version:     version.NewVersion(),
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.version.Version(); got != tt.wantVersion {
				t.Errorf("Version() = %q, want %q", got, tt.wantVersion)
			}

			if got := tt.version.Date(); got != tt.wantDate {
				t.Errorf("Date() = %q, want %q", got, tt.wantDate)
			}

			if got := tt.version.Commit(); got != tt.wantCommit {
				t.Errorf("Commit() = %q, want %q", got, tt.wantCommit)
			}
		})
	}
}
