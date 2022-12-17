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

// Test if PostClasswork() functions correctly
// with valid inputs with no marking periods
// specified.
func TestPostClasswork_AllValidInputs_NoMarkingPeriods(t *testing.T) {
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

	// Register PostClasswork() as the handler
	// for the default route.
	server.App.Post("/", utils.WrapController(server, PostClasswork))

	// Create request data.
	bodyData := models.ClassworkRequestBody{
		BaseRequestBody: models.BaseRequestBody{
			Username: repository.FakeUsername,
			Password: repository.FakePassword,
			Base:     repository.FakeBase,
		},
		MarkingPeriods: nil,
	}
	body, _ := sonic.Marshal(bodyData)

	// Create a test request.
	req := httptest.NewRequest("POST", "http://fake.url/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Test the request.
	resp, _ := server.App.Test(req)

	// Parse the body.
	resBody, _ := io.ReadAll(resp.Body)
	res := models.ClassworkResponse{}

	sonic.Unmarshal(resBody, &res)

	// Make expected body.
	expected := utils.ExpectedServerResponse[models.ClassworkResponse]{
		Status: fiber.StatusOK,
		Body: models.ClassworkResponse{
			Classwork: []models.Classwork{{}},
		},
	}

	// Convert response to a comparable struct.
	got := utils.ExpectedServerResponse[models.ClassworkResponse]{
		Status: resp.StatusCode,
		Body:   res,
	}

	// Test.
	if diff := cmp.Diff(expected, got); diff != "" {
		t.Fatalf("Failed for PostClasswork() All Valid Inputs, No Marking Periods (-want, +got)\n%s", diff)
	}
}

// Test if PostClasswork() functions correctly
// with valid inputs with 6 marking periods
// specified.
func TestPostClasswork_AllValidInputs_WithMarkingPeriods(t *testing.T) {
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

	// Register PostClasswork() as the handler
	// for the default route.
	server.App.Post("/", utils.WrapController(server, PostClasswork))

	// Create request data.
	bodyData := models.ClassworkRequestBody{
		BaseRequestBody: models.BaseRequestBody{
			Username: repository.FakeUsername,
			Password: repository.FakePassword,
			Base:     repository.FakeBase,
		},
		MarkingPeriods: []int{1, 2, 3, 4, 5, 6},
	}
	body, _ := sonic.Marshal(bodyData)

	// Create a test request.
	req := httptest.NewRequest("POST", "http://fake.url/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Test the request.
	resp, _ := server.App.Test(req)

	// Parse the body.
	resBody, _ := io.ReadAll(resp.Body)
	res := models.ClassworkResponse{}

	sonic.Unmarshal(resBody, &res)

	// Make expected body.
	expected := utils.ExpectedServerResponse[models.ClassworkResponse]{
		Status: fiber.StatusOK,
		Body: models.ClassworkResponse{
			Classwork: []models.Classwork{{}, {}, {}, {}, {}, {}},
		},
	}

	// Convert response to a comparable struct.
	got := utils.ExpectedServerResponse[models.ClassworkResponse]{
		Status: resp.StatusCode,
		Body:   res,
	}

	// Test.
	if diff := cmp.Diff(expected, got); diff != "" {
		t.Fatalf("Failed for PostClasswork() All Valid Inputs, With Marking Periods (-want, +got)\n%s", diff)
	}
}

// Test if PostClasswork() errors out due to bad
// body parameters.
func TestPostClasswork_BadBodyParams(t *testing.T) {
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

	// Register PostClasswork() as the handler
	// for the default route.
	server.App.Post("/", utils.WrapController(server, PostClasswork))

	// Create request data.
	bodyData := models.ClassworkRequestBody{
		BaseRequestBody: models.BaseRequestBody{
			Username: repository.FakeUsername,
			Password: repository.FakePassword,
			Base:     repository.FakeBase,
		},
		MarkingPeriods: nil,
	}
	body, _ := sonic.Marshal(bodyData)

	// Create a test request. Purposefully leave out content type to force error.
	req := httptest.NewRequest("POST", "http://fake.url/", bytes.NewReader(body))

	// Test the request.
	resp, _ := server.App.Test(req)

	// Parse the body.
	resBody, _ := io.ReadAll(resp.Body)
	res := models.ClassworkResponse{}

	sonic.Unmarshal(resBody, &res)

	// Make expected body.
	expected := utils.ExpectedServerResponse[models.ClassworkResponse]{
		Status: fiber.StatusBadRequest,
		Body: models.ClassworkResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorBadBodyParams.Error(),
			},
			Classwork: nil,
		},
	}

	// Convert response to a comparable struct.
	got := utils.ExpectedServerResponse[models.ClassworkResponse]{
		Status: resp.StatusCode,
		Body:   res,
	}

	// Test.
	if diff := cmp.Diff(expected, got); diff != "" {
		t.Fatalf("Failed for PostClasswork() Bad Body Parameters (-want, +got)\n%s", diff)
	}
}

// Test if PostClasswork() errors out due to invalid
// body parameters, specifically adding too many marking
// periods.
func TestPostClasswork_InvalidBodyParams_GreaterThan6MarkingPeriods(t *testing.T) {
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

	// Register PostClasswork() as the handler
	// for the default route.
	server.App.Post("/", utils.WrapController(server, PostClasswork))

	// Create request data.
	bodyData := models.ClassworkRequestBody{
		BaseRequestBody: models.BaseRequestBody{
			Username: repository.FakeUsername,
			Password: repository.FakePassword,
			Base:     repository.FakeBase,
		},
		MarkingPeriods: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	}
	body, _ := sonic.Marshal(bodyData)

	// Create a test request. Purposefully leave out content type to force error.
	req := httptest.NewRequest("POST", "http://fake.url/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Test the request.
	resp, _ := server.App.Test(req)

	// Parse the body.
	resBody, _ := io.ReadAll(resp.Body)
	res := models.ClassworkResponse{}

	sonic.Unmarshal(resBody, &res)

	// Make expected body.
	expected := utils.ExpectedServerResponse[models.ClassworkResponse]{
		Status: fiber.StatusBadRequest,
		Body: models.ClassworkResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorBadBodyParams.Error(),
			},
			Classwork: nil,
		},
	}

	// Convert response to a comparable struct.
	got := utils.ExpectedServerResponse[models.ClassworkResponse]{
		Status: resp.StatusCode,
		Body:   res,
	}

	// Test.
	if diff := cmp.Diff(expected, got); diff != "" {
		t.Fatalf("Failed for PostClasswork() Invalid Body Parameters, >6 Marking Periods (-want, +got)\n%s", diff)
	}
}

// Test if PostClasswork() errors out due to invalid
// body parameters, specifically having invalid marking
// periods.
func TestPostClasswork_InvalidBodyParams_MarkingPeriodOutOfRange(t *testing.T) {
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

	// Register PostClasswork() as the handler
	// for the default route.
	server.App.Post("/", utils.WrapController(server, PostClasswork))

	// Create request data.
	bodyData := models.ClassworkRequestBody{
		BaseRequestBody: models.BaseRequestBody{
			Username: repository.FakeUsername,
			Password: repository.FakePassword,
			Base:     repository.FakeBase,
		},
		MarkingPeriods: []int{0, 7, 10, -1, -2},
	}
	body, _ := sonic.Marshal(bodyData)

	// Create a test request. Purposefully leave out content type to force error.
	req := httptest.NewRequest("POST", "http://fake.url/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Test the request.
	resp, _ := server.App.Test(req)

	// Parse the body.
	resBody, _ := io.ReadAll(resp.Body)
	res := models.ClassworkResponse{}

	sonic.Unmarshal(resBody, &res)

	// Make expected body.
	expected := utils.ExpectedServerResponse[models.ClassworkResponse]{
		Status: fiber.StatusBadRequest,
		Body: models.ClassworkResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorBadBodyParams.Error(),
			},
			Classwork: nil,
		},
	}

	// Convert response to a comparable struct.
	got := utils.ExpectedServerResponse[models.ClassworkResponse]{
		Status: resp.StatusCode,
		Body:   res,
	}

	// Test.
	if diff := cmp.Diff(expected, got); diff != "" {
		t.Fatalf("Failed for PostClasswork() Invalid Body Parameters, Out Of Range Marking Periods (-want, +got)\n%s", diff)
	}
}

// Test if PostClasswork() errors out due to invalid
// credentials.
func TestPostClasswork_InvalidCredentials(t *testing.T) {
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

	// Register PostClasswork() as the handler
	// for the default route.
	server.App.Post("/", utils.WrapController(server, PostClasswork))

	// Create request data.
	bodyData := models.ClassworkRequestBody{
		BaseRequestBody: models.BaseRequestBody{
			Username: "bad username",
			Password: "bad password",
			Base:     "bad base",
		},
		MarkingPeriods: nil,
	}
	body, _ := sonic.Marshal(bodyData)

	// Create a test request. Purposefully leave out content type to force error.
	req := httptest.NewRequest("POST", "http://fake.url/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Test the request.
	resp, _ := server.App.Test(req)

	// Parse the body.
	resBody, _ := io.ReadAll(resp.Body)
	res := models.ClassworkResponse{}

	sonic.Unmarshal(resBody, &res)

	// Make expected body.
	expected := utils.ExpectedServerResponse[models.ClassworkResponse]{
		Status: fiber.StatusBadRequest,
		Body: models.ClassworkResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorInvalidAuthentication.Error(),
			},
			Classwork: nil,
		},
	}

	// Convert response to a comparable struct.
	got := utils.ExpectedServerResponse[models.ClassworkResponse]{
		Status: resp.StatusCode,
		Body:   res,
	}

	// Test.
	if diff := cmp.Diff(expected, got); diff != "" {
		t.Fatalf("Failed for PostClasswork() Invalid Authentication (-want, +got)\n%s", diff)
	}
}

// Test if PostClasswork() errors out due to an
// internal error.
func TestPostClasswork_InternalError(t *testing.T) {
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

	// Register PostClasswork() as the handler
	// for the default route.
	server.App.Post("/", utils.WrapController(server, PostClasswork))

	// Create request data.
	bodyData := models.ClassworkRequestBody{
		BaseRequestBody: models.BaseRequestBody{
			Username: repository.FakeUsername,
			Password: repository.FakePassword,
			Base:     repository.FakeBase,
		},
		MarkingPeriods: nil,
	}
	body, _ := sonic.Marshal(bodyData)

	// Create a test request. Purposefully leave out content type to force error.
	req := httptest.NewRequest("POST", "http://fake.url/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Test the request.
	resp, _ := server.App.Test(req)

	// Parse the body.
	resBody, _ := io.ReadAll(resp.Body)
	res := models.ClassworkResponse{}

	sonic.Unmarshal(resBody, &res)

	// Make expected body.
	expected := utils.ExpectedServerResponse[models.ClassworkResponse]{
		Status: fiber.StatusInternalServerError,
		Body: models.ClassworkResponse{
			HTTPError: models.HTTPError{
				Error:   true,
				Message: repository.ErrorInternalError.Error(),
			},
			Classwork: nil,
		},
	}

	// Convert response to a comparable struct.
	got := utils.ExpectedServerResponse[models.ClassworkResponse]{
		Status: resp.StatusCode,
		Body:   res,
	}

	// Test.
	if diff := cmp.Diff(expected, got); diff != "" {
		t.Fatalf("Failed for PostClasswork() Internal Error (-want, +got)\n%s", diff)
	}
}
