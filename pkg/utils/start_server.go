package utils

import (
	"log"
	"os"
	"os/signal"

	"github.com/Threqt1/HACApi/pkg/repository"
)

func StartForProd(server *repository.Server) {
	// Create channel to confirm when connections are closed
	connsClosedChan := make(chan struct{})

	go func() {
		// Catch os signals
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// Gracefully shutdown
		if err := server.App.Shutdown(); err != nil {
			log.Fatalf("Server failed to shutdown. Reason: %v", err)
		}

		close(connsClosedChan)
	}()

	// Build fiber URL
	fiberConnURL, _ := BuildConnectionURL("fiber")

	// Start server
	if err := server.App.Listen(fiberConnURL); err != nil {
		log.Fatalf("Server failed to load. Error: %v", err)
	}

	// Wait till conns are closed before stopping
	<-connsClosedChan
}

func StartForDev(server *repository.Server) {
	// Build fiber URL
	fiberConnURL, _ := BuildConnectionURL("fiber")

	// Start server
	if err := server.App.Listen(fiberConnURL); err != nil {
		log.Fatalf("Server failed to load. Error: %v", err)
	}
}
