package controllers

import (
	"github.com/gofiber/fiber/v2"
)

type ClassworkQueryParams struct {
	// The username to login with
	Username string `json:"username" validate:"required" example:"j152123"`
	// The password to login with
	Password string `json:"password" validate:"required" example:"ltia2392"`
	// The base URL for the Homeaccess Center
	Base string `json:"base" validate:"required" example:"homeaccess.katyisd.org"`
	// The six weeks to pull from
	SixWeeks []int `json:"six_weeks" swaggerignore:"true"`
}

// GetClasswork func gets classwork for a specified user for a specific six weeks, if specified.
// @Description Gets classwork for a user
// @Tags Classwork
// @Param request query ClassworkQueryParams true "The query params"
// @Param six_weeks query []int false "The six weeks to pull classwork from" Enums(1, 2, 3, 4, 5, 6, 7) example("1,2,3")
// @Accept json
// @Produce json
// @Success 200 {object} models.Classwork
// @Failure 400 {object} utils.HTTPError
// @Failure 404 {object} utils.HTTPError
// @Router /v1/classwork [get]
func GetClasswork(context *fiber.Ctx) error {
	return nil
}
