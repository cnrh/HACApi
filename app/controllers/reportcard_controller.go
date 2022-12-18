package controllers

import (
	"fmt"

	"github.com/Threqt1/HACApi/app/models"
	"github.com/Threqt1/HACApi/pkg/repository"
	"github.com/gofiber/fiber/v2"
)

// PostReportCard handles POST requests to the report card endpoint.
//
//	@Description	Returns report card data for the user.
//	@Tags			reportcard
//	@Param			request	body	models.ReportCardRequestBody	false	"Body params"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.ReportCardResponse
//	@Router			/reportcard [post]
func PostReportCard(server *repository.Server, ctx *fiber.Ctx) error {
	// Parse body.
	params := new(models.ReportCardRequestBody)

	// Check if the body was parsed successfully.
	if err := ctx.BodyParser(params); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(models.ReportCardResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorBadBodyParams.Error(),
			},
		})
	}

	// Verify the validity of body the parameters.
	if err := server.Validator.Struct(params); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(models.ReportCardResponse{
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

	// Check if the login was successful.
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(models.ReportCardResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorInvalidAuthentication.Error(),
			},
		})
	}

	// Get the report card.
	reportCard, err := server.Querier.GetReportCard(collector, *params)

	// Check if getting the report card was successful.
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(models.ReportCardResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorInternalError.Error(),
			},
		})
	}

	// Return the report card.
	return ctx.Status(fiber.StatusOK).JSON(models.ReportCardResponse{
		ReportCard: reportCard,
	})
}
