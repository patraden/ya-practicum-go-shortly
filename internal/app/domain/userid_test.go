package domain_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
)

func TestUserIDString(t *testing.T) {
	t.Parallel()

	uid := domain.UserID(uuid.New())
	expected := uuid.UUID(uid).String()

	assert.Equal(t, expected, uid.String())
}

func TestUserIDIsNil(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		userID   domain.UserID
		expected bool
	}{
		{"Non-nil UserID", domain.NewUserID(), false},
		{"Nil UserID", domain.UserID(uuid.Nil), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.expected, tt.userID.IsNil())
		})
	}
}

func TestNewUserID(t *testing.T) {
	t.Parallel()

	uid := domain.NewUserID()
	assert.False(t, uid.IsNil())
}

func TestParseUserID(t *testing.T) {
	t.Parallel()

	validUUID := uuid.New().String()

	t.Run("Valid UUID", func(t *testing.T) {
		t.Parallel()

		parsed, err := domain.ParseUserID(validUUID)
		require.NoError(t, err)
		assert.Equal(t, validUUID, parsed.String())
	})

	t.Run("Invalid UUID", func(t *testing.T) {
		t.Parallel()

		_, err := domain.ParseUserID("invalid-uuid")
		require.Error(t, err)
	})
}
