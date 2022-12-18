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

// Test if PostLogin() functions correctly
// given valid input.
func TestPostLogin_AllValidInputs(t *testing.T) {
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

	// Register PostLogin() as the handler
	// for the default route.
	server.App.Post("/", utils.WrapController(server, PostLogin))

	// Create request data.
	bodyData := models.LoginRequestBody{
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
	res := models.LoginResponse{}

	sonic.Unmarshal(resBody, &res)

	// Make expected body.
	expected := utils.ExpectedServerResponse[models.LoginResponse]{
		Status: fiber.StatusOK,
		Body: models.LoginResponse{
			Login: []models.Login{{
				Username: repository.FakeUsername,
				Password: repository.FakePassword,
				Base:     repository.FakeBase,
			}},
		},
	}

	// Convert response to a comparable struct.
	got := utils.ExpectedServerResponse[models.LoginResponse]{
		Status: resp.StatusCode,
		Body:   res,
	}

	// Test.
	if diff := cmp.Diff(expected, got); diff != "" {
		t.Fatalf("Failed for PostLogin() All Valid Inputs (-want, +got)\n%s", diff)
	}
}

// Test if PostLogin() errors out
// for bad body parameters.
func TestPostLogin_BadBodyParams(t *testing.T) {
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

	// Register PostLogin() as the handler
	// for the default route.
	server.App.Post("/", utils.WrapController(server, PostLogin))

	// Create request data.
	bodyData := models.LoginRequestBody{
		BaseRequestBody: models.BaseRequestBody{
			Username: repository.FakeUsername,
			Password: repository.FakePassword,
			Base:     repository.FakeBase,
		},
	}
	body, _ := sonic.Marshal(bodyData)

	// Create a test request. Remove content type header to force an error.
	req := httptest.NewRequest("POST", "http://fake.url/", bytes.NewReader(body))

	// Test the request.
	resp, _ := server.App.Test(req)

	// Parse the body.
	resBody, _ := io.ReadAll(resp.Body)
	res := models.LoginResponse{}

	sonic.Unmarshal(resBody, &res)

	// Make expected body.
	expected := utils.ExpectedServerResponse[models.LoginResponse]{
		Status: fiber.StatusBadRequest,
		Body: models.LoginResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorBadBodyParams.Error(),
			},
		},
	}

	// Convert response to a comparable struct.
	got := utils.ExpectedServerResponse[models.LoginResponse]{
		Status: resp.StatusCode,
		Body:   res,
	}

	// Test.
	if diff := cmp.Diff(expected, got); diff != "" {
		t.Fatalf("Failed for PostLogin() Bad Body Params (-want, +got)\n%s", diff)
	}
}

// Test if PostLogin() errors out
// for invalid request model.
func TestPostLogin_BadBodyParams_InvalidModel(t *testing.T) {
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

	// Register PostLogin() as the handler
	// for the default route.
	server.App.Post("/", utils.WrapController(server, PostLogin))

	// Create request data.
	bodyData := models.LoginRequestBody{
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
	res := models.LoginResponse{}

	sonic.Unmarshal(resBody, &res)

	// Make expected body.
	expected := utils.ExpectedServerResponse[models.LoginResponse]{
		Status: fiber.StatusBadRequest,
		Body: models.LoginResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorBadBodyParams.Error(),
			},
		},
	}

	// Convert response to a comparable struct.
	got := utils.ExpectedServerResponse[models.LoginResponse]{
		Status: resp.StatusCode,
		Body:   res,
	}

	// Test.
	if diff := cmp.Diff(expected, got); diff != "" {
		t.Fatalf("Failed for PostLogin() Bad Body Params, Invalid Request Model (-want, +got)\n%s", diff)
	}
}

// Test if PostLogin() errors out
// for invalid credentials.
func TestPostLogin_InvalidCredentials(t *testing.T) {
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

	// Register PostLogin() as the handler
	// for the default route.
	server.App.Post("/", utils.WrapController(server, PostLogin))

	// Create request data.
	bodyData := models.LoginRequestBody{
		BaseRequestBody: models.BaseRequestBody{
			Username: "bad username",
			Password: "bad password",
			Base:     "bad base",
		},
	}
	body, _ := sonic.Marshal(bodyData)

	// Create a test request. Remove content type header to force an error.
	req := httptest.NewRequest("POST", "http://fake.url/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Test the request.
	resp, _ := server.App.Test(req)

	// Parse the body.
	resBody, _ := io.ReadAll(resp.Body)
	res := models.LoginResponse{}

	sonic.Unmarshal(resBody, &res)

	// Make expected body.
	expected := utils.ExpectedServerResponse[models.LoginResponse]{
		Status: fiber.StatusBadRequest,
		Body: models.LoginResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorInvalidAuthentication.Error(),
			},
		},
	}

	// Convert response to a comparable struct.
	got := utils.ExpectedServerResponse[models.LoginResponse]{
		Status: resp.StatusCode,
		Body:   res,
	}

	// Test.
	if diff := cmp.Diff(expected, got); diff != "" {
		t.Fatalf("Failed for PostLogin() Invalid Credentials (-want, +got)\n%s", diff)
	}
}

// Test if PostLogin() errors out
// for internal errors.
func TestPostLogin_InternalError(t *testing.T) {
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

	// Register PostLogin() as the handler
	// for the default route.
	server.App.Post("/", utils.WrapController(server, PostLogin))

	// Create request data.
	bodyData := models.LoginRequestBody{
		BaseRequestBody: models.BaseRequestBody{
			Username: repository.FakeUsername,
			Password: repository.FakePassword,
			Base:     repository.FakeBase,
		},
	}
	body, _ := sonic.Marshal(bodyData)

	// Create a test request. Remove content type header to force an error.
	req := httptest.NewRequest("POST", "http://fake.url/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Test the request.
	resp, _ := server.App.Test(req)

	// Parse the body.
	resBody, _ := io.ReadAll(resp.Body)
	res := models.LoginResponse{}

	sonic.Unmarshal(resBody, &res)

	// Make expected body.
	expected := utils.ExpectedServerResponse[models.LoginResponse]{
		Status: fiber.StatusInternalServerError,
		Body: models.LoginResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorInternalError.Error(),
			},
		},
	}

	// Convert response to a comparable struct.
	got := utils.ExpectedServerResponse[models.LoginResponse]{
		Status: resp.StatusCode,
		Body:   res,
	}

	// Test.
	if diff := cmp.Diff(expected, got); diff != "" {
		t.Fatalf("Failed for PostLogin() Internal Error (-want, +got)\n%s", diff)
	}
}
