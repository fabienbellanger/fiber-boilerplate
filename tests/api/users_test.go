package api

import (
	"strings"
	"testing"

	"github.com/fabienbellanger/fiber-boilerplate/deliveries/user"
	"github.com/fabienbellanger/fiber-boilerplate/entities"
	"github.com/fabienbellanger/fiber-boilerplate/tests"
	"github.com/gofiber/fiber/v2"
)

func TestUserCreation(t *testing.T) {
	tdb := tests.Init("../../.env")
	defer tdb.Drop()

	useCases := []tests.Test{
		{
			Description: "User creation",
			Route:       "/api/v1/register",
			Method:      "POST",
			Body: strings.NewReader(tests.JsonToString(entities.UserForm{
				Username:  "test1@gmail.com",
				Password:  "11111111",
				Lastname:  "Test",
				Firstname: "Creation",
			})),
			Headers: []tests.Header{
				{Key: "Content-Type", Value: fiber.MIMEApplicationJSONCharsetUTF8},
				{Key: "Authorization", Value: "Bearer " + tdb.Token},
			},
			CheckCode:    true,
			ExpectedCode: 200,
		},
		{
			Description: "User creation with invalid password",
			Route:       "/api/v1/register",
			Method:      "POST",
			Body: strings.NewReader(tests.JsonToString(entities.UserForm{
				Username:  "test1@gmail.com",
				Password:  "1111111",
				Lastname:  "Test",
				Firstname: "Creation",
			})),
			Headers: []tests.Header{
				{Key: "Content-Type", Value: fiber.MIMEApplicationJSONCharsetUTF8},
				{Key: "Authorization", Value: "Bearer " + tdb.Token},
			},
			CheckBody:    true,
			CheckCode:    true,
			ExpectedCode: 400,
			ExpectedBody: `{"code":400,"message":"Invalid body","details":[{"FailedField":"Password","Tag":"min","Value":"8"}]}`,
		},
		{
			Description: "User creation with invalid username",
			Route:       "/api/v1/register",
			Method:      "POST",
			Body: strings.NewReader(tests.JsonToString(entities.UserForm{
				Username:  "test1",
				Password:  "11111111",
				Lastname:  "Test",
				Firstname: "Creation",
			})),
			Headers: []tests.Header{
				{Key: "Content-Type", Value: fiber.MIMEApplicationJSONCharsetUTF8},
				{Key: "Authorization", Value: "Bearer " + tdb.Token},
			},
			CheckBody:    true,
			CheckCode:    true,
			ExpectedCode: 400,
			ExpectedBody: `{"code":400,"message":"Invalid body","details":[{"FailedField":"Username","Tag":"email","Value":""}]}`,
		},
	}

	tests.Execute(t, tdb.DB, useCases, "../../templates")
}

func TestUserLogin(t *testing.T) {
	tdb := tests.Init("../../.env")
	defer tdb.Drop()

	useCases := []tests.Test{
		{
			Description: "User login",
			Route:       "/api/v1/login",
			Method:      "POST",
			Body: strings.NewReader(tests.JsonToString(user.UserAuth{
				Username: tests.UserUsername,
				Password: tests.UserPassword,
			})),
			Headers: []tests.Header{
				{Key: "Content-Type", Value: fiber.MIMEApplicationJSONCharsetUTF8},
			},
			CheckCode:    true,
			ExpectedCode: 200,
		},
	}

	tests.Execute(t, tdb.DB, useCases, "../../templates")
}
