package tasks_service_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/raizonw/todo-app/internal/core/domain"
	core_errors "github.com/raizonw/todo-app/internal/core/errors"
	tasks_service "github.com/raizonw/todo-app/internal/features/tasks/service"
)

type mockTasksRepository struct {
	createTaskFunc func(ctx context.Context, task domain.Task) (domain.Task, error)
	getTasksFunc   func(ctx context.Context, userID *int, limit *int, offset *int) ([]domain.Task, error)
	getTaskFunc    func(ctx context.Context, id int) (domain.Task, error)
	deleteTaskFunc func(ctx context.Context, id int) error
	patchTaskFunc  func(ctx context.Context, id int, task domain.Task) (domain.Task, error)
}

func (m *mockTasksRepository) CreateTask(ctx context.Context, task domain.Task) (domain.Task, error) {
	if m.createTaskFunc != nil {
		return m.createTaskFunc(ctx, task)
	}
	panic("mock method CreateTask not configured")
}

func (m *mockTasksRepository) GetTasks(ctx context.Context, userID *int, limit *int, offset *int) ([]domain.Task, error) {
	if m.getTasksFunc != nil {
		return m.getTasksFunc(ctx, userID, limit, offset)
	}
	panic("mock method GetTasks not configured")
}

func (m *mockTasksRepository) GetTask(ctx context.Context, id int) (domain.Task, error) {
	if m.getTaskFunc != nil {
		return m.getTaskFunc(ctx, id)
	}
	panic("mock method GetTask not configured")
}

func (m *mockTasksRepository) DeleteTask(ctx context.Context, id int) error {
	if m.deleteTaskFunc != nil {
		return m.deleteTaskFunc(ctx, id)
	}
	panic("mock method DeleteTask not configured")
}

func (m *mockTasksRepository) PatchTask(ctx context.Context, id int, task domain.Task) (domain.Task, error) {
	if m.patchTaskFunc != nil {
		return m.patchTaskFunc(ctx, id, task)
	}
	panic("mock method PatchTask not configured")
}

func TestCreateTask(t *testing.T) {
	type testCase struct {
		name         string
		inputTask    domain.Task
		setupMock    func(m *mockTasksRepository)
		expectedTask domain.Task
		expectedErr  error
	}

	tests := []testCase{
		{
			name: "Success - task created",
			inputTask: domain.Task{
				Title:        "Buy groceries",
				AuthorUserID: 1,
			},
			setupMock: func(m *mockTasksRepository) {
				m.createTaskFunc = func(ctx context.Context, task domain.Task) (domain.Task, error) {
					task.ID = 1
					task.Version = 1
					task.CreatedAt = time.Now()
					return task, nil
				}
			},
			expectedTask: domain.Task{
				ID:           1,
				Version:      1,
				Title:        "Buy groceries",
				AuthorUserID: 1,
			},
			expectedErr: nil,
		},
		{
			name: "Validation Error - title is empty",
			inputTask: domain.Task{
				Title:        "",
				AuthorUserID: 1,
			},
			setupMock:   func(m *mockTasksRepository) {},
			expectedErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "Validation Error - completed without CompletedAt",
			inputTask: domain.Task{
				Title:        "Do laundry",
				Completed:    true,
				CompletedAt:  nil,
				AuthorUserID: 1,
			},
			setupMock:   func(m *mockTasksRepository) {},
			expectedErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "Repository Error",
			inputTask: domain.Task{
				Title:        "Buy groceries",
				AuthorUserID: 1,
			},
			setupMock: func(m *mockTasksRepository) {
				m.createTaskFunc = func(ctx context.Context, task domain.Task) (domain.Task, error) {
					return domain.Task{}, errors.New("db error")
				}
			},
			expectedErr: errors.New("create task: db error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockRepo := &mockTasksRepository{}
			tc.setupMock(mockRepo)

			service := tasks_service.NewTasksService(mockRepo)
			res, err := service.CreateTask(context.Background(), tc.inputTask)

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

			if res.ID != tc.expectedTask.ID || res.Title != tc.expectedTask.Title {
				t.Errorf("expected task %+v, got %+v", tc.expectedTask, res)
			}
		})
	}
}

func TestGetTask(t *testing.T) {
	type testCase struct {
		name         string
		id           int
		setupMock    func(m *mockTasksRepository)
		expectedTask domain.Task
		expectedErr  error
	}

	tests := []testCase{
		{
			name: "Success",
			id:   1,
			setupMock: func(m *mockTasksRepository) {
				m.getTaskFunc = func(ctx context.Context, id int) (domain.Task, error) {
					return domain.Task{ID: 1, Title: "Buy groceries"}, nil
				}
			},
			expectedTask: domain.Task{ID: 1, Title: "Buy groceries"},
			expectedErr:  nil,
		},
		{
			name: "Not Found",
			id:   999,
			setupMock: func(m *mockTasksRepository) {
				m.getTaskFunc = func(ctx context.Context, id int) (domain.Task, error) {
					return domain.Task{}, core_errors.ErrNotFound
				}
			},
			expectedErr: core_errors.ErrNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockRepo := &mockTasksRepository{}
			tc.setupMock(mockRepo)

			service := tasks_service.NewTasksService(mockRepo)
			res, err := service.GetTask(context.Background(), tc.id)

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

			if res.ID != tc.expectedTask.ID || res.Title != tc.expectedTask.Title {
				t.Errorf("expected task %+v, got %+v", tc.expectedTask, res)
			}
		})
	}
}

func TestGetTasks(t *testing.T) {
	type testCase struct {
		name          string
		userID        *int
		limit         *int
		offset        *int
		setupMock     func(m *mockTasksRepository)
		expectedTasks []domain.Task
		expectedErr   error
	}

	tests := []testCase{
		{
			name:   "Success - all tasks",
			userID: nil,
			limit:  nil,
			offset: nil,
			setupMock: func(m *mockTasksRepository) {
				m.getTasksFunc = func(ctx context.Context, u, l, o *int) ([]domain.Task, error) {
					return []domain.Task{{ID: 1, Title: "Task 1"}, {ID: 2, Title: "Task 2"}}, nil
				}
			},
			expectedTasks: []domain.Task{{ID: 1, Title: "Task 1"}, {ID: 2, Title: "Task 2"}},
			expectedErr:   nil,
		},
		{
			name:        "Error - negative limit",
			userID:      nil,
			limit:       ptr(-1),
			offset:      nil,
			setupMock:   func(m *mockTasksRepository) {},
			expectedErr: core_errors.ErrInvalidArgument,
		},
		{
			name:        "Error - negative offset",
			userID:      nil,
			limit:       nil,
			offset:      ptr(-5),
			setupMock:   func(m *mockTasksRepository) {},
			expectedErr: core_errors.ErrInvalidArgument,
		},
		{
			name:   "Repository Error",
			userID: nil,
			limit:  nil,
			offset: nil,
			setupMock: func(m *mockTasksRepository) {
				m.getTasksFunc = func(ctx context.Context, u, l, o *int) ([]domain.Task, error) {
					return nil, errors.New("db error")
				}
			},
			expectedErr: errors.New("get tasks from repository: db error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockRepo := &mockTasksRepository{}
			tc.setupMock(mockRepo)

			service := tasks_service.NewTasksService(mockRepo)
			res, err := service.GetTasks(context.Background(), tc.userID, tc.limit, tc.offset)

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

			if len(res) != len(tc.expectedTasks) {
				t.Fatalf("expected %d tasks, got %d", len(tc.expectedTasks), len(res))
			}
		})
	}
}

func TestDeleteTask(t *testing.T) {
	type testCase struct {
		name        string
		id          int
		setupMock   func(m *mockTasksRepository)
		expectedErr error
	}

	tests := []testCase{
		{
			name: "Success",
			id:   1,
			setupMock: func(m *mockTasksRepository) {
				m.deleteTaskFunc = func(ctx context.Context, id int) error {
					return nil
				}
			},
			expectedErr: nil,
		},
		{
			name: "Repository Error",
			id:   1,
			setupMock: func(m *mockTasksRepository) {
				m.deleteTaskFunc = func(ctx context.Context, id int) error {
					return errors.New("db error")
				}
			},
			expectedErr: errors.New("delete task from repository: db error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockRepo := &mockTasksRepository{}
			tc.setupMock(mockRepo)

			service := tasks_service.NewTasksService(mockRepo)
			err := service.DeleteTask(context.Background(), tc.id)

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

func TestPatchTask(t *testing.T) {
	type testCase struct {
		name         string
		id           int
		patch        domain.TaskPatch
		setupMock    func(m *mockTasksRepository)
		expectedTask domain.Task
		expectedErr  error
	}

	tests := []testCase{
		{
			name: "Success - update title",
			id:   1,
			patch: domain.TaskPatch{
				Title: domain.Nullable[string]{Set: true, Value: ptr("New Title")},
			},
			setupMock: func(m *mockTasksRepository) {
				m.getTaskFunc = func(ctx context.Context, id int) (domain.Task, error) {
					return domain.Task{ID: 1, Title: "Old Title", Version: 1}, nil
				}
				m.patchTaskFunc = func(ctx context.Context, id int, task domain.Task) (domain.Task, error) {
					return task, nil
				}
			},
			expectedTask: domain.Task{ID: 1, Title: "New Title", Version: 1},
			expectedErr:  nil,
		},
		{
			name: "Error - task not found in repository",
			id:   999,
			patch: domain.TaskPatch{
				Title: domain.Nullable[string]{Set: true, Value: ptr("New Title")},
			},
			setupMock: func(m *mockTasksRepository) {
				m.getTaskFunc = func(ctx context.Context, id int) (domain.Task, error) {
					return domain.Task{}, core_errors.ErrNotFound
				}
			},
			expectedErr: core_errors.ErrNotFound,
		},
		{
			name: "Error - apply patch validation fails (title too short)",
			id:   1,
			patch: domain.TaskPatch{
				Title: domain.Nullable[string]{Set: true, Value: ptr("")},
			},
			setupMock: func(m *mockTasksRepository) {
				m.getTaskFunc = func(ctx context.Context, id int) (domain.Task, error) {
					return domain.Task{ID: 1, Title: "Old Title", Version: 1}, nil
				}
			},
			expectedErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "Error - repository patch error",
			id:   1,
			patch: domain.TaskPatch{
				Title: domain.Nullable[string]{Set: true, Value: ptr("New Title")},
			},
			setupMock: func(m *mockTasksRepository) {
				m.getTaskFunc = func(ctx context.Context, id int) (domain.Task, error) {
					return domain.Task{ID: 1, Title: "Old Title", Version: 1}, nil
				}
				m.patchTaskFunc = func(ctx context.Context, id int, task domain.Task) (domain.Task, error) {
					return domain.Task{}, errors.New("db patch error")
				}
			},
			expectedErr: errors.New("patch task: db patch error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockRepo := &mockTasksRepository{}
			tc.setupMock(mockRepo)

			service := tasks_service.NewTasksService(mockRepo)
			res, err := service.PatchTask(context.Background(), tc.id, tc.patch)

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

			if res.Title != tc.expectedTask.Title {
				t.Errorf("expected task title '%s', got '%s'", tc.expectedTask.Title, res.Title)
			}
		})
	}
}

func ptr[T any](v T) *T {
	return &v
}
