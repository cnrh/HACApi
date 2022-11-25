package controllers

///ipr/recent for recent, /ipr/all for all (option for just dates), /ipr for specific

import (
	"fmt"
	"time"

	"github.com/Threqt1/HACApi/app/queries"
	"github.com/Threqt1/HACApi/pkg/utils"
	"github.com/Threqt1/HACApi/platform/cache"
	"github.com/gofiber/fiber/v2"
)

// iprRequestBody represents the body that is to
// be passed with a POST request to the IPR
// endpoint.
type iprRequestBody struct {
	utils.BaseRequestBody
	//The date of the IPR to return
	Date string `json:"date" validate:"optional" example:"09/06/2022"`
}

// PostIPR handles POST requests to the IPR endpoint.
// @Description Returns the IPR(s) for the user. If the date parameter is not passed into the body or is invalid, the most recent IPR is returned.
// @Description It is important the format of the date follows the format "01/02/2006" (01 = month, 02 = day, 2006 = year), with leading zeros like shown in the format.
// @Description For all possible dates, refer to the "/ipr/all" endpoint.
// @Tags        ipr
// @Param       request body iprRequestBody false "Body Params"
// @Accept      json
// @Produce     json
// @Success     200 {object} models.IPRResponse
// @Router      /ipr [post]
func PostIPR(ctx *fiber.Ctx) error {
	//Parse body
	params := new(iprRequestBody)

	//If parsing body fails, error out
	if err := ctx.BodyParser(params); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err": true,
			"msg": "Bad body params",
			"ipr": nil,
		})
	}

	//Verify validity of body params
	bodyParamsValid := true

	//Store parsed date
	var parsedDate time.Time

	//Confirm no required body params are empty
	if params.Username == "" || params.Password == "" || params.Base == "" {
		bodyParamsValid = false
	}

	//Check for valid date
	date, err := time.Parse("01/02/2006", params.Date)

	if err == nil {
		parsedDate = date
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

	//Get IPR
	ipr, err := queries.GetIPR(collector.Value(), params.Base, parsedDate)

	//Check if returned value was nil
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"err": true,
			"msg": "IPR not found. Might be an internal error",
			"ipr": nil,
		})
	}

	//All is well
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"err": false,
		"msg": "",
		"ipr": ipr,
	})
}

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
