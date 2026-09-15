package tasks_transport_http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/raizonw/todo-app/internal/core/domain"
	core_logger "github.com/raizonw/todo-app/internal/core/logger"
	"go.uber.org/zap"
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

func TestHandler(t *testing.T) {
	log := &core_logger.Logger{
		Logger: zap.NewNop(),
	}
	t.Run("successsful create", func(t *testing.T) {

		called := false
		var receivedTask domain.Task

		body := strings.NewReader(`{
			"title": "Купить молоко",
			"description": "2 литра",
			"author_user_id": 1
		}`)
		title, description := "Купить молоко", "2 литра"

		req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", body)

		req = req.WithContext(core_logger.ToContext(req.Context(), log))
		rec := httptest.NewRecorder()
		service := fakeService{
			createTask: func(ctx context.Context, task domain.Task) (domain.Task, error) {
				called = true
				receivedTask = task
				return task, nil
			},
		}
		handler := NewTasksHTTPHandler(service)

		handler.CreateTask(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusCreated)
		}

		if !called {
			t.Error("service dont get called")
		}

		if receivedTask.Title != title {
			t.Errorf("title = %q, want %q", receivedTask.Title, title)
		}

		if receivedTask.AuthorUserID != 1 {
			t.Errorf("author user id = %d, want 1", receivedTask.AuthorUserID)
		}

		if receivedTask.Description == nil || *receivedTask.Description != description {
			t.Errorf("unexpected description: %v", receivedTask.Description)
		}

	})

	t.Run("invalid JSON", func(t *testing.T) {
		body := strings.NewReader(`{"title":}`)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", body)
		req = req.WithContext(core_logger.ToContext(req.Context(), log))
		res := httptest.NewRecorder()

		service := fakeService{
			createTask: func(ctx context.Context, task domain.Task) (domain.Task, error) {
				t.Fatal("sevice must not be called")
				return domain.Task{}, nil
			},
		}

		handler := NewTasksHTTPHandler(service)

		handler.CreateTask(res, req)

		if res.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", res.Code, http.StatusBadRequest)
		}
	})

}
