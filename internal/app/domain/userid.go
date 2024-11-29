package domain

import (
	"github.com/google/uuid"
)

//go:generate easyjson -all
//easyjson:json
type UserID uuid.UUID

func (u UserID) String() string {
	return uuid.UUID(u).String()
}

func (u UserID) IsNil() bool {
	return uuid.UUID(u) == uuid.Nil
}

func NewUserID() UserID {
	return UserID(uuid.New())
}

func ParseUserID(id string) (UserID, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return UserID(uuid.Nil), err
	}

	return UserID(uid), nil
}
