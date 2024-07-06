package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time-tracker/internal/models"
	"time-tracker/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const PassportSerieLength = 4
const PassportNumberLength = 6

var UserParamsArr = [5]string{"passport_number", "name", "surname", "patronymic", "address"}

func (h *Handler) CreateUser(c *gin.Context) {
	var payload models.CreateUserPayload
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if err := validateCreateUserPayload(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.TODO()

	user, err := h.service.IUserService.CreateUser(ctx, &payload)
	if err != nil {
		if errors.Is(err, service.ErrDuplicateEntry) {
			c.JSON(http.StatusConflict, gin.H{"error": "User with this passport series and number already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) GetUsers(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset"})
		return
	}

	filters := make(map[string]string, len(UserParamsArr))
	for _, field := range UserParamsArr {
		if value := c.Query(field); value != "" {
			filters[field] = value
		}
	}

	ctx := context.TODO()

	users, err := h.service.IUserService.GetUsers(ctx, limit, offset, filters)
	if err != nil {
		if errors.Is(err, service.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "No users found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *Handler) GetUserByPassportNumber(c *gin.Context) {
	passportSerieParam := c.Query("passportSerie")
	passportNumberParam := c.Query("passportNumber")

	if passportSerieParam == "" || passportNumberParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passport series and number are required"})
		return
	}

	passportNumber := fmt.Sprintf("%s %s", passportSerieParam, passportNumberParam)

	if err := validatePassportNumber(passportNumber); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.TODO()

	user, err := h.service.IUserService.GetUserByPassportNumber(ctx, passportNumber)
	if err != nil {
		if errors.Is(err, service.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "No user found with this passport series and number"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) GetUser(c *gin.Context) {
	uuidString := c.Param("uuid")
	uuid, err := uuid.Parse(uuidString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	ctx := context.TODO()

	user, err := h.service.IUserService.GetUserByUUID(ctx, uuid)
	if err != nil {
		if errors.Is(err, service.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "No user found with this UUID"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	var payload models.UpdateUserPayload
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	uuidString := c.Param("uuid")
	uuid, err := uuid.Parse(uuidString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	if payload.PassportNumber != nil ||
		payload.Name != nil ||
		payload.Surname != nil ||
		payload.Patronymic != nil ||
		payload.Address != nil {

		ctx := context.TODO()

		user, err := h.service.IUserService.UpdateUserByUUID(ctx, uuid, &payload)
		if err != nil {
			if errors.Is(err, service.ErrDuplicateEntry) {
				c.JSON(http.StatusConflict, gin.H{"error": "user with this passport series and number already exists"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
		c.JSON(http.StatusOK, user)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "No fields to update"})
}

func (h *Handler) DeleteUser(c *gin.Context) {
	uuidString := c.Param("uuid")
	uuid, err := uuid.Parse(uuidString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	ctx := context.TODO()

	err = h.service.IUserService.DeleteUserByUUID(ctx, uuid)
	if err != nil {
		if errors.Is(err, service.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "No user found with this passport series and number"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func validateCreateUserPayload(payload *models.CreateUserPayload) error {
	if payload.PassportNumber == "" {
		return fmt.Errorf("passportNumber is required")
	}
	if payload.Name == "" {
		return fmt.Errorf("name is required")
	}
	if payload.Surname == "" {
		return fmt.Errorf("surname is required")
	}
	if payload.Address == "" {
		return fmt.Errorf("address is required")
	}
	return validatePassportNumber(payload.PassportNumber)
}

func validatePassportNumber(passportNumber string) error {
	matched, err := regexp.MatchString(`^\d{4} \d{6}$`, passportNumber)
	if err != nil || !matched {
		return fmt.Errorf("invalid passport number format")
	}
	return nil
}
