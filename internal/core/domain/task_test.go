package domain

import (
	"errors"
	"strings"
	"testing"
	"time"

	core_errors "github.com/raizonw/todo-app/internal/core/errors"
)

func TestTask(t *testing.T) {
	createdAt := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	completedAt := createdAt.Add(time.Hour)
	invalidCompletedAt := createdAt.Add(-time.Hour)

	tests := []struct {
		name    string
		task    Task
		wantErr bool
	}{
		{
			name: "valid not completed task",
			task: NewTask(1, 1, "learn homework", nil, false, createdAt, nil, 10),
		},
		{
			name:    "empty title",
			task:    NewTask(1, 1, "", nil, false, createdAt, &completedAt, 10),
			wantErr: true,
		},
		{
			name:    "too long title",
			task:    NewTask(1, 1, strings.Repeat("a", 101), nil, false, createdAt, &completedAt, 10),
			wantErr: true,
		},
		{
			name:    "completed without completed_at",
			task:    NewTask(1, 1, "done", nil, true, createdAt, nil, 10),
			wantErr: true,
		},
		{
			name:    "completed before created_at",
			task:    NewTask(1, 1, "completed before created", nil, true, createdAt, &invalidCompletedAt, 10),
			wantErr: true,
		},
		{
			name: "completed with completed_at",
			task: NewTask(1, 1, "done", nil, true, createdAt, &completedAt, 10),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.task.Validate()

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
