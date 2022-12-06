package configs

import (
	"github.com/Threqt1/HACApi/app/queries"
	"github.com/Threqt1/HACApi/app/queries/parsers"
	"github.com/Threqt1/HACApi/pkg/repository"
	"github.com/Threqt1/HACApi/pkg/utils"
	"github.com/Threqt1/HACApi/platform/cache"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func ServerConfig() *repository.Server {
	scraper := utils.NewScraper()
	cache := cache.NewCache(scraper)
	parser := parsers.NewParser()
	querier := queries.NewQuerier(scraper, parser)
	app := fiber.New(FiberConfig())
	validator := validator.New()

	return &repository.Server{
		Scraper:   scraper,
		Cache:     cache,
		App:       app,
		Validator: validator,
		Querier:   querier,
		Parser:    parser,
	}
}
