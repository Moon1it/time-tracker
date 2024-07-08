package models

import (
	"time"

	"github.com/google/uuid"
)

type CreateUserPayload struct {
	PassportNumber string  `json:"passportNumber"`
	Surname        string  `json:"surname"`
	Name           string  `json:"name"`
	Patronymic     *string `json:"patronymic"`
	Address        string  `json:"address"`
}

type UpdateUserPayload struct {
	PassportNumber *string `json:"passportNumber"`
	Surname        *string `json:"surname"`
	Name           *string `json:"name"`
	Patronymic     *string `json:"patronymic"`
	Address        *string `json:"address"`
}

type User struct {
	UUID           uuid.UUID `json:"uuid"`
	PassportNumber string    `json:"passportNumber"`
	Surname        string    `json:"surname"`
	Name           string    `json:"name"`
	Patronymic     *string   `json:"patronymic,omitempty"`
	Address        string    `json:"address"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
