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
	// Parse body.
	params := new(models.ClassworkRequestBody)

	// Check if parsing body parameters succeeded.
	if err := ctx.BodyParser(params); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(models.ClassworkResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorBadBodyParams.Error(),
			},
		})
	}

	// Verify the validity of the body params.
	if err := server.Validator.Struct(params); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(models.ClassworkResponse{
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

	// Error out if the login fails.
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(models.ClassworkResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorInvalidAuthentication.Error(),
			},
		})
	}

	// Get the classwork.
	classwork, err := server.Querier.GetClasswork(collector, *params)

	// Check if returned value was nil, and if so error out.
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(models.ClassworkResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorInternalError.Error(),
			},
		})
	}

	// Return the recieved classwork.
	return ctx.Status(fiber.StatusOK).JSON(models.ClassworkResponse{
		Classwork: classwork,
	})
}
