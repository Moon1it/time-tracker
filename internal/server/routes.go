package server

import (
	"net/http"
	"time-tracker/internal/handler"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) RegisterRoutes(h *handler.Handler) http.Handler {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	{
		users := api.Group("/users")
		{
			users.GET("/info", h.GetUserByPassportNumber) // Check for user existence by passportNumber

			users.POST("", h.CreateUser)         // Add a new user
			users.GET("", h.GetAllUsers)         // Get a list of users with filtering and pagination
			users.GET("/:uuid", h.GetUser)       // Get user data by UUID
			users.PATCH("/:uuid", h.UpdateUser)  // Update user data
			users.DELETE("/:uuid", h.DeleteUser) // Delete a user

			tasks := users.Group("/:uuid/tasks")
			{
				tasks.POST("/start", h.StartTimeTask)  // Start task time tracking for a user
				tasks.POST("/stop", h.StopTimeTask)    // Stop task time tracking for a user
				tasks.GET("/result", h.GetTasksResult) // Get users result for a period
			}
		}
	}

	return r
}
