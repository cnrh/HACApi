package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// FiberMiddleware sets up fiber's middleware for
// the API.
func FiberMiddleware(app *fiber.App) {
	app.Use(
		// Enable CORS
		cors.New(),

		// Add a logger
		logger.New(),
	)
}
