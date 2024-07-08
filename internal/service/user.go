package service

import (
	"context"
	"errors"
	"fmt"
	db "time-tracker/internal/db/sqlc"
	"time-tracker/internal/models"
	"time-tracker/pkg/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrUsersNotFound       = errors.New("no users found")
	ErrForeignKeyViolation = errors.New("foreign key violation")
)

type UserService struct {
	repository db.Querier
}

func NewUserService(repository db.Querier) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (ps *UserService) CreateUser(ctx context.Context, payload *models.CreateUserPayload) (*models.User, error) {
	var patronymic pgtype.Text
	if payload.Patronymic != nil {
		patronymic = pgtype.Text{String: *payload.Patronymic, Valid: true}
	}

	params := db.CreateUserParams{
		PassportNumber: payload.PassportNumber,
		Name:           payload.Name,
		Surname:        payload.Surname,
		Patronymic:     patronymic,
		Address:        payload.Address,
	}

	userRaw, err := ps.repository.CreateUser(ctx, params)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return nil, ErrUserAlreadyExists
		}
		return nil, err
	}

	user, err := utils.ConvertDBUserToModelsUser(userRaw)
	if err != nil {
		return nil, fmt.Errorf("error converting user: %v", err)
	}

	return user, nil
}

func (ps *UserService) GetUsers(ctx context.Context, limit, offset int, filters map[string]string) ([]models.User, error) {
	params := db.GetUsersParams{
		UserLimit:  int32(limit),
		UserOffset: int32(offset),
	}

	if passportNumber, ok := filters["passport_number"]; ok {
		params.PassportNumber = utils.ToPgText(&passportNumber)
	}
	if name, ok := filters["name"]; ok {
		params.Name = utils.ToPgText(&name)
	}
	if surname, ok := filters["surname"]; ok {
		params.Surname = utils.ToPgText(&surname)
	}
	if patronymic, ok := filters["patronymic"]; ok {
		params.Patronymic = utils.ToPgText(&patronymic)
	}
	if address, ok := filters["address"]; ok {
		params.Address = utils.ToPgText(&address)
	}

	usersRaw, err := ps.repository.GetUsers(ctx, params)
	if err != nil {
		return nil, err
	}

	if len(usersRaw) == 0 {
		return nil, ErrUsersNotFound
	}

	users := make([]models.User, len(usersRaw))
	for i, userRaw := range usersRaw {
		user, err := utils.ConvertDBUserToModelsUser(userRaw)
		if err != nil {
			return nil, fmt.Errorf("error converting user: %v", err)
		}
		users[i] = *user
	}
	return users, nil
}

func (ps *UserService) GetUserByUUID(ctx context.Context, UUID uuid.UUID) (*models.User, error) {
	pgUUID := pgtype.UUID{Bytes: UUID, Valid: true}

	userRaw, err := ps.repository.GetUserByUUID(ctx, pgUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	user, err := utils.ConvertDBUserToModelsUser(userRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to convert user: %w", err)
	}

	return user, nil
}

func (ps *UserService) GetUserByPassportNumber(ctx context.Context, passportNumber string) (*models.User, error) {
	userRaw, err := ps.repository.GetUserByPassportNumber(ctx, passportNumber)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	user, err := utils.ConvertDBUserToModelsUser(userRaw)
	if err != nil {
		return nil, fmt.Errorf("Error converting user: %v", err)
	}

	return user, nil
}

func (ps *UserService) UpdateUserByUUID(ctx context.Context, UUID uuid.UUID, payload *models.UpdateUserPayload) (*models.User, error) {
	params := db.UpdateUserByUUIDParams{
		UserUuid:       pgtype.UUID{Bytes: UUID, Valid: true},
		Name:           utils.ToPgText(payload.Name),
		Surname:        utils.ToPgText(payload.Surname),
		Patronymic:     utils.ToPgText(payload.Patronymic),
		Address:        utils.ToPgText(payload.Address),
		PassportNumber: utils.ToPgText(payload.PassportNumber),
	}

	userRaw, err := ps.repository.UpdateUserByUUID(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return nil, ErrUserAlreadyExists
		}
		return nil, err
	}

	user, err := utils.ConvertDBUserToModelsUser(userRaw)
	if err != nil {
		return nil, fmt.Errorf("Error converting user: %v", err)
	}

	return user, nil
}

func (ps *UserService) DeleteUserByUUID(ctx context.Context, UUID uuid.UUID) error {
	pgUUID := pgtype.UUID{Bytes: UUID, Valid: true}

	_, err := ps.repository.GetUserByUUID(ctx, pgUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}

	if err := ps.repository.DeleteUserByUUID(ctx, pgUUID); err != nil {
		return err
	}

	return nil
}
