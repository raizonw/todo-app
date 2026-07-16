package domain

import (
	"errors"
	"strings"
	"testing"
	"time"

	core_errors "github.com/raizonw/todo-app/internal/core/errors"
)

func ptr[T any](v T) *T {
	return &v
}

func validTask(t *testing.T) Task {
	t.Helper()

	created_at := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	return NewTask(1, 1, "learn homework", ptr("description"), false, created_at, nil, 10)
}

func TestTask(t *testing.T) {
	createdAt := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	completedAt := createdAt.Add(time.Hour)
	invalidCompletedAt := createdAt.Add(-time.Hour)

	tests := []struct {
		name    string
		update  func(task *Task)
		wantErr bool
	}{
		{
			name:   "valid not completed task",
			update: func(task *Task) {},
		},
		{
			name: "empty title",
			update: func(task *Task) {
				task.Title = ""
			},
			wantErr: true,
		},
		{
			name: "too long title",
			update: func(task *Task) {
				task.Title = strings.Repeat("a", 101)
			},
			wantErr: true,
		},
		{
			name: "completed without completed_at",
			update: func(task *Task) {
				task.Completed = true
			},
			wantErr: true,
		},
		{
			name: "completed before created_at",
			update: func(task *Task) {
				task.CompletedAt = &invalidCompletedAt
			},
			wantErr: true,
		},
		{
			name: "not completed with completed_at",
			update: func(task *Task) {
				task.CompletedAt = &completedAt
			},
			wantErr: true,
		},
		{
			name: "completed with completed_at",
			update: func(task *Task) {
				task.Completed = true
				task.CompletedAt = &completedAt
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := validTask(t)
			tt.update(&task)
			err := task.Validate()

			if tt.wantErr {
				if !errors.Is(err, core_errors.ErrInvalidArgument) {
					t.Fatalf("expected ErrInvalidArgument, got: %v", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestTaskApplyPatch(t *testing.T) {
	task := validTask(t)

	newTitle := "New title"
	patch := NewTaskPatch(
		Nullable[string]{Value: &newTitle, Set: true},
		Nullable[string]{},
		Nullable[bool]{},
	)

	err := task.ApplyPatch(patch)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if task.Title != newTitle {
		t.Errorf("Title = %q, want %q", task.Title, newTitle)
	}
}
