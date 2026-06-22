package statistics_service_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/raizonw/todo-app/internal/core/domain"
	core_errors "github.com/raizonw/todo-app/internal/core/errors"
	statistics_service "github.com/raizonw/todo-app/internal/features/statistics/service"
)

type mockStatisticsRepository struct {
	getTasksFunc func(ctx context.Context, userID *int, from *time.Time, to *time.Time) ([]domain.Task, error)
}

func (m *mockStatisticsRepository) GetTasks(ctx context.Context, userID *int, from *time.Time, to *time.Time) ([]domain.Task, error) {
	if m.getTasksFunc != nil {
		return m.getTasksFunc(ctx, userID, from, to)
	}
	panic("mock method GetTasks not configured")
}

func TestGetStatistics(t *testing.T) {
	now := time.Now()
	oneHourAgo := now.Add(-1 * time.Hour)
	twoHoursAgo := now.Add(-2 * time.Hour)

	type testCase struct {
		name          string
		userID        *int
		from          *time.Time
		to            *time.Time
		setupMock     func(m *mockStatisticsRepository)
		expectedStats domain.Statistics
		expectedErr   error
	}

	tests := []testCase{
		{
			name:   "Success - empty tasks list",
			userID: ptr(1),
			from:   &twoHoursAgo,
			to:     &now,
			setupMock: func(m *mockStatisticsRepository) {
				m.getTasksFunc = func(ctx context.Context, userID *int, from *time.Time, to *time.Time) ([]domain.Task, error) {
					return []domain.Task{}, nil
				}
			},
			expectedStats: domain.Statistics{
				TaskCreated:               0,
				TaskCompleted:             0,
				TaskCompletedRate:         nil,
				TaskAverageCompletionTime: nil,
			},
			expectedErr: nil,
		},
		{
			name:   "Success - calculated statistics correctly",
			userID: ptr(1),
			from:   &twoHoursAgo,
			to:     &now,
			setupMock: func(m *mockStatisticsRepository) {
				m.getTasksFunc = func(ctx context.Context, userID *int, from *time.Time, to *time.Time) ([]domain.Task, error) {
					return []domain.Task{
						// Task 1: Completed, completion time = 1 hour
						{
							ID:           1,
							Completed:    true,
							CreatedAt:    twoHoursAgo,
							CompletedAt:  &oneHourAgo,
							AuthorUserID: 1,
						},
						// Task 2: Completed, completion time = 2 hours
						{
							ID:           2,
							Completed:    true,
							CreatedAt:    twoHoursAgo,
							CompletedAt:  &now,
							AuthorUserID: 1,
						},
						// Task 3: Not completed
						{
							ID:           3,
							Completed:    false,
							CreatedAt:    twoHoursAgo,
							CompletedAt:  nil,
							AuthorUserID: 1,
						},
					}, nil
				}
			},
			expectedStats: domain.Statistics{
				TaskCreated:               3,
				TaskCompleted:             2,
				TaskCompletedRate:         ptr(66.66666666666666),
				TaskAverageCompletionTime: ptr(90 * time.Minute), // (1h + 2h) / 2 = 1.5h
			},
			expectedErr: nil,
		},
		{
			name:        "Error - to is before from",
			userID:      ptr(1),
			from:        &now,
			to:          &twoHoursAgo,
			setupMock:   func(m *mockStatisticsRepository) {},
			expectedErr: core_errors.ErrInvalidArgument,
		},
		{
			name:        "Error - to is equal to from",
			userID:      ptr(1),
			from:        &now,
			to:          &now,
			setupMock:   func(m *mockStatisticsRepository) {},
			expectedErr: core_errors.ErrInvalidArgument,
		},
		{
			name:   "Error - repository failure",
			userID: ptr(1),
			from:   &twoHoursAgo,
			to:     &now,
			setupMock: func(m *mockStatisticsRepository) {
				m.getTasksFunc = func(ctx context.Context, userID *int, from *time.Time, to *time.Time) ([]domain.Task, error) {
					return nil, errors.New("db error")
				}
			},
			expectedErr: errors.New("get tasks from repository: db error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockRepo := &mockStatisticsRepository{}
			tc.setupMock(mockRepo)

			service := statistics_service.NewStatisticService(mockRepo)
			res, err := service.GetStatistics(context.Background(), tc.userID, tc.from, tc.to)

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

			if res.TaskCreated != tc.expectedStats.TaskCreated {
				t.Errorf("expected TaskCreated %d, got %d", tc.expectedStats.TaskCreated, res.TaskCreated)
			}
			if res.TaskCompleted != tc.expectedStats.TaskCompleted {
				t.Errorf("expected TaskCompleted %d, got %d", tc.expectedStats.TaskCompleted, res.TaskCompleted)
			}

			if tc.expectedStats.TaskCompletedRate == nil {
				if res.TaskCompletedRate != nil {
					t.Errorf("expected TaskCompletedRate to be nil, got %v", *res.TaskCompletedRate)
				}
			} else {
				if res.TaskCompletedRate == nil {
					t.Errorf("expected TaskCompletedRate to be %f, got nil", *tc.expectedStats.TaskCompletedRate)
				} else if *res.TaskCompletedRate != *tc.expectedStats.TaskCompletedRate {
					t.Errorf("expected TaskCompletedRate %f, got %f", *tc.expectedStats.TaskCompletedRate, *res.TaskCompletedRate)
				}
			}

			if tc.expectedStats.TaskAverageCompletionTime == nil {
				if res.TaskAverageCompletionTime != nil {
					t.Errorf("expected TaskAverageCompletionTime to be nil, got %v", *res.TaskAverageCompletionTime)
				}
			} else {
				if res.TaskAverageCompletionTime == nil {
					t.Errorf("expected TaskAverageCompletionTime to be %v, got nil", *tc.expectedStats.TaskAverageCompletionTime)
				} else if *res.TaskAverageCompletionTime != *tc.expectedStats.TaskAverageCompletionTime {
					t.Errorf("expected TaskAverageCompletionTime %v, got %v", *tc.expectedStats.TaskAverageCompletionTime, *res.TaskAverageCompletionTime)
				}
			}
		})
	}
}

func ptr[T any](v T) *T {
	return &v
}
