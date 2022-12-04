package controllers

import (
	"fmt"

	"github.com/Threqt1/HACApi/app/queries"
	"github.com/Threqt1/HACApi/pkg/repository"
	"github.com/Threqt1/HACApi/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type transcriptRequestBody struct {
	utils.BaseRequestBody
}

// PostTranscript handles POST request to the transcript endpoint.
// @Description Returns the transcript for the user.
// @Tags        transcript
// @Param       request body transcriptRequestBody false "Body params"
// @Accept      json
// @Produce     json
// @Success     200 {object} models.TranscriptResponse
// @Router      /transcript [post]
func PostTranscript(server *repository.Server, ctx *fiber.Ctx) error {
	// Parse body
	params := new(transcriptRequestBody)

	if err := ctx.BodyParser(params); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":        true,
			"msg":        "Bad body params",
			"transcript": nil,
		})
	}

	// Verify validity of body params
	if err := server.Validator.Struct(params); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":        true,
			"msg":        "Bad body params",
			"transcript": nil,
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
			"transcript": nil,
		})
	}

	// Get transcript
	transcript, err := queries.GetTranscript(server, collector, params.Base)

	// Check if returned value was nil
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":        true,
			"msg":        "Transcript not found. Might be an internal error",
			"transcript": nil,
		})
	}

	// All is well
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"err":        false,
		"msg":        "",
		"transcript": transcript,
	})
}
