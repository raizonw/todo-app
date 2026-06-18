package tasks_service

import (
	"context"

	"github.com/raizonw/todo-app/internal/core/domain"
)

type TasksService struct {
	tasksRepository TasksRepository
}

type TasksRepository interface {
	CreateTask(
		ctx context.Context,
		task domain.Task,
	) (domain.Task, error)

	GetTasks(
		ctx context.Context,
		user_id *int,
		limit *int,
		offset *int,
	) ([]domain.Task, error)

	GetTask(
		ctx context.Context,
		id int,
	) (domain.Task, error)

	DeleteTask(
		ctx context.Context,
		id int,
	) error
}

func NewTasksService(
	taskRepository TasksRepository,
) *TasksService {
	return &TasksService{
		tasksRepository: taskRepository,
	}
}
