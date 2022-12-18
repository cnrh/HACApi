package controllers

import (
	"fmt"

	"github.com/Threqt1/HACApi/app/models"
	"github.com/Threqt1/HACApi/pkg/repository"
	"github.com/gofiber/fiber/v2"
)

// PostTranscript handles POST request to the transcript endpoint.
//
//	@Description	Returns the transcript for the user.
//	@Tags			transcript
//	@Param			request	body	models.TranscriptRequestBody	false	"Body params"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.TranscriptResponse
//	@Router			/transcript [post]
func PostTranscript(server *repository.Server, ctx *fiber.Ctx) error {
	// Parse body.
	params := new(models.TranscriptRequestBody)

	// Check if the parsing was successful.
	if err := ctx.BodyParser(params); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(models.TranscriptResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorBadBodyParams.Error(),
			},
		})
	}

	// Verify the validity of the body params.
	if err := server.Validator.Struct(params); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(models.TranscriptResponse{
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

	// Check if the login went through.
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(models.TranscriptResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorInvalidAuthentication.Error(),
			},
		})
	}

	// Get the transcript.
	transcript, err := server.Querier.GetTranscript(collector, *params)

	// Check if getting the transcript was successful.
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(models.TranscriptResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorInternalError.Error(),
			},
		})
	}

	// Return the transcript.
	return ctx.Status(fiber.StatusOK).JSON(models.TranscriptResponse{
		Transcript: transcript,
	})
}
