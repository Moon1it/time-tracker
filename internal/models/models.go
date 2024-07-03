package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreatePeoplePayload struct {
	PassportNumber string `json:"passportNumber"`
}

type People struct {
	ID             uuid.UUID   `db:"id"`
	PassportSerie  int         `db:"passport_serie"`
	PassportNumber int         `db:"passport_number"`
	Surname        pgtype.Text `db:"surname"`
	Name           pgtype.Text `db:"name"`
	Patronymic     pgtype.Text `db:"patronymic"`
	Address        pgtype.Text `db:"address"`
	CreatedAt      time.Time   `db:"created_at"`
	UpdatedAt      time.Time   `db:"updated_at"`
}
