package tests

import (
	"io/ioutil"
	"net/http"
	"testing"

	server "github.com/fabienbellanger/fiber-boilerplate"
	"github.com/stretchr/testify/assert"
)

// Example: https://github.com/gofiber/recipes/blob/master/unit-test/main_test.go

func TestWebRoutes(t *testing.T) {
	// Setup the app as it is done in the main function
	app := server.Setup(nil, nil)

	// Create a new http request with the route
	// from the test case
	req, _ := http.NewRequest("GET", "/health-check", nil)

	// Perform the request plain with the app.
	// The -1 disables request latency.
	res, err := app.Test(req, -1)

	assert.Nilf(t, err, "Health Check route")

	// Read the response body
	body, err := ioutil.ReadAll(res.Body)

	assert.Nilf(t, err, "Health Check route")

	// Verify, that the reponse body equals the expected body
	assert.Equalf(t, "OK", string(body), "Health Check route")
}
