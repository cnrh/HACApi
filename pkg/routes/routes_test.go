package routes

import (
	"testing"

	"github.com/Threqt1/HACApi/pkg/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// Filter out fields that should not be compared.
var testRoute_Comparer = cmpopts.IgnoreFields(fiber.Route{}, "pos", "use", "star", "root", "path", "routeParser", "group", "Name", "Handlers")

// TestPublicRoutes tests if public routes are registered
// properly.
func TestPublicRoutes(t *testing.T) {
	// Create a testing server.
	server := repository.Server{App: fiber.New(fiber.Config{})}

	// Register the routes.
	PublicRoutes(&server)

	// Confirm routes were registered.
	registered := server.App.GetRoutes(true)

	// Make expected output.
	apiRoute := "/api/v1"
	expected := []fiber.Route{
		// Login.
		{
			Method: "POST",
			Path:   apiRoute + "/login",
			Params: nil,
		},
		// Classwork.
		{
			Method: "POST",
			Path:   apiRoute + "/classwork",
			Params: nil,
		},
		// IPR.
		{
			Method: "POST",
			Path:   apiRoute + "/ipr",
			Params: nil,
		},
		// IPR All.
		{
			Method: "POST",
			Path:   apiRoute + "/ipr/all",
			Params: nil,
		},
		// Report Card.
		{
			Method: "POST",
			Path:   apiRoute + "/reportcard",
			Params: nil,
		},
		// Schedule.
		{
			Method: "POST",
			Path:   apiRoute + "/schedule",
			Params: nil,
		},
		// Transcript.
		{
			Method: "POST",
			Path:   apiRoute + "/transcript",
			Params: nil,
		},
	}

	// Compare them.
	if diff := cmp.Diff(expected, registered, testRoute_Comparer); diff != "" {
		t.Fatalf("Failed for TestPublicRoutes() (-want, +got)\n%s", diff)
	}
}

// TestSwaggerRoute tests if the swagger route is
// registered properly.
func TestSwaggerRoute(t *testing.T) {
	// Create a testing server.
	server := repository.Server{App: fiber.New(fiber.Config{})}

	// Register the routes.
	SwaggerRoute(&server)

	// Confirm routes were registered.
	registered := server.App.GetRoutes(true)

	// Make expected output.
	apiRoute := "/docs"
	expected := []fiber.Route{
		// Swagger.
		{
			Method: "GET",
			Path:   apiRoute + "/*",
			Params: []string{"*1"},
		},
		{
			Method: "HEAD",
			Path:   apiRoute + "/*",
			Params: []string{"*1"},
		},
	}

	// Compare them.
	if diff := cmp.Diff(expected, registered, testRoute_Comparer); diff != "" {
		t.Fatalf("Failed for TestSwaggerRoute() (-want, +got)\n%s", diff)
	}
}
