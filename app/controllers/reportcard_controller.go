package controllers

import (
	"fmt"

	"github.com/Threqt1/HACApi/app/queries"
	"github.com/Threqt1/HACApi/pkg/repository"
	"github.com/Threqt1/HACApi/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

// requestCardRequestBody represents the
// request body that is to be passed in
// with the POST request to this endpoint.
type reportCardRequestBody struct {
	utils.BaseRequestBody
}

// PostReportCard handles POST requests to the report card endpoint.
// @Description Returns report card data for the user.
// @Tags        reportcard
// @Param       request body reportCardRequestBody false "Body params"
// @Accept      json
// @Produce     json
// @Success     200 {object} models.ReportCardResponse
// @Router      /reportcard [post]
func PostReportCard(server *repository.Server, ctx *fiber.Ctx) error {
	// Parse body
	params := new(reportCardRequestBody)

	// If parsing fails, error out
	if err := ctx.BodyParser(params); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":        true,
			"msg":        "Bad body params",
			"reportCard": nil,
		})
	}

	// Verify validity of body params
	if err := server.Validator.Struct(params); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":        true,
			"msg":        "Bad body params",
			"reportCard": nil,
		})
	}

	// Form cache key
	cacheKey := fmt.Sprintf("%s\n%s\n%s", params.Username, params.Password, params.Base)

	// Try logging in, or grab cached collector
	collector, err := server.Cache.GetOrLogin(cacheKey)

	// Error out if login fails
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":        true,
			"msg":        "Invalid username/password/base",
			"reportCard": nil,
		})
	}

	// Get report card
	reportCard, err := queries.GetReportCard(server, collector, params.Base)

	// Check if returned value was nil
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":        true,
			"msg":        "Report Card not found. Might be an internal error",
			"reportCard": nil,
		})
	}

	// All is well
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"err":        false,
		"msg":        "",
		"reportCard": reportCard,
	})
}
