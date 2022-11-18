package routes

import (
	"github.com/Threqt1/HACApi/app/controllers"
	"github.com/gofiber/fiber/v2"
)

// PublicRoutes func for describing groups of public routes
func PublicRoutes(app *fiber.App) {
	//Create group
	route := app.Group("/api/v1")

	//Routes for GET methods
	route.Get("/classwork", controllers.GetClasswork)
}
