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
	// Parse body.
	params := new(models.ScheduleRequestBody)

	// Check if parsing was successful.
	if err := ctx.BodyParser(params); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(models.ScheduleResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorBadBodyParams.Error(),
			},
		})
	}

	// Check for body parameter validity.
	if err := server.Validator.Struct(params); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(models.ScheduleResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorBadBodyParams.Error(),
			},
		})
	}

	// Form a cache key.
	cacheKey := fmt.Sprintf("%s\n%s\n%s", params.Username, params.Password, params.Base)

	// Try logging in, or grab the cached collector.
	collector, err := server.Cache.GetOrLogin(cacheKey)

	// Check if the login was successful.
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(models.ScheduleResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorInvalidAuthentication.Error(),
			},
		})
	}

	// Get the schedule.
	schedule, err := server.Querier.GetSchedule(collector, *params)

	// Check if getting the schedule succeeded.
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(models.ScheduleResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorInternalError.Error(),
			},
		})
	}

	// Return the schedule.
	return ctx.Status(fiber.StatusOK).JSON(models.ScheduleResponse{
		Schedule: schedule,
	})
}
