package tasks_service

import (
	"context"
	"errors"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/raizonw/todo-app/internal/core/domain"
	core_errors "github.com/raizonw/todo-app/internal/core/errors"
)

func ptr[T any](v T) *T {
	return &v
}

type fakeRepository struct {
	createTaskFunc func(ctx context.Context, task domain.Task) (domain.Task, error)
	getTasksFunc   func(ctx context.Context, userID *int, limit *int, offset *int) ([]domain.Task, error)
	getTaskFunc    func(ctx context.Context, id int) (domain.Task, error)
	deleteTaskFunc func(ctx context.Context, id int) error
	patchTaskFunc  func(ctx context.Context, id int, task domain.Task) (domain.Task, error)
	getTaskCalls   int
	patchTaskCalls int
}

func (f *fakeRepository) CreateTask(ctx context.Context, task domain.Task) (domain.Task, error) {
	return f.createTaskFunc(ctx, task)
}

func (f *fakeRepository) GetTasks(ctx context.Context, user_id *int, limit *int, offset *int) ([]domain.Task, error) {
	return f.getTasksFunc(ctx, user_id, limit, offset)
}

func (f *fakeRepository) GetTask(ctx context.Context, id int) (domain.Task, error) {
	f.getTaskCalls++
	return f.getTaskFunc(ctx, id)
}

func (f *fakeRepository) DeleteTask(ctx context.Context, id int) error {
	return f.deleteTaskFunc(ctx, id)
}

func (f *fakeRepository) PatchTask(ctx context.Context, id int, task domain.Task) (domain.Task, error) {
	f.patchTaskCalls++
	return f.patchTaskFunc(ctx, id, task)
}

func newTasksService(repo *fakeRepository) *TasksService {
	return NewTasksService(repo)
}

func TestCreateTask(t *testing.T) {
	createdAt := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	task := domain.NewTask(1, 1, "title", nil, false, createdAt, nil, 10)
	invalidTask := domain.NewTask(1, 1, strings.Repeat("a", 101), nil, false, createdAt, nil, 10)

	t.Run("validTask", func(t *testing.T) {
		called := false
		repo := fakeRepository{
			createTaskFunc: func(ctx context.Context, task domain.Task) (domain.Task, error) {
				called = true
				return task, nil
			},
		}
		service := newTasksService(&repo)

		result, err := service.CreateTask(t.Context(), task)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !called {
			t.Error("service doesn't call repository")
		}

		if result.ID != task.ID {
			t.Errorf("id got :%d, want: %d", result.ID, task.ID)
		}

		if result.Title != task.Title {
			t.Errorf("title got :%q, want: %q", result.Title, task.Title)
		}
	})
	t.Run("invalidTask", func(t *testing.T) {
		called := false
		repo := fakeRepository{
			createTaskFunc: func(ctx context.Context, task domain.Task) (domain.Task, error) {
				called = true
				return task, nil
			},
		}
		service := newTasksService(&repo)

		_, err := service.CreateTask(t.Context(), invalidTask)

		if !errors.Is(err, core_errors.ErrInvalidArgument) {
			t.Fatalf("expected ErrInvalidArgument, got: %v", err)
		}

		if called {
			t.Fatalf("repository must not be called for invalid task")
		}
	})

	t.Run("repository error", func(t *testing.T) {
		repositoryErr := errors.New("database unavailable")
		repo := fakeRepository{
			createTaskFunc: func(ctx context.Context, task domain.Task) (domain.Task, error) { return domain.Task{}, repositoryErr },
		}

		service := newTasksService(&repo)

		_, err := service.CreateTask(t.Context(), task)

		if !errors.Is(err, repositoryErr) {
			t.Fatalf("expected repositoryErr, got: %v", err)
		}
	})

	t.Run("passes correct task to repository", func(t *testing.T) {
		var receivedTask domain.Task

		repo := fakeRepository{
			createTaskFunc: func(ctx context.Context, task domain.Task) (domain.Task, error) {
				receivedTask = task
				return task, nil
			},
		}

		service := newTasksService(&repo)
		_, err := service.CreateTask(t.Context(), task)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if receivedTask != task {
			t.Errorf("repository recived task: %v, want: %v", receivedTask, task)
		}
	})
}

func TestGetTasks(t *testing.T) {
	createdAt := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	completedAt := createdAt.Add(time.Hour)

	tasks := []domain.Task{
		{
			ID:           1,
			Version:      1,
			Title:        "task1",
			Description:  nil,
			Completed:    true,
			CreatedAt:    createdAt,
			CompletedAt:  &completedAt,
			AuthorUserID: 10,
		},
		{
			ID:           2,
			Version:      1,
			Title:        "task2",
			Description:  ptr("description"),
			Completed:    false,
			CreatedAt:    createdAt,
			CompletedAt:  nil,
			AuthorUserID: 10,
		},
	}
	t.Run("valid parameters", func(t *testing.T) {
		user_id, limit, offset := 10, 1, 1
		var receivedID, receivedLimit, receivedOffset *int

		repo := fakeRepository{getTasksFunc: func(ctx context.Context, userID, limit, offset *int) ([]domain.Task, error) {
			receivedID = userID
			receivedLimit = limit
			receivedOffset = offset
			return tasks, nil
		}}

		service := newTasksService(&repo)
		result, err := service.GetTasks(t.Context(), &user_id, &limit, &offset)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !slices.Equal(result, tasks) {
			t.Errorf("slice wanted: %v, slice got: %v", tasks, result)
		}
		if receivedID != &user_id || receivedLimit != &limit || receivedOffset != &offset {
			t.Error("service passed incorrect parameters to repository")
		}
	})

	t.Run("negative limit", func(t *testing.T) {
		user_id, limit, offset := 10, -1, 1
		var called = false
		repo := fakeRepository{getTasksFunc: func(ctx context.Context, userID, limit, offset *int) ([]domain.Task, error) {
			called = true
			return tasks, nil
		}}

		service := newTasksService(&repo)
		_, err := service.GetTasks(t.Context(), &user_id, &limit, &offset)
		if !errors.Is(err, core_errors.ErrInvalidArgument) {
			t.Fatalf("expected ErrInvalidArgument, got: %v", err)
		}
		if called {
			t.Errorf("repository must not be called for invalid limit/offset")
		}
	})

	t.Run("repository error", func(t *testing.T) {
		repositoryErr := errors.New("repository error")
		user_id, limit, offset := 10, 1, 1
		repo := fakeRepository{getTasksFunc: func(ctx context.Context, userID, limit, offset *int) ([]domain.Task, error) {
			return []domain.Task{}, repositoryErr
		}}

		service := newTasksService(&repo)
		_, err := service.GetTasks(t.Context(), &user_id, &limit, &offset)

		if !errors.Is(err, repositoryErr) {
			t.Fatalf("expected repositoryErr, got: %v", err)
		}
	})
}

func TestPatchTask(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		var receivedGetID int
		var receivedPatchID int
		var receivedTask domain.Task
		createdAt := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
		task := domain.NewTask(1, 1, "Title", nil, false, createdAt, nil, 10)
		patch := domain.NewTaskPatch(domain.Nullable[string]{Value: ptr("ChangedTitle"), Set: true}, domain.Nullable[string]{Set: false}, domain.Nullable[bool]{Set: false})
		patchedTask := domain.NewTask(1, 2, "ChangedTitle", nil, false, createdAt, nil, 10)

		repo := fakeRepository{
			getTaskFunc: func(ctx context.Context, id int) (domain.Task, error) {
				receivedGetID = id
				return task, nil
			},
			patchTaskFunc: func(ctx context.Context, id int, task domain.Task) (domain.Task, error) {
				receivedPatchID = id
				receivedTask = task
				return patchedTask, nil
			},
		}

		service := newTasksService(&repo)
		result, err := service.PatchTask(t.Context(), task.ID, patch)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result != patchedTask {
			t.Errorf("result = %+v, want %+v", result, patchedTask)
		}

		if receivedGetID != task.ID {
			t.Errorf("GetTask id = %d, want %d", receivedGetID, task.ID)
		}

		if receivedPatchID != task.ID {
			t.Errorf("PatchTask id = %d, want %d", receivedPatchID, task.ID)
		}

		expectedTaskForRepo := domain.NewTask(
			1, 1, "ChangedTitle", nil, false, createdAt, nil, 10,
		)

		if receivedTask != expectedTaskForRepo {
			t.Errorf("repository received task = %v, want %v", receivedTask, expectedTaskForRepo)
		}

		if repo.getTaskCalls != 1 {
			t.Errorf("GetTask calls = %d, want 1", repo.getTaskCalls)
		}

		if repo.patchTaskCalls != 1 {
			t.Errorf("PatchTask calls = %d, want 1", repo.patchTaskCalls)
		}
	})

	t.Run("GetTask error", func(t *testing.T) {
		getTaskErr := errors.New("get task error")
		repo := fakeRepository{
			getTaskFunc: func(ctx context.Context, id int) (domain.Task, error) {
				return domain.Task{}, getTaskErr
			},
			patchTaskFunc: func(ctx context.Context, id int, task domain.Task) (domain.Task, error) {
				return domain.Task{}, nil
			},
		}

		service := newTasksService(&repo)
		_, err := service.PatchTask(t.Context(), 1, domain.TaskPatch{})

		if !errors.Is(err, getTaskErr) {
			t.Fatalf("expected getTaskErr, got %v", err)
		}
		if repo.getTaskCalls != 1 {
			t.Errorf("GetTask calls = %d, want 1", repo.getTaskCalls)
		}
		if repo.patchTaskCalls != 0 {
			t.Errorf("PatchTask calls = %d, want 0", repo.patchTaskCalls)
		}
	})

	t.Run("invalid Patch", func(t *testing.T) {
		task := domain.NewTask(1, 1, "title", nil, false, time.Now(), nil, 10)
		patch := domain.NewTaskPatch(domain.Nullable[string]{Value: ptr(strings.Repeat("a", 101)), Set: true}, domain.Nullable[string]{}, domain.Nullable[bool]{})
		repo := fakeRepository{
			getTaskFunc: func(ctx context.Context, id int) (domain.Task, error) {
				return task, nil
			},
			patchTaskFunc: func(ctx context.Context, id int, task domain.Task) (domain.Task, error) {
				return domain.Task{}, nil
			},
		}

		service := newTasksService(&repo)
		_, err := service.PatchTask(t.Context(), task.ID, patch)

		if !errors.Is(err, core_errors.ErrInvalidArgument) {
			t.Fatalf("expected ErrInvalidArgument, got %v", err)
		}
		if repo.getTaskCalls != 1 {
			t.Errorf("GetTask calls = %d, want 1", repo.getTaskCalls)
		}
		if repo.patchTaskCalls != 0 {
			t.Errorf("PatchTask calls = %d, want 0", repo.patchTaskCalls)
		}
	})

	t.Run("PatchTask repository error", func(t *testing.T) {
		repositoryErr := errors.New("repository error")
		task := domain.NewTask(1, 1, "title", nil, false, time.Now(), nil, 10)
		patch := domain.NewTaskPatch(domain.Nullable[string]{Value: ptr("ChangedTitle"), Set: true}, domain.Nullable[string]{Set: false}, domain.Nullable[bool]{Set: false})

		repo := fakeRepository{
			getTaskFunc: func(ctx context.Context, id int) (domain.Task, error) {
				return task, nil
			},
			patchTaskFunc: func(ctx context.Context, id int, task domain.Task) (domain.Task, error) {
				return domain.Task{}, repositoryErr
			},
		}

		service := newTasksService(&repo)
		_, err := service.PatchTask(t.Context(), task.ID, patch)

		if !errors.Is(err, repositoryErr) {
			t.Fatalf("expected repository error, got : %v", err)
		}

		if repo.getTaskCalls != 1 {
			t.Errorf("GetTask calls = %d, want 1", repo.getTaskCalls)
		}

		if repo.patchTaskCalls != 1 {
			t.Errorf("PatchTask calls = %d, want 1", repo.patchTaskCalls)
		}
	})

	t.Run("Context check", func(t *testing.T) {
		type contextKey string

		const key contextKey = "test-key"
		ctx := context.WithValue(t.Context(), key, "test-value")
		var getTaskContext any
		var patchTaskContext any
		createdAt := time.Now()
		task := domain.NewTask(1, 1, "Title", nil, false, createdAt, nil, 10)
		patch := domain.NewTaskPatch(domain.Nullable[string]{Value: ptr("ChangedTitle"), Set: true}, domain.Nullable[string]{Set: false}, domain.Nullable[bool]{Set: false})
		patchedTask := domain.NewTask(1, 2, "ChangedTitle", nil, false, createdAt, nil, 10)

		repo := fakeRepository{
			getTaskFunc: func(ctx context.Context, id int) (domain.Task, error) {
				getTaskContext = ctx.Value(key)
				return task, nil
			},
			patchTaskFunc: func(ctx context.Context, id int, task domain.Task) (domain.Task, error) {
				patchTaskContext = ctx.Value(key)
				return patchedTask, nil
			},
		}

		service := newTasksService(&repo)

		_, err := service.PatchTask(ctx, task.ID, patch)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if getTaskContext != ctx.Value(key) {
			t.Errorf("Get task context expected: %v, got: %v", getTaskContext, ctx.Value(key))
		}

		if patchTaskContext != ctx.Value(key) {
			t.Errorf("Get task context expected: %v, got: %v", patchTaskContext, ctx.Value(key))
		}
	})
}
