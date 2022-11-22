package controllers

import (
	"fmt"

	"github.com/Threqt1/HACApi/app/queries"
	"github.com/Threqt1/HACApi/platform/cache"
	"github.com/gofiber/fiber/v2"
)

// classworkRequestBody represents the body that is to be passed with
// the POST request to the classwork endpoint.
type classworkRequestBody struct {
	//The username to log in with
	Username string `json:"username" validate:"required" example:"j1732901"`
	//The password to log in with
	Password string `json:"password" validate:"required" example:"j382704"`
	//The base URL for the PowerSchool HAC service
	Base string `json:"base" validate:"required" example:"homeaccess.katyisd.org"`
	//The marking period to pull data from
	MarkingPeriods []int `json:"markingPeriods" validate:"optional" example:"1,2"`
}

// PostClasswork handles POST requests to the classwork endpoint.
// @Description Returns classwork for the user.
// @Tags        data
// @Param       request body classworkRequestBody false "Body Params"
// @Accept      json
// @Produce     json
// @Success     200 {object} models.ClassworkResponse
// @Router      /classwork [post]
func PostClasswork(ctx *fiber.Ctx) error {
	//Parse body
	params := new(classworkRequestBody)

	//If parsing body params failed, return error
	if err := ctx.BodyParser(params); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":       true,
			"msg":       "Bad body params",
			"classwork": nil,
		})
	}

	//Verify validity of body params
	bodyParamsValid := true

	//Confirm no required body parameters are empty
	if params.Username == "" || params.Password == "" || params.Base == "" {
		bodyParamsValid = false
	}

	//Confirm length of marking periods array is 0 <= len <= 6
	if mpLen := len(params.MarkingPeriods); mpLen > 6 {
		bodyParamsValid = false
	}

	//Confirm marking periods array has values of 1 <= val <= 6
	for _, mp := range params.MarkingPeriods {
		if mp < 1 || mp > 6 {
			bodyParamsValid = false
			break
		}
	}

	//If body params not valid, return error
	if !bodyParamsValid {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":       true,
			"msg":       "Bad body params",
			"classwork": nil,
		})
	}

	//Form cache key
	cacheKey := fmt.Sprintf("%s\n%s\n%s", params.Username, params.Password, params.Base)

	//Try logging in, or grab cached collector
	collector := cache.CurrentCache().Get(cacheKey)

	if collector == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":       true,
			"msg":       "Invalid username/password/base",
			"classwork": nil,
		})
	}

	//Get classwork
	classwork := queries.GetClasswork(collector.Value(), params.Base, params.MarkingPeriods)

	//Check if returned value was nil
	if classwork == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"err":       true,
			"msg":       "Classwork not found. Might be an internal error.",
			"classwork": nil,
		})
	}

	//All is well
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"err":       false,
		"msg":       "",
		"classwork": classwork,
	})
}
