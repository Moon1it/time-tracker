package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time-tracker/internal/models"
	"time-tracker/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var validPeriods = map[string]interface{}{
	"day":   nil,
	"week":  nil,
	"month": nil,
	"year":  nil,
}

// @Summary      Start a time task
// @Description  Create a new task for a user
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id     path      string                        true  "User id"
// @Param        payload  body      models.CreateTaskPayload      true  "Task Payload"
// @Success      201      {object}  models.Task                   "Task created successfully"
// @Failure      400      {object}  errorResponse                 "Bad request"
// @Failure      409      {object}  errorResponse                 "Task with this user id already exists. Please complete the active task first."
// @Failure      500      {object}  errorResponse                 "Internal server error"
// @Router       /users/{id}/tasks/start [post]
func (h *Handler) StartTimeTask(c *gin.Context) {
	var payload models.CreateTaskPayload
	if err := c.BindJSON(&payload); err != nil {
		logrus.Errorf("Invalid JSON: %v", err)
		newErrorResponse(c, http.StatusBadRequest, "Bad request")
		return
	}

	userIDParam := c.Param("id")
	userUUID, err := uuid.Parse(userIDParam)
	if err != nil {
		logrus.Errorf("Invalid UUID format: %v", err)
		newErrorResponse(c, http.StatusBadRequest, "Bad request")
		return
	}

	ctx := c.Request.Context()
	task, err := h.service.ITaskService.CreateTask(ctx, userUUID, &payload)
	if err != nil {
		if errors.Is(err, service.ErrForeignKeyViolation) {
			logrus.Warnf("Error creating task: %v", err)
			newErrorResponse(c, http.StatusBadRequest, "Bad request")
		}
		if errors.Is(err, service.ErrTaskAlreadyExists) {
			logrus.Warnf("Task with user UUID %s already exists: %v", userUUID, err)
			newErrorResponse(c, http.StatusConflict, "Task with this user UUID already exists. Please complete the active task first.")
		}
		logrus.Errorf("Error starting task: %v", err)
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	logrus.Infof("Task created successfully: %v", task)
	c.JSON(http.StatusCreated, task)
}

// @Summary      Stop a time task
// @Description  Stop an active task for a user
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id     path      string                        true  "User id"
// @Success      200      {object}  models.Task                   "Task stopped successfully"
// @Failure      400      {object}  errorResponse                 "Bad request"
// @Failure      404      {object}  errorResponse                 "No users found or this user does not have an active task yet."
// @Failure      500      {object}  errorResponse                 "Internal server error"
// @Router /users/{id}/tasks/stop [post]
func (h *Handler) StopTimeTask(c *gin.Context) {
	userIDParam := c.Param("id")
	userUUID, err := uuid.Parse(userIDParam)
	if err != nil {
		logrus.Errorf("Invalid UUID format: %v", err)
		newErrorResponse(c, http.StatusBadRequest, "Bad request")
		return
	}

	ctx := c.Request.Context()
	task, err := h.service.ITaskService.FinishTask(ctx, userUUID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			logrus.Infof("No user found for UUID: %s", userUUID)
			newErrorResponse(c, http.StatusNotFound, "No users found")
			return
		}
		if errors.Is(err, service.ErrTaskNotFound) {
			logrus.Warnf("No active task for user UUID %s: %v", userUUID, err)
			newErrorResponse(c, http.StatusNotFound, "This user does not have an active task yet.")
			return
		}
		logrus.Errorf("Error finishing task: %v", err)
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	logrus.Infof("Task stopped successfully for user UUID: %s", userUUID)
	c.JSON(http.StatusOK, task)
}

// @Summary Get tasks result
// @Description Retrieve tasks result for a user within a specified time period
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param id path string true "User id"
// @Param timePeriod query string false "Time period ('day', 'week', 'month', 'year')" default(day)
// @Param timeAmount query string false "Amount of time" default(1)
// @Success 200 {object} models.TasksResult "Tasks retrieved successfully"
// @Success 204 {object} nil "No tasks found for the specified period"
// @Failure 400 {object} errorResponse "Bad request"
// @Failure 404 {object} errorResponse "User not found"
// @Failure 500 {object} errorResponse "Internal server error"
// @Router /users/{id}/tasks/result [get]
func (h *Handler) GetTasksResult(c *gin.Context) {
	userIDParam := c.Param("id")
	userUUID, err := uuid.Parse(userIDParam)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Bad request")
		return
	}

	timePeriod := c.DefaultQuery("timePeriod", "day")
	timeAmount := c.DefaultQuery("timeAmount", "1")

	if _, isValid := validPeriods[timePeriod]; !isValid {
		newErrorResponse(c, http.StatusBadRequest, "Bad request")
		return
	}

	days, err := convertPeriodToDays(timePeriod, timeAmount)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Bad request")
		return
	}

	ctx := c.Request.Context()
	task, err := h.service.ITaskService.GetTasksResult(ctx, userUUID, days)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			logrus.Infof("No user found for UUID: %s", userUUID)
			newErrorResponse(c, http.StatusNotFound, "User not found")
			return
		}
		if errors.Is(err, service.ErrTaskNotFound) {
			logrus.Infof("No tasks found for the specified period for user UUID: %s", userUUID)
			c.JSON(http.StatusNoContent, nil)
			return
		}
		logrus.Errorf("Error getting tasks for user UUID: %s", userUUID)
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	c.JSON(http.StatusOK, task)
}

func convertPeriodToDays(period string, amount string) (int, error) {
	amountInt, err := strconv.Atoi(amount)
	if err != nil {
		return 0, fmt.Errorf("invalid time amount")
	}

	switch period {
	case "day":
		return amountInt, nil
	case "week":
		return amountInt * 7, nil
	case "month":
		return amountInt * 30, nil
	case "year":
		return amountInt * 365, nil
	default:
		return 0, fmt.Errorf("invalid period value")
	}
}
