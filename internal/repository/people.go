package repository

import (
	"errors"
	"fmt"
	"time"
	"time-tracker/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

var ErrDuplicateEntry = errors.New("people with this passport series and number already exists")
var ErrNoRows = errors.New("no user found with this passport series and number")

type PeopleRepository struct {
	db *sqlx.DB
}

func NewPeopleRepository(db *sqlx.DB) *PeopleRepository {
	return &PeopleRepository{db}
}

type CreatePeopleParams struct {
	PassportSerie  int       `db:"passport_serie" json:"-"`
	PassportNumber int       `db:"passport_number" json:"-"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

func (pr *PeopleRepository) CreatePeople(passportSerie int, passportNumber int) (*models.People, error) {
	fmt.Println("repository.CreatePeople", passportSerie, passportNumber)

	newPeople := CreatePeopleParams{
		PassportSerie:  passportSerie,
		PassportNumber: passportNumber,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	query := `
		INSERT INTO people (passport_serie, passport_number, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, passport_serie, passport_number, created_at, updated_at
	`

	var people models.People
	err := pr.db.QueryRow(query,
		newPeople.PassportSerie,
		newPeople.PassportNumber,
		newPeople.CreatedAt,
		newPeople.UpdatedAt,
	).Scan(&people.ID, &people.PassportSerie, &people.PassportNumber, &people.CreatedAt, &people.UpdatedAt)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return nil, ErrDuplicateEntry
		}
		return nil, err
	}

	return &people, nil
}

func (pr *PeopleRepository) GetPeople(passportSerie int, passportNumber int) (*models.People, error) {
	var people models.People
	query := "SELECT id, passport_serie, passport_number, surname, name, patronymic, address, created_at, updated_at FROM people WHERE passport_serie = $1 AND passport_number = $2"

	err := pr.db.QueryRow(query,
		passportSerie,
		passportNumber,
	).Scan(&people.ID, &people.PassportSerie, &people.PassportNumber, &people.Surname, &people.Name, &people.Patronymic, &people.Address, &people.CreatedAt, &people.UpdatedAt)
	if err != nil {
		fmt.Println("err: ", err)
		fmt.Println("pgx err: ", pgx.ErrNoRows)
		if err == pgx.ErrNoRows {
			return nil, ErrNoRows
		}
		return nil, err
	}

	return &people, nil
}
