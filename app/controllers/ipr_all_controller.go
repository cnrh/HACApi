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
	// Parse body
	params := new(models.IprAllRequestBody)

	// Error out if fail to parse body
	if err := ctx.BodyParser(params); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err": true,
			"msg": "Bad body params",
			"ipr": nil,
		})
	}

	// Check for body param validity
	if err := server.Validator.Struct(params); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err": true,
			"msg": "Bad body params",
			"ipr": nil,
		})
	}

	// Form cache key
	cacheKey := fmt.Sprintf("%s\n%s\n%s", params.Username, params.Password, params.Base)

	// Try logging in, or grab cached collector
	collector, err := server.Cache.GetOrLogin(cacheKey)

	// Error out if login fails
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err": true,
			"msg": "Invalid username/password/base",
			"ipr": nil,
		})
	}

	// Get IPRs
	iprs, err := server.Queries.GetIPRAll(collector, params)

	// Check if returned value was nil
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": true,
			"msg": "IPRs not found. Might be an internal error",
			"ipr": nil,
		})
	}

	// All is well
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"err": false,
		"msg": "",
		"ipr": iprs,
	})
}
