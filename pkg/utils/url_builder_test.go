package utils

import (
	"errors"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// urlEnvVar represents an environment variable.
type urlEnvVar struct {
	Key   string
	Value string
}

// urlTest represents an expected output from a test.
type urlTest struct {
	Value string // The returned URL.
	Error error  // The expected error, if any.
}

// urlCase represents an individual test case for BuildURLConnection().
type urlCase struct {
	Test  urlTest
	Input string      // The input for the URL building function.
	Env   []urlEnvVar // The environment to simulate.
}

// TestBuildConnectionURL tests BuildURLConnection() for different inputs.
func TestBuildConnectionURL(t *testing.T) {
	// Set up all test cases.
	cases := []urlCase{
		// Check for host URL.
		{Input: "fiber", Test: urlTest{Value: "127.0.0.1:3000", Error: nil}, Env: []urlEnvVar{{Key: "SERVER_HOST", Value: "127.0.0.1"}, {Key: "SERVER_PORT", Value: "3000"}}},
		// Check for error on empty/invalid input.
		{Input: "", Test: urlTest{Value: "", Error: ErrorURLConnection}},
	}

	// Custom comparer for errors.
	errorComp := cmp.Comparer(errors.Is)

	for _, test := range cases {
		// Set environment variables, store previous values for after tests.
		prevEnv := make([]urlEnvVar, len(test.Env))
		for i, envVar := range test.Env {
			prevEnv[i] = urlEnvVar{Key: envVar.Key, Value: os.Getenv(envVar.Key)}
			os.Setenv(envVar.Key, envVar.Value)
		}

		// Test.
		url, err := BuildConnectionURL(test.Input)

		if diff := cmp.Diff(test.Test, urlTest{Value: url, Error: err}, errorComp); diff != "" {
			t.Fatalf("Failed for BuildConnectionURL() (-want, +got):\n%s", diff)
		}

		// Reset environment variables.
		for _, prevVar := range prevEnv {
			os.Setenv(prevVar.Key, prevVar.Value)
		}
	}
}
