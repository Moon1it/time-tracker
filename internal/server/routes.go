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

			users.POST("", h.CreateUser) // Add a new user
			users.GET("", h.GetAllUsers) // Get a list of users with filtering and pagination

			userID := users.Group("/:id")
			{
				userID.GET("", h.GetUser)       // Get a user data by user id
				userID.PATCH("", h.UpdateUser)  // Update a user data by user id
				userID.DELETE("", h.DeleteUser) // Delete a user by user id

				tasks := userID.Group("/tasks")
				{
					tasks.POST("/start", h.StartTimeTask)  // Start task time tracking for a user
					tasks.POST("/stop", h.StopTimeTask)    // Stop task time tracking for a user
					tasks.GET("/result", h.GetTasksResult) // Get users result for a period
				}
			}
		}
	}

	return r
}
