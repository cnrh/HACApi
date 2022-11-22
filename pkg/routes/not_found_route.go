package routes

import "github.com/gofiber/fiber/v2"

//Send an error for any route not registered
func NotFoundRoute(app *fiber.App) {
	app.Use(func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"err": true,
			"msg": "No endpoint found",
		})
	})
}
