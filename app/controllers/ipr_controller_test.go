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

// Test if PostIPR() functions correctly
// with valid inputs and no date.
func TestPostIPR_AllValidInputs_NoDate(t *testing.T) {
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

	// Register PostIPR() as the handler
	// for the default route.
	server.App.Post("/", utils.WrapController(server, PostIPR))

	// Create request data.
	bodyData := models.IprRequestBody{
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
	res := models.IPRResponse{}

	sonic.Unmarshal(resBody, &res)

	// Make expected body.
	expected := utils.ExpectedServerResponse[models.IPRResponse]{
		Status: fiber.StatusOK,
		Body: models.IPRResponse{
			IPR: []models.IPR{{}},
		},
	}

	// Convert response to a comparable struct.
	got := utils.ExpectedServerResponse[models.IPRResponse]{
		Status: resp.StatusCode,
		Body:   res,
	}

	// Test.
	if diff := cmp.Diff(expected, got); diff != "" {
		t.Fatalf("Failed for PostIPR() All Valid Inputs, No Date (-want, +got)\n%s", diff)
	}
}

// Test if PostIPR() functions correctly
// with valid inputs and a date.
func TestPostIPR_AllValidInputs_WithDate(t *testing.T) {
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

	// Register PostIPR() as the handler
	// for the default route.
	server.App.Post("/", utils.WrapController(server, PostIPR))

	// Create request data.
	bodyData := models.IprRequestBody{
		BaseRequestBody: models.BaseRequestBody{
			Username: repository.FakeUsername,
			Password: repository.FakePassword,
			Base:     repository.FakeBase,
		},
		Date: "01/02/2006",
	}
	body, _ := sonic.Marshal(bodyData)

	// Create a test request.
	req := httptest.NewRequest("POST", "http://fake.url/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Test the request.
	resp, _ := server.App.Test(req)

	// Parse the body.
	resBody, _ := io.ReadAll(resp.Body)
	res := models.IPRResponse{}

	sonic.Unmarshal(resBody, &res)

	// Make expected body.
	expected := utils.ExpectedServerResponse[models.IPRResponse]{
		Status: fiber.StatusOK,
		Body: models.IPRResponse{
			IPR: []models.IPR{{}},
		},
	}

	// Convert response to a comparable struct.
	got := utils.ExpectedServerResponse[models.IPRResponse]{
		Status: resp.StatusCode,
		Body:   res,
	}

	// Test.
	if diff := cmp.Diff(expected, got); diff != "" {
		t.Fatalf("Failed for PostIPR() All Valid Inputs, With Date (-want, +got)\n%s", diff)
	}
}

// Test if PostIPR() errors out
// with bad body parameters.
func TestPostIPR_BadBodyParams(t *testing.T) {
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

	// Register PostIPR() as the handler
	// for the default route.
	server.App.Post("/", utils.WrapController(server, PostIPR))

	// Create request data.
	bodyData := models.IprRequestBody{
		BaseRequestBody: models.BaseRequestBody{
			Username: repository.FakeUsername,
			Password: repository.FakePassword,
			Base:     repository.FakeBase,
		},
	}
	body, _ := sonic.Marshal(bodyData)

	// Create a test request. Leave out content type header to force error.
	req := httptest.NewRequest("POST", "http://fake.url/", bytes.NewReader(body))

	// Test the request.
	resp, _ := server.App.Test(req)

	// Parse the body.
	resBody, _ := io.ReadAll(resp.Body)
	res := models.IPRResponse{}

	sonic.Unmarshal(resBody, &res)

	// Make expected body.
	expected := utils.ExpectedServerResponse[models.IPRResponse]{
		Status: fiber.StatusBadRequest,
		Body: models.IPRResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorBadBodyParams.Error(),
			},
		},
	}

	// Convert response to a comparable struct.
	got := utils.ExpectedServerResponse[models.IPRResponse]{
		Status: resp.StatusCode,
		Body:   res,
	}

	// Test.
	if diff := cmp.Diff(expected, got); diff != "" {
		t.Fatalf("Failed for PostIPR() Bad Body Parameters (-want, +got)\n%s", diff)
	}
}

// Test if PostIPR() errors out
// with an invalid request model.
func TestPostIPR_BadBodyParams_InvalidModel(t *testing.T) {
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

	// Register PostIPR() as the handler
	// for the default route.
	server.App.Post("/", utils.WrapController(server, PostIPR))

	// Create request data.
	bodyData := models.IprRequestBody{
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
	res := models.IPRResponse{}

	sonic.Unmarshal(resBody, &res)

	// Make expected body.
	expected := utils.ExpectedServerResponse[models.IPRResponse]{
		Status: fiber.StatusBadRequest,
		Body: models.IPRResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorBadBodyParams.Error(),
			},
		},
	}

	// Convert response to a comparable struct.
	got := utils.ExpectedServerResponse[models.IPRResponse]{
		Status: resp.StatusCode,
		Body:   res,
	}

	// Test.
	if diff := cmp.Diff(expected, got); diff != "" {
		t.Fatalf("Failed for PostIPR() Bad Body Parameters, Invalid Request Model (-want, +got)\n%s", diff)
	}
}

// Test if PostIPR() errors out
// with an invalid date
func TestPostIPR_InvalidBodyParams_InvalidDate(t *testing.T) {
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

	// Register PostIPR() as the handler
	// for the default route.
	server.App.Post("/", utils.WrapController(server, PostIPR))

	// Create request data.
	bodyData := models.IprRequestBody{
		BaseRequestBody: models.BaseRequestBody{
			Username: repository.FakeUsername,
			Password: repository.FakePassword,
			Base:     repository.FakeBase,
		},
		Date: "1/6/2016",
	}
	body, _ := sonic.Marshal(bodyData)

	// Create a test request.
	req := httptest.NewRequest("POST", "http://fake.url/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Test the request.
	resp, _ := server.App.Test(req)

	// Parse the body.
	resBody, _ := io.ReadAll(resp.Body)
	res := models.IPRResponse{}

	sonic.Unmarshal(resBody, &res)

	// Make expected body.
	expected := utils.ExpectedServerResponse[models.IPRResponse]{
		Status: fiber.StatusBadRequest,
		Body: models.IPRResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorBadBodyParams.Error(),
			},
		},
	}

	// Convert response to a comparable struct.
	got := utils.ExpectedServerResponse[models.IPRResponse]{
		Status: resp.StatusCode,
		Body:   res,
	}

	// Test.
	if diff := cmp.Diff(expected, got); diff != "" {
		t.Fatalf("Failed for PostIPR() Invalid Body Parameters, Bad Date (-want, +got)\n%s", diff)
	}
}

// Test if PostIPR() errors out
// given invalid credentials.
func TestPostIPR_InvalidCredentials(t *testing.T) {
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

	// Register PostIPR() as the handler
	// for the default route.
	server.App.Post("/", utils.WrapController(server, PostIPR))

	// Create request data.
	bodyData := models.IprRequestBody{
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
	res := models.IPRResponse{}

	sonic.Unmarshal(resBody, &res)

	// Make expected body.
	expected := utils.ExpectedServerResponse[models.IPRResponse]{
		Status: fiber.StatusBadRequest,
		Body: models.IPRResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorInvalidAuthentication.Error(),
			},
		},
	}

	// Convert response to a comparable struct.
	got := utils.ExpectedServerResponse[models.IPRResponse]{
		Status: resp.StatusCode,
		Body:   res,
	}

	// Test.
	if diff := cmp.Diff(expected, got); diff != "" {
		t.Fatalf("Failed for PostIPR() Invalid Credentials (-want, +got)\n%s", diff)
	}
}

// Test if PostIPR() errors out
// due to an internal error.
func TestPostIPR_InternalError(t *testing.T) {
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

	// Register PostIPR() as the handler
	// for the default route.
	server.App.Post("/", utils.WrapController(server, PostIPR))

	// Create request data.
	bodyData := models.IprRequestBody{
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
	res := models.IPRResponse{}

	sonic.Unmarshal(resBody, &res)

	// Make expected body.
	expected := utils.ExpectedServerResponse[models.IPRResponse]{
		Status: fiber.StatusInternalServerError,
		Body: models.IPRResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorInternalError.Error(),
			},
		},
	}

	// Convert response to a comparable struct.
	got := utils.ExpectedServerResponse[models.IPRResponse]{
		Status: resp.StatusCode,
		Body:   res,
	}

	// Test.
	if diff := cmp.Diff(expected, got); diff != "" {
		t.Fatalf("Failed for PostIPR() Internal Error (-want, +got)\n%s", diff)
	}
}
