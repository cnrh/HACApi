package routes

import (
	"github.com/Threqt1/HACApi/pkg/repository"
	"github.com/gofiber/fiber/v2"
)

// SwaggerRoute sets up the swagger docs endpoint.
func SwaggerRoute(server *repository.Server) {
	// Create swagger route (docs).
	route := server.App.Group("docs")

	// Redirect to the GitHub Pages link.
	route.Get("*", func(c *fiber.Ctx) error {
		return c.Redirect("https://threqt1.github.io/HACApi/", fiber.StatusSeeOther)
	})

	// route.Get("*", swagger.New(swagger.Config{
	// 	CustomStyle: "../../../docs/swagger-flattop-theme.css",
	// })).Name("docs")
}
