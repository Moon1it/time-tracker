package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time-tracker/internal/models"
	"time-tracker/internal/repository"

	"github.com/gin-gonic/gin"
)

const PassportSerieLength = 4
const PassportNumberLength = 6

func (h *Handler) CreatePeople(c *gin.Context) {
	var newPeople models.CreatePeoplePayload

	if err := c.BindJSON(&newPeople); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if newPeople.PassportNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passport number is required"})
		return
	}

	passportNumberParts := strings.Split(newPeople.PassportNumber, " ")
	if len(passportNumberParts) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid passport number format"})
		return
	}

	passportSerieString := passportNumberParts[0]
	passportNumberString := passportNumberParts[1]

	if len(passportSerieString) != PassportSerieLength || len(passportNumberString) != PassportNumberLength {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid passport series or number length"})
		return
	}

	passportSerie, err := strconv.Atoi(passportSerieString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid passport series format"})
		return
	}

	passportNumber, err := strconv.Atoi(passportNumberString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid passport number format"})
		return
	}

	people, err := h.service.CreatePeople(passportSerie, passportNumber)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateEntry) {
			c.JSON(http.StatusConflict, gin.H{"error": "People with this passport series and number already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, people)
}

func (h *Handler) GetPeople(c *gin.Context) {
	passportSerie := c.Query("passportSerie")
	passportNumber := c.Query("passportNumber")

	if passportSerie == "" || passportNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passport series and number are required"})
		return
	}

	passportSerieInt, err := strconv.Atoi(passportSerie)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid passport series"})
		return
	}

	passportNumberInt, err := strconv.Atoi(passportNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid passport number"})
		return
	}

	people, err := h.service.GetPeople(passportSerieInt, passportNumberInt)
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "No user found with this passport series and number"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, people)
}

// func (h *Handler) GetPeoples(c *gin.Context) {
// 	c.JSON(http.StatusOK, &People{
// 		Surname:    "Иванов",
// 		Name:       "Иван",
// 		Patronymic: "Иванович",
// 		Address:    "г. Москва, ул. Ленина, д. 5, кв. 1",
// 	})
// }

// func (h *Handler) UpdatePeople(c *gin.Context) {
// 	c.JSON(http.StatusOK, &People{
// 		Name:    "Kirill",
// 		Surname: "Smirnov",
// 	})
// }

// func (h *Handler) DeletePeople(c *gin.Context) {
// 	c.JSON(http.StatusOK, &People{
// 		Surname:    "Иванов",
// 		Name:       "Иван",
// 		Patronymic: "Иванович",
// 		Address:    "г. Москва, ул. Ленина, д. 5, кв. 1",
// 	})
// }
