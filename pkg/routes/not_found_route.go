package routes

import (
	"github.com/Threqt1/HACApi/pkg/repository"
	"github.com/gofiber/fiber/v2"
)

// Send an error for any route not registered
func NotFoundRoute(server *repository.Server) {
	server.App.Use(func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"err": true,
			"msg": "No endpoint found",
		})
	})
}
