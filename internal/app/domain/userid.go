package domain

import (
	"github.com/google/uuid"
)

//go:generate easyjson -all
//easyjson:json

// UserID represents a unique identifier for a user.
type UserID uuid.UUID

// String returns the string representation of the UserID.
func (u UserID) String() string {
	return uuid.UUID(u).String()
}

// IsNil checks whether the UserID is nil (zero value).
func (u UserID) IsNil() bool {
	return uuid.UUID(u) == uuid.Nil
}

// NewUserID generates and returns a new unique UserID.
func NewUserID() UserID {
	return UserID(uuid.New())
}

// ParseUserID converts a string representation of a UUID into a UserID.
func ParseUserID(id string) (UserID, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return UserID(uuid.Nil), err
	}

	return UserID(uid), nil
}
