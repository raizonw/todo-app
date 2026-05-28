package users_transport_http

type UsersService interface {
}

func NewUsersHTTPHandler(
	userService UsersService,
) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		usersService: userService,
	}
}

type UsersHTTPHandler struct {
	usersService UsersService
}
