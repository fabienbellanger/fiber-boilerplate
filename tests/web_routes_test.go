package tests

import (
	"strings"
	"testing"
)

func TestWebRoutes(t *testing.T) {
	Init("../.env")

	useCases := []Test{
		{
			Description:  "Health Check route",
			Route:        "/health-check",
			Method:       "GET",
			CheckCode:    true,
			CheckBody:    true,
			ExpectedCode: 200,
			ExpectedBody: "OK",
		},
		{
			Description: "Non existing route",
			Route:       "/not-exists",
			Method:      "GET",
			Body:        strings.NewReader("v=1"),
			Headers: []Header{
				{Key: "Content-Type", Value: "application/x-www-form-urlencoded"},
			},
			CheckCode:     true,
			CheckBody:     true,
			ExpectedError: false,
			ExpectedCode:  401,
			ExpectedBody:  "{\"code\":401,\"message\":\"Unauthorized\"}",
		},
	}

	Execute(t, nil, useCases, "../templates")
}
