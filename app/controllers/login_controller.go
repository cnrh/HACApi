package controllers

import (
	"fmt"

	"github.com/Threqt1/HACApi/pkg/repository"
	"github.com/Threqt1/HACApi/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

// loginRequestBody represents the body that is to be
// passed along with the POST request to the login
// endpoint.
type loginRequestBody struct {
	utils.BaseRequestBody
}

// PostLogin handles POST requests to the login endpoint.
// @Description Pre-registers the user with the API by logging them into HAC early, and caching the cookies.
// @Description Subsequent requests using the same credentials will use these stored cookies, leading to faster response times for other endpoints.
// @Tags        auth
// @Param       request body loginRequestBody false "Body Params"
// @Accept      json
// @Produce     json
// @Success     200 {object} models.LoginResponse
// @Router      /login [post]
func PostLogin(server *repository.Server, ctx *fiber.Ctx) error {
	// Parse body
	params := new(loginRequestBody)

	// If parsing body params failed, return error
	if err := ctx.BodyParser(params); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err": true,
			"msg": "Bad body params",
		})
	}

	// Verify validity of body params
	if err := server.Validator.Struct(params); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err": true,
			"msg": "Bad body params",
		})
	}

	// Form cache key
	cacheKey := fmt.Sprintf("%s\n%s\n%s", params.Username, params.Password, params.Base)

	// Cache the user, if not cached already
	_, err := server.Cache.GetOrLogin(cacheKey)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err": true,
			"msg": "Invalid username/password",
		})
	}

	// Success
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"err": false,
		"msg": "",
	})
}
