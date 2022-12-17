package routes

import (
	"github.com/Threqt1/HACApi/pkg/repository"
	"github.com/gofiber/swagger"
)

// SwaggerRoute sets up the swagger docs endpoint.
func SwaggerRoute(server *repository.Server) {
	// Create swagger route (docs).
	route := server.App.Group("docs")

	// Let swagger handle routes.
	route.Get("*", swagger.New(swagger.Config{
		CustomStyle: "../../../docs/swagger-flattop-theme.css",
	})).Name("docs")
}
