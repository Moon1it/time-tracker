package service

import (
	"context"
	sqlc "time-tracker/internal/db/sqlc"
	"time-tracker/internal/models"

	"github.com/google/uuid"
)

//go:generate mockery --name IUserService
type IUserService interface {
	CreateUser(ctx context.Context, payload *models.CreateUserPayload) (*models.User, error)
	GetUsers(ctx context.Context, limit, offset int, filters map[string]string) ([]models.User, error)
	GetUserByUUID(ctx context.Context, UUID uuid.UUID) (*models.User, error)
	GetUserByPassportNumber(ctx context.Context, passportNumber string) (*models.User, error)
	UpdateUserByUUID(ctx context.Context, UUID uuid.UUID, payload *models.UpdateUserPayload) (*models.User, error)
	DeleteUserByUUID(ctx context.Context, UUID uuid.UUID) error
}

//go:generate mockery --name ITaskService
type ITaskService interface {
	CreateTask(ctx context.Context, userUUID uuid.UUID, payload *models.CreateTaskPayload) (*models.Task, error)
	FinishTask(ctx context.Context, userUUID uuid.UUID) (*models.CompletedTask, error)
	GetTasksResult(ctx context.Context, userUUID uuid.UUID, days int) (*models.TasksResult, error)
}

type Service struct {
	IUserService
	ITaskService
}

func NewService(repository sqlc.Querier) *Service {
	return &Service{
		IUserService: NewUserService(repository),
		ITaskService: NewTaskService(repository),
	}
}
