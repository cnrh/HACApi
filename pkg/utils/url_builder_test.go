package utils

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// testUrl_EnvVar represents an environment variable.
type testUrl_EnvVar struct {
	Key   string
	Value string
}

// testUrl_Test represents an expected output from a test.
type testUrl_Test struct {
	Value string // The returned URL.
	Error error  // The expected error, if any.
}

// testUrl_Case represents an individual test case for BuildURLConnection().
type testUrl_Case struct {
	Test  testUrl_Test
	Input string           // The input for the URL building function.
	Env   []testUrl_EnvVar // The environment to simulate.
}

// TestBuildConnectionURL tests BuildURLConnection() for different inputs.
func TestBuildConnectionURL(t *testing.T) {
	// Set up all test cases.
	cases := []testUrl_Case{
		// Check for host URL.
		{Input: "fiber", Test: testUrl_Test{Value: "127.0.0.1:3000", Error: nil}, Env: []testUrl_EnvVar{{Key: "SERVER_HOST", Value: "127.0.0.1"}, {Key: "SERVER_PORT", Value: "3000"}}},
		// Check for error on empty/invalid input.
		{Input: "", Test: testUrl_Test{Value: "", Error: ErrorURLConnection}},
	}

	for _, test := range cases {
		// Set environment variables, store previous values for after tests.
		prevEnv := make([]testUrl_EnvVar, len(test.Env))
		for i, envVar := range test.Env {
			prevEnv[i] = testUrl_EnvVar{Key: envVar.Key, Value: os.Getenv(envVar.Key)}
			os.Setenv(envVar.Key, envVar.Value)
		}

		// Test.
		url, err := BuildConnectionURL(test.Input)

		if diff := cmp.Diff(test.Test, testUrl_Test{Value: url, Error: err}, cmpopts.EquateErrors()); diff != "" {
			t.Fatalf("Failed for BuildConnectionURL() (-want, +got):\n%s", diff)
		}

		// Reset environment variables.
		for _, prevVar := range prevEnv {
			os.Setenv(prevVar.Key, prevVar.Value)
		}
	}
}
