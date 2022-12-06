package queries

import (
	"github.com/Threqt1/HACApi/app/models"
	"github.com/Threqt1/HACApi/pkg/repository"
	"github.com/gocolly/colly"
)

// getLogin returns the logged-in credentials for the user inputted.
func getLogin(scraper repository.ScraperProvider, parser repository.ParserProvider, collector *colly.Collector, params models.LoginRequestBody) ([]models.Login, error) {
	// Form the response
	loginRes := models.Login{
		Username: params.Username,
		Password: params.Password,
		Base:     params.Base,
	}

	return []models.Login{loginRes}, nil
}
