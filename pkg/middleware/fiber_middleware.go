package middleware

import (
	"github.com/Threqt1/HACApi/pkg/repository"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// FiberMiddleware sets up fiber's middleware for
// the API.
func FiberMiddleware(server *repository.Server) {
	server.App.Use(
		// Enable CORS
		cors.New(),

		// Add a logger
		logger.New(),
	)
}
