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
	scraperService := utils.NewScraper()
	cacheService := cache.NewCache(scraperService)
	parserService := parsers.NewParser()
	queryService := queries.NewQuerier(scraperService, parserService)
	appService := fiber.New(FiberConfig())
	validatorService := validator.New()

	return &repository.Server{
		Scraper:   scraperService,
		Cache:     cacheService,
		App:       appService,
		Validator: validatorService,
		Querier:   queryService,
		Parser:    parserService,
	}
}
