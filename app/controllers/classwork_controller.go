package controllers

import (
	"fmt"

	"github.com/Threqt1/HACApi/app/models"
	"github.com/Threqt1/HACApi/pkg/repository"
	"github.com/gofiber/fiber/v2"
)

// PostClasswork handles POST requests to the classwork endpoint.
//
//	@Description	Returns classwork for the marking periods specified.
//	@Description	If no marking periods are specified, the classwork for the current marking period is returned.
//	@Tags			classwork
//	@Param			request	body	models.ClassworkRequestBody	false	"Body Params"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.ClassworkResponse
//	@Router			/classwork [post]
func PostClasswork(server *repository.Server, ctx *fiber.Ctx) error {
	// Parse body
	params := new(models.ClassworkRequestBody)

	// If parsing body params failed, return error
	if err := ctx.BodyParser(params); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":       true,
			"msg":       "Bad body params",
			"classwork": nil,
		})
	}

	// Verify validity of body params
	if err := server.Validator.Struct(params); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":       true,
			"msg":       "Bad body params",
			"classwork": nil,
		})
	}

	// Form cache key
	cacheKey := fmt.Sprintf("%s\n%s\n%s", params.Username, params.Password, params.Base)

	// Try logging in, or grab cached collector
	collector, err := server.Cache.GetOrLogin(cacheKey)

	// Error out if login fails
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":       true,
			"msg":       "Invalid username/password/base",
			"classwork": nil,
		})
	}

	// Get classwork
	classwork, err := server.Querier.GetClasswork(collector, *params)

	// Check if returned value was nil
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":       true,
			"msg":       "Classwork not found. Might be an internal error",
			"classwork": nil,
		})
	}

	// All is well
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"err":       false,
		"msg":       "",
		"classwork": classwork,
	})
}
