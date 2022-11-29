package controllers

import (
	"fmt"

	"github.com/Threqt1/HACApi/app/queries"
	"github.com/Threqt1/HACApi/pkg/utils"
	"github.com/Threqt1/HACApi/platform/cache"
	"github.com/gofiber/fiber/v2"
)

// iprAllRequestBody represents the body that is to
// be passed with a POST request to the /ipr/all
// endpoint.
type iprAllRequestBody struct {
	utils.BaseRequestBody
	//Whether to return only dates or all the IPRs
	DatesOnly bool `json:"datesOnly" validate:"optional" example:"true" default:"false"`
}

// PostIPRAll handles POST requests to the IPR/All endpoint.
// @Description Returns all the IPRs for the user, or just the dates depending on the DatesOnly parameter's value in the body.
// @Tags        ipr
// @Param       request body iprAllRequestBody false "Body Params"
// @Accept      json
// @Produce     json
// @Success     200 {object} models.IPRResponse
// @Router      /ipr/all [post]
func PostIPRAll(ctx *fiber.Ctx) error {
	//Parse body
	params := new(iprAllRequestBody)

	//Error out if fail to parse body
	if err := ctx.BodyParser(params); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err": true,
			"msg": "Bad body params",
			"ipr": nil,
		})
	}

	//Check for body param validity
	bodyParamsValid := true

	//Confirm no required body params are empty
	if params.Username == "" || params.Password == "" || params.Base == "" {
		bodyParamsValid = false
	}

	//If body params are invalid, error out
	if !bodyParamsValid {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err": true,
			"msg": "Bad body params",
			"ipr": nil,
		})
	}

	//Form cache key
	cacheKey := fmt.Sprintf("%s\n%s\n%s", params.Username, params.Password, params.Base)

	//Try logging in, or grab cached collector
	collector := cache.CurrentCache().Get(cacheKey)

	//Error out if login fails
	if collector == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err": true,
			"msg": "Invalid username/password/base",
			"ipr": nil,
		})
	}

	//Get IPRs
	iprs, err := queries.GetAllIPRs(collector.Value(), params.Base, params.DatesOnly)

	//Check if returned value was nil
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"err": true,
			"msg": "IPRs not found. Might be an internal error",
			"ipr": nil,
		})
	}

	//All is well
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"err": false,
		"msg": "",
		"ipr": iprs,
	})
}
