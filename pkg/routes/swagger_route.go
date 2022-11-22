package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// SwaggerRoute sets up the swagger docs endpoint.
func SwaggerRoute(app *fiber.App) {
	// Create swagger route (docs)
	route := app.Group("docs")

	// Let swagger handle routes
	route.Get("*", swagger.HandlerDefault)
}
