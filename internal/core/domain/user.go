package domain

import (
	"fmt"

	core_errors "github.com/raizonw/todo-app/internal/core/errors"
)

type User struct {
	ID          int
	Version     int
	FullName    string
	PhoneNumber *string
}

func NewUser(
	id int,
	version int,
	fullName string,
	phoneNumber *string,
) User {
	return User{
		ID:          id,
		Version:     version,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}
}

func NewUserUnitialized(fullName string,
	phoneNumber *string,
) User {
	return NewUser(
		UninitializedId,
		UnitializedVersion,
		fullName,
		phoneNumber,
	)
}

func (u *User) Validate() error {
	FullNameLength := len([]rune(u.FullName))
	if FullNameLength < 3 || FullNameLength > 100 {
		return fmt.Errorf("invalid `FullName` len: %d: %w",
			FullNameLength,
			core_errors.ErrInvalidArgument)
	}

	if u.PhoneNumber != nil {
		phoneNumberLentgh := len([]rune(*u.PhoneNumber))
		if phoneNumberLentgh < 10 || phoneNumberLentgh > 15 {
			return fmt.Errorf(
				"invalid `PhoneNumber` len: %d: %w",
				phoneNumberLentgh,
				core_errors.ErrInvalidArgument,
			)
		}
	}

	return nil
}
