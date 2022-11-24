package controllers

import (
	"fmt"

	"github.com/Threqt1/HACApi/app/queries"
	"github.com/Threqt1/HACApi/platform/cache"
	"github.com/gofiber/fiber/v2"
)

type scheduleRequestBody struct {
	baseRequestBody
}

// PostSchedule handles POST request to the schedule endpoint.
// @Description Returns the schedule for the user.
// @Tags        schedule
// @Param       request body scheduleRequestBody false "Body params"
// @Accept      json
// @Produce     json
// @Success     200 {object} models.ScheduleResponse
// @Router      /schedule [post]
func PostSchedule(ctx *fiber.Ctx) error {
	//Parse body
	params := new(scheduleRequestBody)

	//If parsing fails, error out
	if err := ctx.BodyParser(params); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":      true,
			"msg":      "Bad body params",
			"schedule": nil,
		})
	}

	//Check for body param validity
	bodyParamsValid := true

	//Confirm no required body parameters are empty
	if params.Username == "" || params.Password == "" || params.Base == "" {
		bodyParamsValid = false
	}

	//If body params not valid, return error
	if !bodyParamsValid {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":      true,
			"msg":      "Bad body params",
			"schedule": nil,
		})
	}

	//Form cache key
	cacheKey := fmt.Sprintf("%s\n%s\n%s", params.Username, params.Password, params.Base)

	//Try logging in, or grab cached collector
	collector := cache.CurrentCache().Get(cacheKey)

	//Error out if login fails
	if collector == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":      true,
			"msg":      "Invalid username/password/base",
			"schedule": nil,
		})
	}

	//Get schedule
	schedule, err := queries.GetSchedule(collector.Value(), params.Base)

	//Check if returned value was nil
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"err":      true,
			"msg":      "Schedule not found. Might be an internal error",
			"schedule": nil,
		})
	}

	//All is well
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"err":      false,
		"msg":      "",
		"schedule": schedule,
	})
}
