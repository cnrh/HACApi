package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/Threqt1/HACApi/pkg/configs"
	"github.com/Threqt1/HACApi/pkg/middleware"
	"github.com/Threqt1/HACApi/pkg/routes"
	"github.com/Threqt1/HACApi/pkg/utils"
	"github.com/joho/godotenv"

	_ "github.com/Threqt1/HACApi/docs" // load API Docs files (Swagger)
)

//	@title				HAC Information API
//	@version			1.0
//	@description		An API to fetch data from Home Access Center.
//	@BasePath			/api/v1
//	@accept				json
//	@produce			json
//
//	@tag.name			auth
//	@tag.description	Caching a login with the API
//
//	@tag.name			classwork
//	@tag.description	Get data about classwork
//
//	@tag.name			ipr
//	@tag.description	Get data about interim progress report(s)
//
//	@tag.name			reportcard
//	@tag.description	Get data about the report card
//
//	@tag.name			schedule
//	@tag.description	Get data about the schedule
//
//	@tag.name			transcript
//	@tag.description	Get data about the transcript

func main() {
	// Register .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("DotEnv failed to load. Error: %v", err)
	}

	// Make new server
	server := configs.ServerConfig()

	// Register middleware(s)
	middleware.FiberMiddleware(server)

	// Register routes
	routes.SwaggerRoute(server)
	routes.PublicRoutes(server)
	routes.NotFoundRoute(server)

	// Start server

	// Create channel to confirm when connections are closed
	connsClosedChan := make(chan struct{})

	go func() {
		// Catch os signals
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// Gracefully shutdown
		if err := server.App.Shutdown(); err != nil {
			log.Fatalf("Server failed to shutdown. Reason: %v", err)
		}

		close(connsClosedChan)
	}()

	// Build fiber URL
	fiberConnURL, _ := utils.BuildConnectionURL("fiber")

	// Start server
	err = server.App.Listen(fiberConnURL)

	if err != nil {
		log.Fatalf(err.Error())
	}

	// Wait till conns are closed before stopping
	<-connsClosedChan
}
