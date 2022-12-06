package utils

import (
	"github.com/Threqt1/HACApi/pkg/repository"
	"github.com/gofiber/fiber/v2"
)

// WrapController wraps an endpoint controller method into a function acceptable by Fiber.
func WrapController(server *repository.Server, wrapped func(*repository.Server, *fiber.Ctx) error) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		return wrapped(server, ctx)
	}
}
