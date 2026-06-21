package statistics_service

import (
	"context"
	"time"

	"github.com/raizonw/todo-app/internal/core/domain"
)

type StatisticsService struct {
	statisticsRepository StatisticsRepository
}

type StatisticsRepository interface {
	GetTasks(
		ctx context.Context,
		userID *int,
		from *time.Time,
		to *time.Time,
	) ([]domain.Task, error)
}

func NewStatisticService(
	statisticRepository StatisticsRepository,
) *StatisticsService {
	return &StatisticsService{
		statisticsRepository: statisticRepository,
	}
}
