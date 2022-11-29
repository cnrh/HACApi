package controllers

import (
	"fmt"

	"github.com/Threqt1/HACApi/app/queries"
	"github.com/Threqt1/HACApi/pkg/utils"
	"github.com/Threqt1/HACApi/platform/cache"
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
func PostTranscript(ctx *fiber.Ctx) error {
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
	bodyParamsValid := true

	// Confirm no required body parameters are empty
	if params.Username == "" || params.Password == "" || params.Base == "" {
		bodyParamsValid = false
	}

	if !bodyParamsValid {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":        true,
			"msg":        "Bad body params",
			"transcript": nil,
		})
	}

	// Form cache key
	cacheKey := fmt.Sprintf("%s\n%s\n%s", params.Username, params.Password, params.Base)

	// Try logging in, or grab cached collector
	collector := cache.CurrentCache().Get(cacheKey)

	// Error out if login fails
	if collector == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":        true,
			"msg":        "Invalid username/password/base",
			"transcript": nil,
		})
	}

	// Get transcript
	transcript, err := queries.GetTranscript(collector.Value(), params.Base)

	// Check if returned value was nil
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
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
