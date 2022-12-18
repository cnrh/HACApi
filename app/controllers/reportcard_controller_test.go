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

// Test if PostReportCard() functions correctly
// with valid inputs and no date.
func TestPostReportCard(t *testing.T) {
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

	// Register PostReportCard() as the handler
	// for the default route.
	server.App.Post("/", utils.WrapController(server, PostReportCard))

	// Create request data.
	bodyData := models.ReportCardRequestBody{
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
	res := models.ReportCardResponse{}

	sonic.Unmarshal(resBody, &res)

	// Make expected body.
	expected := utils.ExpectedServerResponse[models.ReportCardResponse]{
		Status: fiber.StatusOK,
		Body: models.ReportCardResponse{
			ReportCard: []models.ReportCard{{}},
		},
	}

	// Convert response to a comparable struct.
	got := utils.ExpectedServerResponse[models.ReportCardResponse]{
		Status: resp.StatusCode,
		Body:   res,
	}

	// Test.
	if diff := cmp.Diff(expected, got); diff != "" {
		t.Fatalf("Failed for PostReportCard() All Valid Inputs (-want, +got)\n%s", diff)
	}
}

// Test if PostReportCard() errors out
// if the body params are bad.
func TestPostReportCard_BadBodyParams(t *testing.T) {
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

	// Register PostReportCard() as the handler
	// for the default route.
	server.App.Post("/", utils.WrapController(server, PostReportCard))

	// Create request data.
	bodyData := models.ReportCardRequestBody{
		BaseRequestBody: models.BaseRequestBody{
			Username: repository.FakeUsername,
			Password: repository.FakePassword,
			Base:     repository.FakeBase,
		},
	}
	body, _ := sonic.Marshal(bodyData)

	// Create a test request. Don't include header to force an error.
	req := httptest.NewRequest("POST", "http://fake.url/", bytes.NewReader(body))

	// Test the request.
	resp, _ := server.App.Test(req)

	// Parse the body.
	resBody, _ := io.ReadAll(resp.Body)
	res := models.ReportCardResponse{}

	sonic.Unmarshal(resBody, &res)

	// Make expected body.
	expected := utils.ExpectedServerResponse[models.ReportCardResponse]{
		Status: fiber.StatusBadRequest,
		Body: models.ReportCardResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorBadBodyParams.Error(),
			},
		},
	}

	// Convert response to a comparable struct.
	got := utils.ExpectedServerResponse[models.ReportCardResponse]{
		Status: resp.StatusCode,
		Body:   res,
	}

	// Test.
	if diff := cmp.Diff(expected, got); diff != "" {
		t.Fatalf("Failed for PostReportCard() Bad Body Params (-want, +got)\n%s", diff)
	}
}

// Test if PostReportCard() errors out
// if the request model is wrong.
func TestPostReportCard_BadBodyParams_InvalidModel(t *testing.T) {
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

	// Register PostReportCard() as the handler
	// for the default route.
	server.App.Post("/", utils.WrapController(server, PostReportCard))

	// Create request data.
	bodyData := models.ReportCardRequestBody{
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
	res := models.ReportCardResponse{}

	sonic.Unmarshal(resBody, &res)

	// Make expected body.
	expected := utils.ExpectedServerResponse[models.ReportCardResponse]{
		Status: fiber.StatusBadRequest,
		Body: models.ReportCardResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorBadBodyParams.Error(),
			},
		},
	}

	// Convert response to a comparable struct.
	got := utils.ExpectedServerResponse[models.ReportCardResponse]{
		Status: resp.StatusCode,
		Body:   res,
	}

	// Test.
	if diff := cmp.Diff(expected, got); diff != "" {
		t.Fatalf("Failed for PostReportCard() Bad Body Params, Invalid Request Body (-want, +got)\n%s", diff)
	}
}

// Test if PostReportCard() errors out
// if the credentials are wrong
func TestPostReportCard_InvalidCredentials(t *testing.T) {
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

	// Register PostReportCard() as the handler
	// for the default route.
	server.App.Post("/", utils.WrapController(server, PostReportCard))

	// Create request data.
	bodyData := models.ReportCardRequestBody{
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
	res := models.ReportCardResponse{}

	sonic.Unmarshal(resBody, &res)

	// Make expected body.
	expected := utils.ExpectedServerResponse[models.ReportCardResponse]{
		Status: fiber.StatusBadRequest,
		Body: models.ReportCardResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorInvalidAuthentication.Error(),
			},
		},
	}

	// Convert response to a comparable struct.
	got := utils.ExpectedServerResponse[models.ReportCardResponse]{
		Status: resp.StatusCode,
		Body:   res,
	}

	// Test.
	if diff := cmp.Diff(expected, got); diff != "" {
		t.Fatalf("Failed for PostReportCard() Invalid Credentials (-want, +got)\n%s", diff)
	}
}

// Test if PostReportCard() errors out
// if the credentials are wrong
func TestPostReportCard_InternalError(t *testing.T) {
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

	// Register PostReportCard() as the handler
	// for the default route.
	server.App.Post("/", utils.WrapController(server, PostReportCard))

	// Create request data.
	bodyData := models.ReportCardRequestBody{
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
	res := models.ReportCardResponse{}

	sonic.Unmarshal(resBody, &res)

	// Make expected body.
	expected := utils.ExpectedServerResponse[models.ReportCardResponse]{
		Status: fiber.StatusInternalServerError,
		Body: models.ReportCardResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorInternalError.Error(),
			},
		},
	}

	// Convert response to a comparable struct.
	got := utils.ExpectedServerResponse[models.ReportCardResponse]{
		Status: resp.StatusCode,
		Body:   res,
	}

	// Test.
	if diff := cmp.Diff(expected, got); diff != "" {
		t.Fatalf("Failed for PostReportCard() Internal Error (-want, +got)\n%s", diff)
	}
}
