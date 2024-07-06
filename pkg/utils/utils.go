package utils

import (
	db "time-tracker/internal/db/sqlc"
	"time-tracker/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func ConvertDBUserToModelsUser(user db.User) (*models.User, error) {
	var id uuid.UUID
	err := id.UnmarshalBinary(user.Uuid.Bytes[:])
	if err != nil {
		return nil, err
	}

	var patronymic *string
	if user.Patronymic.Valid {
		patronymic = &user.Patronymic.String
	}

	return &models.User{
		UUID:           id,
		PassportNumber: user.PassportNumber,
		Surname:        user.Surname,
		Name:           user.Name,
		Patronymic:     patronymic,
		Address:        user.Address,
		CreatedAt:      user.CreatedAt.Time,
		UpdatedAt:      user.UpdatedAt.Time,
	}, nil
}

func ToPgText(s *string) pgtype.Text {
	if s != nil {
		return pgtype.Text{String: *s, Valid: true}
	}
	return pgtype.Text{Valid: false}
}
