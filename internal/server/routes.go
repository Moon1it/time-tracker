package server

import (
	"net/http"
	"time-tracker/internal/handlers"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes(h *handlers.Handler) http.Handler {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/", h.HelloWorldHandler)
	}

	return r
}
