package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type TaskResponse struct {
	Name    string
	Surname string
}

func (h *Handler) CreateTask(c *gin.Context) {
	c.JSON(http.StatusOK, &TaskResponse{
		Name:    "Kirill",
		Surname: "Smirnov",
	})
}

func (h *Handler) GetTasks(c *gin.Context) {
	c.JSON(http.StatusOK, &TaskResponse{
		Name:    "Kirill",
		Surname: "Smirnov",
	})
}

func (h *Handler) GetTask(c *gin.Context) {
	c.JSON(http.StatusOK, &TaskResponse{
		Name:    "Kirill",
		Surname: "Smirnov",
	})
}

func (h *Handler) UpdateTask(c *gin.Context) {
	c.JSON(http.StatusOK, &TaskResponse{
		Name:    "Kirill",
		Surname: "Smirnov",
	})
}

func (h *Handler) DeleteTask(c *gin.Context) {
	c.JSON(http.StatusOK, &TaskResponse{
		Name:    "Kirill",
		Surname: "Smirnov",
	})
}
