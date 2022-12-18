package controllers

import (
	"fmt"

	"github.com/Threqt1/HACApi/app/models"
	"github.com/Threqt1/HACApi/pkg/repository"
	"github.com/gofiber/fiber/v2"
)

// PostIPRAll handles POST requests to the IPR/All endpoint.
//
//	@Description	Returns all the IPRs for the user, or just the dates depending on the DatesOnly parameter's value in the body.
//	@Tags			ipr
//	@Param			request	body	models.IprAllRequestBody	false	"Body Params"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.IPRResponse
//	@Router			/ipr/all [post]
func PostIPRAll(server *repository.Server, ctx *fiber.Ctx) error {
	// Parse body.
	params := new(models.IprAllRequestBody)

	// Check if parsing body was successful.
	if err := ctx.BodyParser(params); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(models.IPRResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorBadBodyParams.Error(),
			},
		})
	}

	// Check if the body parameters are valid.
	if err := server.Validator.Struct(params); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(models.IPRResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorBadBodyParams.Error(),
			},
		})
	}

	// Form a cache key.
	cacheKey := fmt.Sprintf("%s\n%s\n%s", params.Username, params.Password, params.Base)

	// Try logging in, or grab the cached collector
	collector, err := server.Cache.GetOrLogin(cacheKey)

	// Check if the login failed.
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(models.IPRResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorInvalidAuthentication.Error(),
			},
		})
	}

	// Get IPRs.
	iprs, err := server.Querier.GetIPRAll(collector, *params)

	// Check if getting IPRs succeeded.
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(models.IPRResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorInternalError.Error(),
			},
		})
	}

	// Return the grabbed IPRs.
	return ctx.Status(fiber.StatusOK).JSON(models.IPRResponse{
		IPR: iprs,
	})
}
