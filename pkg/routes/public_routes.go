package routes

import (
	"github.com/Threqt1/HACApi/app/controllers"
	"github.com/Threqt1/HACApi/pkg/repository"
	"github.com/Threqt1/HACApi/pkg/utils"
)

// PublicRoutes sets up groups of public routes.
func PublicRoutes(server *repository.Server) {
	// Create group.
	route := server.App.Group("/api/v1")

	// Routes for POST methods.

	// login.
	route.Post("/login", utils.WrapController(server, controllers.PostLogin)) // post login

	// classwork.
	route.Post("/classwork", utils.WrapController(server, controllers.PostClasswork)) // post classwork

	// ipr.
	route.Post("/ipr", utils.WrapController(server, controllers.PostIPR))        // post interim progress report
	route.Post("/ipr/all", utils.WrapController(server, controllers.PostIPRAll)) // post all interim progress reports

	// report card.
	route.Post("/reportcard", utils.WrapController(server, controllers.PostReportCard)) // post report card

	// schedule.
	route.Post("/schedule", utils.WrapController(server, controllers.PostSchedule)) // post schedule

	// transcript.
	route.Post("/transcript", utils.WrapController(server, controllers.PostTranscript)) // post transcript
}
