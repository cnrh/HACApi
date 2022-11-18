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
	route.Get("/classwork", controllers.GetClasswork)   //get classwork
	route.Get("/ipr", controllers.GetIPR)               //get interim progress report
	route.Get("/reportcard", controllers.GetReportCard) //get report card
	route.Get("/schedule", controllers.GetSchedule)     //get schedule
	route.Get("/transcript", controllers.GetTranscript) //get transcript
}
