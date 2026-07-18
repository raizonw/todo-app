package tasks_transport_http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/raizonw/todo-app/internal/core/domain"
)

type fakeService struct {
	createTask func(ctx context.Context, task domain.Task) (domain.Task, error)
	getTasks   func(ctx context.Context, userID *int, limit *int, offset *int) ([]domain.Task, error)
	getTask    func(ctx context.Context, id int) (domain.Task, error)
	deleteTask func(ctx context.Context, id int) error
	patchTask  func(ctx context.Context, id int, patch domain.TaskPatch) (domain.Task, error)
}

func (f fakeService) CreateTask(ctx context.Context, task domain.Task) (domain.Task, error) {
	return f.createTask(ctx, task)
}

func (f fakeService) DeleteTask(ctx context.Context, id int) error {
	return f.deleteTask(ctx, id)
}

func (f fakeService) GetTask(ctx context.Context, id int) (domain.Task, error) {
	return f.getTask(ctx, id)
}

func (f fakeService) GetTasks(ctx context.Context, userID *int, limit *int, offset *int) ([]domain.Task, error) {
	return f.getTasks(ctx, userID, limit, offset)
}

func (f fakeService) PatchTask(ctx context.Context, id int, patch domain.TaskPatch) (domain.Task, error) {
	return f.patchTask(ctx, id, patch)
}

func ptr[T any](v T) *T {
	return &v
}

func TestHandler(t *testing.T) {
	t.Run("Successs", func(t *testing.T) {
		createdAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

		body := strings.NewReader(`{
			"title": "Купить молоко",
			"description": "2 литра",
			"author_user_id": 1
		}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", body)
		rec := httptest.NewRecorder()
		task := domain.NewTask(1, 1, "Купить молоко", ptr("2 литра"), false, createdAt, nil, 1)
		service := fakeService{
			createTask: func(ctx context.Context, task domain.Task) (domain.Task, error) {
				return task, nil
			},
		}
		handler := NewTasksHTTPHandler(service)

	})
}
