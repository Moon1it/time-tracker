package server

import (
	"net/http"
	"time-tracker/internal/handler"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes(h *handler.Handler) http.Handler {
	r := gin.Default()

	api := r.Group("/api")
	{
		// api.GET("/users/:id/workloads", h.GetUserWorkloads) // Получение трудозатрат по пользователю за период
		// api.POST("/users/:id/tasks/start", h.StartTask)     // Начать отсчет времени по задаче для пользователя
		// api.POST("/users/:id/tasks/stop", h.StopTask)       // Закончить отсчет времени по задаче для пользователя

		// api.GET("/info", h.GetUserByPassportNumber)

		api.POST("/users", h.CreateUser)         // Добавление нового пользователя
		api.GET("/users", h.GetUsers)            // Получение данных пользователей с фильтрацией и пагинацией
		api.GET("/users/:uuid", h.GetUser)       // Добавление нового пользователя
		api.PATCH("/users/:uuid", h.UpdateUser)  // Изменение данных пользователя
		api.DELETE("/users/:uuid", h.DeleteUser) // Удаление пользователя

	}

	return r
}

// api.GET("/users", h.GetUsers)                         // Получение списка пользователей с фильтрацией и пагинацией
// api.GET("/user/time-tracking", h.GetUserTimeTracking) // Получение трудозатрат по пользователю за период
// api.POST("/user/start-task", h.StartTask)             // Начало отсчета времени по задаче для пользователя
// api.POST("/user/stop-task", h.StopTask)               // Завершение отсчета времени по задаче для пользователя
