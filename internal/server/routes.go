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
		api.POST("/people", h.CreatePeople)
		api.GET("/info", h.GetPeople)

		// users := api.Group("/users")
		// {
		// 	users.POST("/", h.CreateUser)
		// 	users.GET("/", h.GetUsers)
		// 	users.GET("/:id", h.GetUser)
		// 	users.PUT("/:id", h.UpdateUser)
		// 	users.DELETE("/:id", h.DeleteUser)

		// 	// tasks.GET("/:id/tasks/stop", h.HelloWorldHandler)
		// }
		// // :id/execution-times

		// tasks := users.Group("/tasks")
		// {

		// 	tasks.POST("/", h.CreateTask)
		// 	tasks.GET("/", h.HelloWorldHandler)
		// 	tasks.GET("/:id", h.HelloWorldHandler)
		// 	tasks.DELETE("/:id", h.HelloWorldHandler)
		// }
		api.GET("/", h.HelloWorldHandler)
	}

	return r
}
