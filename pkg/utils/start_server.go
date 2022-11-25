package utils

import (
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
)

func StartForProd(app *fiber.App) {
	//Create channel to confirm when connections are closed
	connsClosedChan := make(chan struct{})

	go func() {
		//Catch os signals
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		//Gracefully shutdown
		if err := app.Shutdown(); err != nil {
			log.Fatalf("Server failed to shutdown. Reason: %v", err)
		}

		close(connsClosedChan)
	}()

	// Build fiber URL
	fiberConnURL, _ := BuildConnectionURL("fiber")

	// Start server
	if err := app.Listen(fiberConnURL); err != nil {
		log.Fatalf("Server failed to load. Error: %v", err)
	}

	//Wait till conns are closed before stopping
	<-connsClosedChan
}

func StartForDev(app *fiber.App) {
	// Build fiber URL
	fiberConnURL, _ := BuildConnectionURL("fiber")

	// Start server
	if err := app.Listen(fiberConnURL); err != nil {
		log.Fatalf("Server failed to load. Error: %v", err)
	}
}
