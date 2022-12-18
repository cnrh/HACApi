package controllers

import (
	"fmt"

	"github.com/Threqt1/HACApi/app/models"
	"github.com/Threqt1/HACApi/pkg/repository"
	"github.com/gofiber/fiber/v2"
)

// PostLogin handles POST requests to the login endpoint.
//
//	@Description	Pre-registers the user with the API by logging them into HAC early, and caching the cookies.
//	@Description	Subsequent requests using the same credentials will use these stored cookies, leading to faster response times for other endpoints.
//	@Tags			auth
//	@Param			request	body	models.LoginRequestBody	false	"Body Params"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.LoginResponse
//	@Router			/login [post]
func PostLogin(server *repository.Server, ctx *fiber.Ctx) error {
	// Parse body.
	params := new(models.LoginRequestBody)

	// Check if parsing was successful.
	if err := ctx.BodyParser(params); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(models.LoginResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorBadBodyParams.Error(),
			},
		})
	}

	// Verify the validity of the body parameters.
	if err := server.Validator.Struct(params); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(models.LoginResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorBadBodyParams.Error(),
			},
		})
	}

	// Form a cache key.
	cacheKey := fmt.Sprintf("%s\n%s\n%s", params.Username, params.Password, params.Base)

	// Cache the user, if not cached already.
	collector, err := server.Cache.GetOrLogin(cacheKey)

	// Check if caching succeeded.
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(models.LoginResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorInvalidAuthentication.Error(),
			},
		})
	}

	// Get response from the querier.
	login, err := server.Querier.GetLogin(collector, *params)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(models.LoginResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorInternalError.Error(),
			},
		})
	}

	// Send back information about the login.
	return ctx.Status(fiber.StatusOK).JSON(models.LoginResponse{
		Login: login,
	})
}
