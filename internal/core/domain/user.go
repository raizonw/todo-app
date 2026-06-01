package domain

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
