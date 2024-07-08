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
	"github.com/sirupsen/logrus"
)

const PassportSerieLength = 4
const PassportNumberLength = 6

var UserParamsArr = [5]string{"passport_number", "name", "surname", "patronymic", "address"}

// @Summary Create a new user
// @Tags users
// @Description Create a new user with the given payload.
// @Accept  json
// @Produce  json
// @Param payload body models.CreateUserPayload true "User creation payload"
// @Success 201 {object} models.User "User created successfully"
// @Failure 400 {object} errorResponse "Bad request"
// @Failure 409 {object} errorResponse "User already exists"
// @Failure 500 {object} errorResponse "Internal server error"
// @Router /api/users [post]
func (h *Handler) CreateUser(c *gin.Context) {
	var payload models.CreateUserPayload
	if err := c.BindJSON(&payload); err != nil {
		logrus.Errorf("Error binding JSON: %v", err)
		newErrorResponse(c, http.StatusBadRequest, "Bad request")
		return
	}

	if err := validateCreateUserPayload(&payload); err != nil {
		logrus.Errorf("Validation error: %v", err)
		newErrorResponse(c, http.StatusBadRequest, "Bad request")
		return
	}

	ctx := c.Request.Context()
	user, err := h.service.IUserService.CreateUser(ctx, &payload)
	if err != nil {
		logrus.Errorf("Error creating user: %v", err)
		if errors.Is(err, service.ErrUserAlreadyExists) {
			newErrorResponse(c, http.StatusConflict, "User already exists")
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, "Internal server error")
		return
	}

	logrus.Infof("User created successfully: %v", user)
	c.JSON(http.StatusCreated, user)
}

// @Summary Get all users
// @Tags users
// @Description Retrieve a list of users with optional filters, limit, and offset.
// @Accept  json
// @Produce  json
// @Param limit query int false "Limit the number of users returned" default(10)
// @Param offset query int false "Offset the number of users returned" default(0)
// @Param filters query string false "Optional filters to apply on users"
// @Success 200 {array}  models.User "List of users"
// @Failure 400 {object} errorResponse "Bad request"
// @Failure 404 {object} errorResponse "No users found"
// @Failure 500 {object} errorResponse "Internal server error"
// @Router /api/users [get]
func (h *Handler) GetAllUsers(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		logrus.Errorf("Invalid limit parameter: %v", err)
		newErrorResponse(c, http.StatusBadRequest, "Bad request")
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		logrus.Errorf("Invalid offset parameter: %v", err)
		newErrorResponse(c, http.StatusBadRequest, "Bad request")
		return
	}

	filters := make(map[string]string, len(UserParamsArr))
	for _, field := range UserParamsArr {
		if value := c.Query(field); value != "" {
			filters[field] = value
		}
	}

	ctx := c.Request.Context()
	users, err := h.service.IUserService.GetUsers(ctx, limit, offset, filters)
	if err != nil {
		if errors.Is(err, service.ErrUsersNotFound) {
			logrus.Infof("No users found with filters: %v", filters)
			newErrorResponse(c, http.StatusNotFound, "No users found")
			return
		}
		logrus.Errorf("Error retrieving users: %v", err)
		newErrorResponse(c, http.StatusInternalServerError, "Internal server error")
		return
	}

	logrus.Infof("Retrieved %d users", len(users))
	c.JSON(http.StatusOK, users)
}

// @Summary People info
// @Tags users
// @Description Get user information by providing passport series and number.
// @Accept  json
// @Produce  json
// @Param passportSerie query string true "Passport Series"
// @Param passportNumber query string true "Passport Number"
// @Success 200 {object} models.User "User info"
// @Failure 400 {object} errorResponse "Bad request"
// @Failure 404 {object} errorResponse "User not found"
// @Failure 500 {object} errorResponse "Internal server error"
// @Router /api/users/info [get]
func (h *Handler) GetUserByPassportNumber(c *gin.Context) {
	passportSerieParam := c.Query("passportSerie")
	passportNumberParam := c.Query("passportNumber")

	if passportSerieParam == "" || passportNumberParam == "" {
		logrus.Warn("Missing passportSerie or passportNumber")
		newErrorResponse(c, http.StatusBadRequest, "Bad request")
		return
	}

	passportNumber := fmt.Sprintf("%s %s", passportSerieParam, passportNumberParam)
	if err := validatePassportNumber(passportNumber); err != nil {
		logrus.Errorf("Invalid passport number format: %v", err)
		newErrorResponse(c, http.StatusBadRequest, "Bad request")
		return
	}

	ctx := c.Request.Context()
	user, err := h.service.IUserService.GetUserByPassportNumber(ctx, passportNumber)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			logrus.Infof("User not found for passport number: %s", passportNumber)
			newErrorResponse(c, http.StatusNotFound, "User not found")
			return
		}
		logrus.Errorf("Error retrieving user by passport number: %v", err)
		newErrorResponse(c, http.StatusInternalServerError, "Internal server error")
		return
	}

	logrus.Infof("User retrieved successfully for passport number: %s", passportNumber)
	c.JSON(http.StatusOK, user)
}

// @Summary Get user by UUID
// @Tags users
// @Description Retrieve a user by their UUID.
// @Accept  json
// @Produce  json
// @Param uuid path string true "User UUID"
// @Success 200 {object} models.User "User retrieved successfully"
// @Failure 400 {object} errorResponse "Bad request"
// @Failure 404 {object} errorResponse "No users found"
// @Failure 500 {object} errorResponse "Internal server error"
// @Router /api/users/{uuid} [get]
func (h *Handler) GetUser(c *gin.Context) {
	uuidString := c.Param("uuid")
	uuid, err := uuid.Parse(uuidString)
	if err != nil {
		logrus.Errorf("Invalid UUID format: %v", err)
		newErrorResponse(c, http.StatusBadRequest, "Bad request")
		return
	}

	ctx := c.Request.Context()
	user, err := h.service.IUserService.GetUserByUUID(ctx, uuid)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			logrus.Infof("No user found for UUID: %s", uuidString)
			newErrorResponse(c, http.StatusNotFound, "No users found")
			return
		}
		logrus.Errorf("Error retrieving user by passport number: %v", err)
		newErrorResponse(c, http.StatusInternalServerError, "Internal server error")
		return
	}

	logrus.Infof("User retrieved successfully: %v", user)
	c.JSON(http.StatusOK, user)
}

// @Summary Update user by UUID
// @Tags users
// @Description Update a user's details by their UUID.
// @Accept  json
// @Produce  json
// @Param uuid path string true "User UUID"
// @Param payload body models.UpdateUserPayload true "User Update Payload"
// @Success 200 {object} models.User "User updated successfully"
// @Success 200 {object} statusResponse "No fields to update"
// @Failure 400 {object} errorResponse "Bad request"
// @Failure 404 {object} errorResponse "No users found"
// @Failure 409 {object} errorResponse "User with this passportNumber already exists"
// @Failure 500 {object} errorResponse "Internal server error"
// @Router /api/users/{uuid} [patch]
func (h *Handler) UpdateUser(c *gin.Context) {
	var payload models.UpdateUserPayload
	if err := c.BindJSON(&payload); err != nil {
		logrus.Errorf("Invalid JSON: %v", err)
		newErrorResponse(c, http.StatusBadRequest, "Bad request")
		return
	}

	uuidString := c.Param("uuid")
	uuid, err := uuid.Parse(uuidString)
	if err != nil {
		logrus.Errorf("Invalid UUID format: %v", err)
		newErrorResponse(c, http.StatusBadRequest, "Bad request")
		return
	}

	if payload.PassportNumber != nil ||
		payload.Name != nil ||
		payload.Surname != nil ||
		payload.Patronymic != nil ||
		payload.Address != nil {

		ctx := c.Request.Context()
		user, err := h.service.IUserService.UpdateUserByUUID(ctx, uuid, &payload)
		if err != nil {
			logrus.Errorf("Error updating user: %v", err)
			if errors.Is(err, service.ErrUserNotFound) {
				logrus.Infof("No user found for UUID: %s", uuidString)
				newErrorResponse(c, http.StatusNotFound, "No users found")
				return
			}
			if errors.Is(err, service.ErrUserAlreadyExists) {
				logrus.Warnf("User with passport number already exists: %v", err)
				newErrorResponse(c, http.StatusConflict, "user with this passportNumber already exists")
				return
			}
			newErrorResponse(c, http.StatusInternalServerError, "Internal server error")
			return
		}

		logrus.Infof("User updated successfully: %v", user)
		c.JSON(http.StatusOK, user)
		return
	}

	logrus.Infof("No fields to update for UUID: %s", uuidString)
	c.JSON(http.StatusOK, statusResponse{Description: "No fields to update"})
}

// @Summary Delete user by UUID
// @Tags users
// @Description Delete a user by their UUID.
// @Accept  json
// @Produce  json
// @Param uuid path string true "User UUID"
// @Success 200 {object} statusResponse "User deleted successfully"
// @Failure 400 {object} errorResponse "Bad request"
// @Failure 404 {object} errorResponse "No users found"
// @Failure 500 {object} errorResponse "Internal server error"
// @Router /api/users/{uuid} [delete]
func (h *Handler) DeleteUser(c *gin.Context) {
	uuidStr := c.Param("uuid")
	UUID, err := uuid.Parse(uuidStr)
	if err != nil {
		logrus.Errorf("Invalid UUID format: %v", err)
		newErrorResponse(c, http.StatusBadRequest, "Bad request")
		return
	}

	err = h.service.IUserService.DeleteUserByUUID(context.Background(), UUID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			logrus.Infof("No user found for UUID: %s", uuidStr)
			newErrorResponse(c, http.StatusNotFound, "No users found")
			return
		}
		logrus.Errorf("Error deleting user: %v", err)
		newErrorResponse(c, http.StatusInternalServerError, "Internal server error")
		return
	}

	logrus.Infof("User deleted successfully: UUID=%s", uuidStr)
	c.JSON(http.StatusOK, statusResponse{Description: "User deleted successfully"})
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
