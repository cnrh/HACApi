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

	//login
	route.Post("/login", controllers.PostLogin) //register login with the API

	//classwork
	route.Post("/classwork", controllers.PostClasswork) //post classwork

	//ipr
	route.Post("/ipr", controllers.PostIPR)        //post interim progress report
	route.Post("/ipr/all", controllers.PostIPRAll) //post all interim progress reports

	//report card
	route.Post("/reportcard", controllers.PostReportCard) //post report card

	//schedule
	route.Post("/schedule", controllers.PostSchedule) //post schedule

	//transcript
	route.Post("/transcript", controllers.PostTranscript) //post transcript
}
