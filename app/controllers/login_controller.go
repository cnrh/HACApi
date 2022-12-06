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
	// Parse body
	params := new(models.LoginRequestBody)

	// If parsing body params failed, return error
	if err := ctx.BodyParser(params); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":   true,
			"msg":   "Bad body params",
			"login": nil,
		})
	}

	// Verify validity of body params
	if err := server.Validator.Struct(params); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":   true,
			"msg":   "Bad body params",
			"login": nil,
		})
	}

	// Form cache key
	cacheKey := fmt.Sprintf("%s\n%s\n%s", params.Username, params.Password, params.Base)

	// Cache the user, if not cached already
	collector, err := server.Cache.GetOrLogin(cacheKey)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":   true,
			"msg":   "Invalid username/password",
			"login": nil,
		})
	}

	// Get response
	login, err := server.Queries.GetLogin(collector, params)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":   true,
			"msg":   "Invalid username/password",
			"login": nil,
		})
	}

	// Success
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"err":   false,
		"msg":   "",
		"login": login,
	})
}
