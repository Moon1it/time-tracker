package service

import (
	"time-tracker/internal/models"
	"time-tracker/internal/repository"
)

type PeopleService struct {
	peopleRepo repository.People
}

func NewPeopleService(repo repository.People) *PeopleService {
	return &PeopleService{
		peopleRepo: repo,
	}
}

func (ps *PeopleService) CreatePeople(passportSerie int, passportNumber int) (*models.People, error) {
	return ps.peopleRepo.CreatePeople(passportSerie, passportNumber)
}

func (ps *PeopleService) GetPeople(passportSerie int, passportNumber int) (*models.People, error) {
	return ps.peopleRepo.GetPeople(passportSerie, passportNumber)
}
