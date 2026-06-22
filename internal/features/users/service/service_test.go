package users_service_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/raizonw/todo-app/internal/core/domain"
	core_errors "github.com/raizonw/todo-app/internal/core/errors"
	users_service "github.com/raizonw/todo-app/internal/features/users/service"
)

type mockUsersRepository struct {
	createUserFunc func(ctx context.Context, user domain.User) (domain.User, error)
	getUsersFunc   func(ctx context.Context, limit *int, offset *int) ([]domain.User, error)
	getUserFunc    func(ctx context.Context, id int) (domain.User, error)
	deleteUserFunc func(ctx context.Context, id int) error
	patchUserFunc  func(ctx context.Context, id int, user domain.User) (domain.User, error)
}

func (m *mockUsersRepository) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	if m.createUserFunc != nil {
		return m.createUserFunc(ctx, user)
	}
	panic("mock method CreateUser not configured")
}

func (m *mockUsersRepository) GetUsers(ctx context.Context, limit *int, offset *int) ([]domain.User, error) {
	if m.getUsersFunc != nil {
		return m.getUsersFunc(ctx, limit, offset)
	}
	panic("mock method GetUsers not configured")
}

func (m *mockUsersRepository) GetUser(ctx context.Context, id int) (domain.User, error) {
	if m.getUserFunc != nil {
		return m.getUserFunc(ctx, id)
	}
	panic("mock method GetUser not configured")
}

func (m *mockUsersRepository) DeleteUser(ctx context.Context, id int) error {
	if m.deleteUserFunc != nil {
		return m.deleteUserFunc(ctx, id)
	}
	panic("mock method DeleteUser not configured")
}

func (m *mockUsersRepository) PatchUser(ctx context.Context, id int, user domain.User) (domain.User, error) {
	if m.patchUserFunc != nil {
		return m.patchUserFunc(ctx, id, user)
	}
	panic("mock method PatchUser not configured")
}

func TestCreateUser(t *testing.T) {
	type testCase struct {
		name         string
		inputUser    domain.User
		setupMock    func(m *mockUsersRepository)
		expectedUser domain.User
		expectedErr  error
	}

	validPhone := "+1234567890"

	tests := []testCase{
		{
			name: "Success - user created",
			inputUser: domain.User{
				FullName:    "Ivan Ivanov",
				PhoneNumber: &validPhone,
			},
			setupMock: func(m *mockUsersRepository) {
				m.createUserFunc = func(ctx context.Context, user domain.User) (domain.User, error) {
					user.ID = 1
					user.Version = 1
					return user, nil
				}
			},
			expectedUser: domain.User{
				ID:          1,
				Version:     1,
				FullName:    "Ivan Ivanov",
				PhoneNumber: &validPhone,
			},
			expectedErr: nil,
		},
		{
			name: "Validation Error - name too short",
			inputUser: domain.User{
				FullName: "Io",
			},
			setupMock:   func(m *mockUsersRepository) {},
			expectedErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "Validation Error - phone number too short",
			inputUser: domain.User{
				FullName:    "Ivan Ivanov",
				PhoneNumber: ptr("123"),
			},
			setupMock:   func(m *mockUsersRepository) {},
			expectedErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "Repository Error",
			inputUser: domain.User{
				FullName: "Ivan Ivanov",
			},
			setupMock: func(m *mockUsersRepository) {
				m.createUserFunc = func(ctx context.Context, user domain.User) (domain.User, error) {
					return domain.User{}, errors.New("db error")
				}
			},
			expectedErr: errors.New("create user: db error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockRepo := &mockUsersRepository{}
			tc.setupMock(mockRepo)

			service := users_service.NewUsersService(mockRepo)
			res, err := service.CreateUser(context.Background(), tc.inputUser)

			if tc.expectedErr != nil {
				if err == nil {
					t.Fatalf("expected error containing '%v', got nil", tc.expectedErr)
				}
				if !strings.Contains(err.Error(), tc.expectedErr.Error()) {
					t.Errorf("expected error containing '%v', got '%v'", tc.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if res.ID != tc.expectedUser.ID || res.FullName != tc.expectedUser.FullName {
				t.Errorf("expected user %+v, got %+v", tc.expectedUser, res)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	type testCase struct {
		name         string
		id           int
		setupMock    func(m *mockUsersRepository)
		expectedUser domain.User
		expectedErr  error
	}

	tests := []testCase{
		{
			name: "Success",
			id:   1,
			setupMock: func(m *mockUsersRepository) {
				m.getUserFunc = func(ctx context.Context, id int) (domain.User, error) {
					return domain.User{ID: 1, FullName: "Ivan Ivanov"}, nil
				}
			},
			expectedUser: domain.User{ID: 1, FullName: "Ivan Ivanov"},
			expectedErr:  nil,
		},
		{
			name: "Not Found",
			id:   999,
			setupMock: func(m *mockUsersRepository) {
				m.getUserFunc = func(ctx context.Context, id int) (domain.User, error) {
					return domain.User{}, core_errors.ErrNotFound
				}
			},
			expectedErr: core_errors.ErrNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockRepo := &mockUsersRepository{}
			tc.setupMock(mockRepo)

			service := users_service.NewUsersService(mockRepo)
			res, err := service.GetUser(context.Background(), tc.id)

			if tc.expectedErr != nil {
				if err == nil {
					t.Fatalf("expected error containing '%v', got nil", tc.expectedErr)
				}
				if !strings.Contains(err.Error(), tc.expectedErr.Error()) {
					t.Errorf("expected error containing '%v', got '%v'", tc.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if res.ID != tc.expectedUser.ID || res.FullName != tc.expectedUser.FullName {
				t.Errorf("expected user %+v, got %+v", tc.expectedUser, res)
			}
		})
	}
}

func TestGetUsers(t *testing.T) {
	type testCase struct {
		name          string
		limit         *int
		offset        *int
		setupMock     func(m *mockUsersRepository)
		expectedUsers []domain.User
		expectedErr   error
	}

	tests := []testCase{
		{
			name:   "Success - no pagination params",
			limit:  nil,
			offset: nil,
			setupMock: func(m *mockUsersRepository) {
				m.getUsersFunc = func(ctx context.Context, limit, offset *int) ([]domain.User, error) {
					return []domain.User{{ID: 1, FullName: "User 1"}, {ID: 2, FullName: "User 2"}}, nil
				}
			},
			expectedUsers: []domain.User{{ID: 1, FullName: "User 1"}, {ID: 2, FullName: "User 2"}},
			expectedErr:   nil,
		},
		{
			name:        "Error - negative limit",
			limit:       ptr(-1),
			offset:      nil,
			setupMock:   func(m *mockUsersRepository) {},
			expectedErr: core_errors.ErrInvalidArgument,
		},
		{
			name:        "Error - negative offset",
			limit:       nil,
			offset:      ptr(-5),
			setupMock:   func(m *mockUsersRepository) {},
			expectedErr: core_errors.ErrInvalidArgument,
		},
		{
			name:   "Repository Error",
			limit:  nil,
			offset: nil,
			setupMock: func(m *mockUsersRepository) {
				m.getUsersFunc = func(ctx context.Context, limit, offset *int) ([]domain.User, error) {
					return nil, errors.New("db error")
				}
			},
			expectedErr: errors.New("get users from repository: db error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockRepo := &mockUsersRepository{}
			tc.setupMock(mockRepo)

			service := users_service.NewUsersService(mockRepo)
			res, err := service.GetUsers(context.Background(), tc.limit, tc.offset)

			if tc.expectedErr != nil {
				if err == nil {
					t.Fatalf("expected error containing '%v', got nil", tc.expectedErr)
				}
				if !strings.Contains(err.Error(), tc.expectedErr.Error()) {
					t.Errorf("expected error containing '%v', got '%v'", tc.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(res) != len(tc.expectedUsers) {
				t.Fatalf("expected %d users, got %d", len(tc.expectedUsers), len(res))
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	type testCase struct {
		name        string
		id          int
		setupMock   func(m *mockUsersRepository)
		expectedErr error
	}

	tests := []testCase{
		{
			name: "Success",
			id:   1,
			setupMock: func(m *mockUsersRepository) {
				m.deleteUserFunc = func(ctx context.Context, id int) error {
					return nil
				}
			},
			expectedErr: nil,
		},
		{
			name: "Repository Error",
			id:   1,
			setupMock: func(m *mockUsersRepository) {
				m.deleteUserFunc = func(ctx context.Context, id int) error {
					return errors.New("db error")
				}
			},
			expectedErr: errors.New("delete user: db error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockRepo := &mockUsersRepository{}
			tc.setupMock(mockRepo)

			service := users_service.NewUsersService(mockRepo)
			err := service.DeleteUser(context.Background(), tc.id)

			if tc.expectedErr != nil {
				if err == nil {
					t.Fatalf("expected error containing '%v', got nil", tc.expectedErr)
				}
				if !strings.Contains(err.Error(), tc.expectedErr.Error()) {
					t.Errorf("expected error containing '%v', got '%v'", tc.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestPatchUser(t *testing.T) {
	type testCase struct {
		name         string
		id           int
		patch        domain.UserPatch
		setupMock    func(m *mockUsersRepository)
		expectedUser domain.User
		expectedErr  error
	}

	tests := []testCase{
		{
			name: "Success - update full name",
			id:   1,
			patch: domain.UserPatch{
				FullName: domain.Nullable[string]{Set: true, Value: ptr("Ivan P. Ivanov")},
			},
			setupMock: func(m *mockUsersRepository) {
				m.getUserFunc = func(ctx context.Context, id int) (domain.User, error) {
					return domain.User{ID: 1, FullName: "Ivan Ivanov", Version: 1}, nil
				}
				m.patchUserFunc = func(ctx context.Context, id int, user domain.User) (domain.User, error) {
					return user, nil
				}
			},
			expectedUser: domain.User{ID: 1, FullName: "Ivan P. Ivanov", Version: 1},
			expectedErr:  nil,
		},
		{
			name: "Error - user not found in repo",
			id:   999,
			patch: domain.UserPatch{
				FullName: domain.Nullable[string]{Set: true, Value: ptr("Ivan P. Ivanov")},
			},
			setupMock: func(m *mockUsersRepository) {
				m.getUserFunc = func(ctx context.Context, id int) (domain.User, error) {
					return domain.User{}, core_errors.ErrNotFound
				}
			},
			expectedErr: core_errors.ErrNotFound,
		},
		{
			name: "Error - apply patch validation fails (name too short)",
			id:   1,
			patch: domain.UserPatch{
				FullName: domain.Nullable[string]{Set: true, Value: ptr("Io")},
			},
			setupMock: func(m *mockUsersRepository) {
				m.getUserFunc = func(ctx context.Context, id int) (domain.User, error) {
					return domain.User{ID: 1, FullName: "Ivan Ivanov", Version: 1}, nil
				}
			},
			expectedErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "Error - repository patch error",
			id:   1,
			patch: domain.UserPatch{
				FullName: domain.Nullable[string]{Set: true, Value: ptr("Ivan P. Ivanov")},
			},
			setupMock: func(m *mockUsersRepository) {
				m.getUserFunc = func(ctx context.Context, id int) (domain.User, error) {
					return domain.User{ID: 1, FullName: "Ivan Ivanov", Version: 1}, nil
				}
				m.patchUserFunc = func(ctx context.Context, id int, user domain.User) (domain.User, error) {
					return domain.User{}, errors.New("db patch error")
				}
			},
			expectedErr: errors.New("patch user: db patch error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockRepo := &mockUsersRepository{}
			tc.setupMock(mockRepo)

			service := users_service.NewUsersService(mockRepo)
			res, err := service.PatchUser(context.Background(), tc.id, tc.patch)

			if tc.expectedErr != nil {
				if err == nil {
					t.Fatalf("expected error containing '%v', got nil", tc.expectedErr)
				}
				if !strings.Contains(err.Error(), tc.expectedErr.Error()) {
					t.Errorf("expected error containing '%v', got '%v'", tc.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if res.FullName != tc.expectedUser.FullName {
				t.Errorf("expected user full name '%s', got '%s'", tc.expectedUser.FullName, res.FullName)
			}
		})
	}
}

func ptr[T any](v T) *T {
	return &v
}
