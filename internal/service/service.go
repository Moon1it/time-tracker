package service

import (
	"time-tracker/internal/models"
	"time-tracker/internal/repository"
)

type Service struct {
	People
}

type People interface {
	CreatePeople(passportSerie int, passportNumber int) (*models.People, error)
	GetPeople(passportSerie int, passportNumber int) (*models.People, error)
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		People: NewPeopleService(repos.People),
	}
}
