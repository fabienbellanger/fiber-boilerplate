package tests

import (
	"io"
	"net/http"
	"strings"
	"testing"

	server "github.com/fabienbellanger/fiber-boilerplate"
	"github.com/stretchr/testify/assert"
)

type header struct {
	key   string
	value string
}

// Define a structure for specifying input and output data of a single test case
type test struct {
	description string

	// Test input
	route   string
	method  string
	body    io.Reader
	headers []header

	// Expected output
	expectedError bool
	expectedCode  int
	expectedBody  string
}

func TestWebRoutes(t *testing.T) {
	tests := []test{
		{
			description:   "Health Check route",
			route:         "/health-check",
			method:        "GET",
			expectedError: false,
			expectedCode:  200,
			expectedBody:  "OK",
		},
		{
			description: "Non existing route",
			route:       "/not-exists",
			method:      "GET",
			body:        strings.NewReader("v=1"),
			headers: []header{
				{key: "Content-Type", value: "application/x-www-form-urlencoded"},
			},
			expectedError: false,
			expectedCode:  401,
			expectedBody:  "{\"code\":401,\"message\":\"Unauthorized\"}",
		},
	}

	// Setup the app as it is done in the main function
	app := server.Setup(nil, nil)

	// Iterate through test single test cases
	for _, test := range tests {
		// Create a new http request with the route from the test case
		req, _ := http.NewRequest(test.method, test.route, test.body)
		for _, h := range test.headers {
			req.Header.Add(h.key, h.value)
		}

		// Perform the request plain with the app.
		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		// Verify that no error occured, that is not expected
		assert.Equalf(t, test.expectedError, err != nil, test.description)

		// As expected errors lead to broken responses, the next test case needs to be processed
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)

		// Read the response body
		body, err := io.ReadAll(res.Body)

		// Reading the response body should work everytime, such that
		// the err variable should be nil
		assert.Nilf(t, err, test.description)

		// Verify, that the reponse body equals the expected body
		assert.Equalf(t, test.expectedBody, string(body), test.description)
	}
}
