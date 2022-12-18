package controllers

import (
	"bytes"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/Threqt1/HACApi/app/models"
	"github.com/Threqt1/HACApi/app/queries"
	"github.com/Threqt1/HACApi/pkg/repository"
	"github.com/Threqt1/HACApi/pkg/utils"
	"github.com/Threqt1/HACApi/platform/cache"
	"github.com/bytedance/sonic"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/go-cmp/cmp"
)

// Test if PostTranscript() functions correctly
// with valid inputs and no date.
func TestPostTranscript(t *testing.T) {
	// Set up testing server.
	server := &repository.Server{
		App: fiber.New(fiber.Config{
			JSONEncoder: sonic.Marshal,
			JSONDecoder: sonic.Unmarshal,
		}),
		Querier:   queries.NewTestQuerier(),
		Validator: validator.New(),
		Cache:     cache.NewTestCache(),
	}

	// Register PostTranscript() as the handler
	// for the default route.
	server.App.Post("/", utils.WrapController(server, PostTranscript))

	// Create request data.
	bodyData := models.TranscriptRequestBody{
		BaseRequestBody: models.BaseRequestBody{
			Username: repository.FakeUsername,
			Password: repository.FakePassword,
			Base:     repository.FakeBase,
		},
	}
	body, _ := sonic.Marshal(bodyData)

	// Create a test request.
	req := httptest.NewRequest("POST", "http://fake.url/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Test the request.
	resp, _ := server.App.Test(req)

	// Parse the body.
	resBody, _ := io.ReadAll(resp.Body)
	res := models.TranscriptResponse{}

	sonic.Unmarshal(resBody, &res)

	// Make expected body.
	expected := utils.ExpectedServerResponse[models.TranscriptResponse]{
		Status: fiber.StatusOK,
		Body: models.TranscriptResponse{
			Transcript: []models.Transcript{{}},
		},
	}

	// Convert response to a comparable struct.
	got := utils.ExpectedServerResponse[models.TranscriptResponse]{
		Status: resp.StatusCode,
		Body:   res,
	}

	// Test.
	if diff := cmp.Diff(expected, got); diff != "" {
		t.Fatalf("Failed for PostTranscript() All Valid Inputs (-want, +got)\n%s", diff)
	}
}

// Test if PostTranscript() errors out
// if the body parameters are invalid.
func TestPostTranscript_BadBodyParams(t *testing.T) {
	// Set up testing server.
	server := &repository.Server{
		App: fiber.New(fiber.Config{
			JSONEncoder: sonic.Marshal,
			JSONDecoder: sonic.Unmarshal,
		}),
		Querier:   queries.NewTestQuerier(),
		Validator: validator.New(),
		Cache:     cache.NewTestCache(),
	}

	// Register PostTranscript() as the handler
	// for the default route.
	server.App.Post("/", utils.WrapController(server, PostTranscript))

	// Create request data.
	bodyData := models.TranscriptRequestBody{
		BaseRequestBody: models.BaseRequestBody{
			Username: repository.FakeUsername,
			Password: repository.FakePassword,
			Base:     repository.FakeBase,
		},
	}
	body, _ := sonic.Marshal(bodyData)

	// Create a test request. Remove the header on purpose to force an error.
	req := httptest.NewRequest("POST", "http://fake.url/", bytes.NewReader(body))

	// Test the request.
	resp, _ := server.App.Test(req)

	// Parse the body.
	resBody, _ := io.ReadAll(resp.Body)
	res := models.TranscriptResponse{}

	sonic.Unmarshal(resBody, &res)

	// Make expected body.
	expected := utils.ExpectedServerResponse[models.TranscriptResponse]{
		Status: fiber.StatusBadRequest,
		Body: models.TranscriptResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorBadBodyParams.Error(),
			},
		},
	}

	// Convert response to a comparable struct.
	got := utils.ExpectedServerResponse[models.TranscriptResponse]{
		Status: resp.StatusCode,
		Body:   res,
	}

	// Test.
	if diff := cmp.Diff(expected, got); diff != "" {
		t.Fatalf("Failed for PostTranscript() Bad Body Params (-want, +got)\n%s", diff)
	}
}

// Test if PostTranscript() errors out
// if the request model is invalid.
func TestPostTranscript_BadBodyParams_InvalidModel(t *testing.T) {
	// Set up testing server.
	server := &repository.Server{
		App: fiber.New(fiber.Config{
			JSONEncoder: sonic.Marshal,
			JSONDecoder: sonic.Unmarshal,
		}),
		Querier:   queries.NewTestQuerier(),
		Validator: validator.New(),
		Cache:     cache.NewTestCache(),
	}

	// Register PostTranscript() as the handler
	// for the default route.
	server.App.Post("/", utils.WrapController(server, PostTranscript))

	// Create request data.
	bodyData := models.TranscriptRequestBody{
		BaseRequestBody: models.BaseRequestBody{
			Username: "",
			Password: "",
			Base:     "",
		},
	}
	body, _ := sonic.Marshal(bodyData)

	// Create a test request.
	req := httptest.NewRequest("POST", "http://fake.url/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Test the request.
	resp, _ := server.App.Test(req)

	// Parse the body.
	resBody, _ := io.ReadAll(resp.Body)
	res := models.TranscriptResponse{}

	sonic.Unmarshal(resBody, &res)

	// Make expected body.
	expected := utils.ExpectedServerResponse[models.TranscriptResponse]{
		Status: fiber.StatusBadRequest,
		Body: models.TranscriptResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorBadBodyParams.Error(),
			},
		},
	}

	// Convert response to a comparable struct.
	got := utils.ExpectedServerResponse[models.TranscriptResponse]{
		Status: resp.StatusCode,
		Body:   res,
	}

	// Test.
	if diff := cmp.Diff(expected, got); diff != "" {
		t.Fatalf("Failed for PostTranscript() Bad Body Params, Invalid Request Model (-want, +got)\n%s", diff)
	}
}

// Test if PostTranscript() errors out
// if the credentials are invalid.
func TestPostTranscript_InvalidCredentials(t *testing.T) {
	// Set up testing server.
	server := &repository.Server{
		App: fiber.New(fiber.Config{
			JSONEncoder: sonic.Marshal,
			JSONDecoder: sonic.Unmarshal,
		}),
		Querier:   queries.NewTestQuerier(),
		Validator: validator.New(),
		Cache:     cache.NewTestCache(),
	}

	// Register PostTranscript() as the handler
	// for the default route.
	server.App.Post("/", utils.WrapController(server, PostTranscript))

	// Create request data.
	bodyData := models.TranscriptRequestBody{
		BaseRequestBody: models.BaseRequestBody{
			Username: "bad username",
			Password: "bad password",
			Base:     "bad base",
		},
	}
	body, _ := sonic.Marshal(bodyData)

	// Create a test request.
	req := httptest.NewRequest("POST", "http://fake.url/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Test the request.
	resp, _ := server.App.Test(req)

	// Parse the body.
	resBody, _ := io.ReadAll(resp.Body)
	res := models.TranscriptResponse{}

	sonic.Unmarshal(resBody, &res)

	// Make expected body.
	expected := utils.ExpectedServerResponse[models.TranscriptResponse]{
		Status: fiber.StatusBadRequest,
		Body: models.TranscriptResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorInvalidAuthentication.Error(),
			},
		},
	}

	// Convert response to a comparable struct.
	got := utils.ExpectedServerResponse[models.TranscriptResponse]{
		Status: resp.StatusCode,
		Body:   res,
	}

	// Test.
	if diff := cmp.Diff(expected, got); diff != "" {
		t.Fatalf("Failed for PostTranscript() Invalid Credentials (-want, +got)\n%s", diff)
	}
}

// Test if PostTranscript() errors out
// if tjere is an internal error.
func TestPostTranscript_InternalError(t *testing.T) {
	// Set up testing server.
	server := &repository.Server{
		App: fiber.New(fiber.Config{
			JSONEncoder: sonic.Marshal,
			JSONDecoder: sonic.Unmarshal,
		}),
		Querier:   queries.NewTestErrorQuerier(),
		Validator: validator.New(),
		Cache:     cache.NewTestCache(),
	}

	// Register PostTranscript() as the handler
	// for the default route.
	server.App.Post("/", utils.WrapController(server, PostTranscript))

	// Create request data.
	bodyData := models.TranscriptRequestBody{
		BaseRequestBody: models.BaseRequestBody{
			Username: repository.FakeUsername,
			Password: repository.FakePassword,
			Base:     repository.FakeBase,
		},
	}
	body, _ := sonic.Marshal(bodyData)

	// Create a test request.
	req := httptest.NewRequest("POST", "http://fake.url/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Test the request.
	resp, _ := server.App.Test(req)

	// Parse the body.
	resBody, _ := io.ReadAll(resp.Body)
	res := models.TranscriptResponse{}

	sonic.Unmarshal(resBody, &res)

	// Make expected body.
	expected := utils.ExpectedServerResponse[models.TranscriptResponse]{
		Status: fiber.StatusInternalServerError,
		Body: models.TranscriptResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorInternalError.Error(),
			},
		},
	}

	// Convert response to a comparable struct.
	got := utils.ExpectedServerResponse[models.TranscriptResponse]{
		Status: resp.StatusCode,
		Body:   res,
	}

	// Test.
	if diff := cmp.Diff(expected, got); diff != "" {
		t.Fatalf("Failed for PostTranscript() Internal Error (-want, +got)\n%s", diff)
	}
}
