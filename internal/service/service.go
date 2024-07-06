package service

import (
	"context"
	sqlc "time-tracker/internal/db/sqlc"
	"time-tracker/internal/models"

	"github.com/google/uuid"
)

type IUserService interface {
	CreateUser(ctx context.Context, payload *models.CreateUserPayload) (*models.User, error)
	GetUsers(ctx context.Context, limit, offset int, filters map[string]string) ([]models.User, error)
	GetUserByUUID(ctx context.Context, UUID uuid.UUID) (*models.User, error)
	GetUserByPassportNumber(ctx context.Context, passportNumber string) (*models.User, error)
	UpdateUserByUUID(ctx context.Context, UUID uuid.UUID, upayload *models.UpdateUserPayload) (*models.User, error)
	DeleteUserByUUID(ctx context.Context, UUID uuid.UUID) error
}

type Service struct {
	IUserService
}

func NewService(repository sqlc.Querier) *Service {
	return &Service{
		IUserService: NewUserService(repository),
	}
}
