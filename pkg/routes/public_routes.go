package routes

import (
	"github.com/Threqt1/HACApi/app/controllers"
	"github.com/gofiber/fiber/v2"
)

// PublicRoutes sets up groups of public routes.
func PublicRoutes(app *fiber.App) {
	//Create group
	route := app.Group("/api/v1")

	//Routes for POST methods
	route.Post("/login", controllers.PostLogin)
	route.Post("/classwork", controllers.PostClasswork)   //post classwork
	route.Post("/ipr", controllers.PostIPR)               //post interim progress report
	route.Post("/reportcard", controllers.PostReportCard) //post report card
	route.Post("/schedule", controllers.PostSchedule)     //post schedule
	route.Post("/transcript", controllers.PostTranscript) //post transcript
}
