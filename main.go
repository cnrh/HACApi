package main

import (
	"log"

	"github.com/Threqt1/HACApi/pkg/configs"
	"github.com/Threqt1/HACApi/pkg/middleware"
	"github.com/Threqt1/HACApi/pkg/routes"
	"github.com/gofiber/fiber/v2"

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

func main() {
	//Make new config
	config := configs.FiberConfig()

	//Make app with config
	app := fiber.New(config)

	//Register middleware(s)
	middleware.FiberMiddleware(app)

	//Register routes
	routes.SwaggerRoute(app)
	routes.PublicRoutes(app)
	routes.NotFoundRoute(app)

	//Start server
	log.Fatal(app.Listen(":3000"))
}
