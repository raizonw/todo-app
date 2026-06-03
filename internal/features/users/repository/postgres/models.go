package users_postgres_repository

import "github.com/raizonw/todo-app/internal/core/domain"

type UserModel struct {
	ID          int
	Version     int
	FullName    string
	PhoneNumber *string
}

func userDomainsFromModels(users []UserModel) []domain.User {
	userDomains := make([]domain.User, len(users))

	for i, user := range users {
		userDomains[i] = domain.User{
			ID:          user.ID,
			Version:     user.Version,
			FullName:    user.FullName,
			PhoneNumber: user.PhoneNumber,
		}
	}

	return userDomains
}
