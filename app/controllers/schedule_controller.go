package controllers

import (
	"fmt"

	"github.com/Threqt1/HACApi/app/models"
	"github.com/Threqt1/HACApi/pkg/repository"
	"github.com/gofiber/fiber/v2"
)

// PostSchedule handles POST request to the schedule endpoint.
//
//	@Description	Returns the schedule for the user.
//	@Tags			schedule
//	@Param			request	body	models.ScheduleRequestBody	false	"Body params"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.ScheduleResponse
//	@Router			/schedule [post]
func PostSchedule(server *repository.Server, ctx *fiber.Ctx) error {
	// Parse body
	params := new(models.ScheduleRequestBody)

	// If parsing fails, error out
	if err := ctx.BodyParser(params); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":      true,
			"msg":      "Bad body params",
			"schedule": nil,
		})
	}

	// Check for body param validity
	if err := server.Validator.Struct(params); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":      true,
			"msg":      "Bad body params",
			"schedule": nil,
		})
	}

	// Form cache key
	cacheKey := fmt.Sprintf("%s\n%s\n%s", params.Username, params.Password, params.Base)

	// Try logging in, or grab cached collector
	collector, err := server.Cache.GetOrLogin(cacheKey)

	// Error out if login fails
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":      true,
			"msg":      "Invalid username/password/base",
			"schedule": nil,
		})
	}

	// Get schedule
	schedule, err := server.Queries.GetSchedule(collector, params)

	// Check if returned value was nil
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":      true,
			"msg":      "Schedule not found. Might be an internal error",
			"schedule": nil,
		})
	}

	// All is well
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"err":      false,
		"msg":      "",
		"schedule": schedule,
	})
}
