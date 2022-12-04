package main

import (
	"log"
	"os"

	"github.com/Threqt1/HACApi/pkg/configs"
	"github.com/Threqt1/HACApi/pkg/middleware"
	"github.com/Threqt1/HACApi/pkg/routes"
	"github.com/Threqt1/HACApi/pkg/utils"
	"github.com/joho/godotenv"

	_ "github.com/Threqt1/HACApi/docs" // load API Docs files (Swagger)
)

// @title           HAC Information API
// @version         1.0
// @description     An API to fetch data from Home Access Center.
// @BasePath        /api/v1
// @accept          json
// @produce         json
//
// @tag.name        auth
// @tag.description Caching a login with the API
//
// @tag.name        classwork
// @tag.description Get data about classwork
//
// @tag.name        ipr
// @tag.description Get data about interim progress report(s)
//
// @tag.name        reportcard
// @tag.description Get data about the report card
//
// @tag.name        schedule
// @tag.description Get data about the schedule
//
// @tag.name        transcript
// @tag.description Get data about the transcript

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
	if os.Getenv("DEV_STAGE") == "dev" {
		utils.StartForDev(server)
	} else {
		utils.StartForProd(server)
	}
}
