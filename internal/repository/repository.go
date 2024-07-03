package repository

import (
	"time-tracker/internal/models"

	"github.com/jmoiron/sqlx"
)

type People interface {
	CreatePeople(passportSerie int, passportNumber int) (*models.People, error)
	GetPeople(passportSerie int, passportNumber int) (*models.People, error)
}

type Repository struct {
	People
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		People: NewPeopleRepository(db),
	}
}
