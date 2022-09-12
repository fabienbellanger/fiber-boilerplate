package api

import (
	"bufio"
	"encoding/json"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"gorm.io/gorm"

	"github.com/fabienbellanger/fiber-boilerplate/db"
	entities "github.com/fabienbellanger/fiber-boilerplate/entities"
	"github.com/fabienbellanger/fiber-boilerplate/repositories"
	"github.com/fabienbellanger/fiber-boilerplate/utils"
)

type userLogin struct {
	entities.User
	Token     string `json:"token" xml:"token" form:"token"`
	ExpiresAt string `json:"expires_at" xml:"expires_at" form:"expires_at"`
}

type userAuth struct {
	Username string `json:"username" xml:"username" form:"username" validate:"required,email"`
	Password string `json:"password" xml:"password" form:"password" validate:"required,min=8"`
}

// Login authenticates a user.
func Login(db *db.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		u := new(userAuth)
		if err := c.BodyParser(u); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Invalid body",
			})
		}

		loginErrors := utils.ValidateStruct(*u)
		if loginErrors != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Invalid body",
				Details: loginErrors,
			})
		}

		user, err := repositories.Login(db, u.Username, u.Password)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusUnauthorized).JSON(utils.HTTPError{
					Code:    fiber.StatusUnauthorized,
					Message: "Unauthorized",
				})
			}
			return fiber.NewError(fiber.StatusInternalServerError, "Error during authentication")
		}

		// Create token
		token, expiresAt, err := user.GenerateJWT(viper.GetDuration("JWT_LIFETIME"), viper.GetString("JWT_ALGO"), viper.GetString("JWT_SECRET"))
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error during token generation")
		}

		return c.JSON(userLogin{
			User:      user,
			Token:     token,
			ExpiresAt: expiresAt.Format("2006-01-02T15:04:05.000Z"),
		})
	}
}

// GetAllUsers lists all users.
func GetAllUsers(db *db.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		users, err := repositories.ListAllUsers(db)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.JSON(users)
	}
}

// GetUser return a user.
func GetUser(db *db.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Bad ID",
			})
		}

		user, err := repositories.GetUser(db, id)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error when retrieving user")
		}
		if user.ID == "" {
			return c.Status(fiber.StatusNotFound).JSON(utils.HTTPError{
				Code:    fiber.StatusNotFound,
				Message: "No user found",
			})
		}

		return c.JSON(user)
	}
}

// CreateUser creates a new user.
func CreateUser(db *db.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := new(entities.UserForm)
		if err := c.BodyParser(user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Bad Request",
			})
		}

		// Data validation
		// ---------------
		if user.Firstname == "" || user.Lastname == "" || user.Username == "" || user.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Bad Parameters",
			})
		}

		// Database insertion
		// ------------------
		newUser := entities.User{
			Lastname:  user.Lastname,
			Firstname: user.Firstname,
			Password:  user.Password,
			Username:  user.Username,
		}

		if err := repositories.CreateUser(db, &newUser); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error during user creation")
		}
		return c.JSON(newUser)
	}
}

// DeleteUser return a user.
func DeleteUser(db *db.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Bad ID",
			})
		}

		err := repositories.DeleteUser(db, id)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error when deleting user")
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

// UpdateUser updates user information.
func UpdateUser(db *db.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Bad ID",
			})
		}

		user := new(entities.UserForm)
		if err := c.BodyParser(user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Bad Data",
			})
		}

		updateErrors := utils.ValidateStruct(*user)
		if updateErrors != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Bad Request",
				Details: updateErrors,
			})
		}

		updatedUser, err := repositories.UpdateUser(db, id, user)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error when updating user")
		}

		return c.JSON(updatedUser)
	}
}

// StreamUsers returns users list in a stream.
func StreamUsers(db *db.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)

		c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
			w.WriteString("[")
			enc := json.NewEncoder(w)
			n := 100_000
			for i := 0; i < n; i++ {
				user := entities.User{
					ID:        uuid.New().String(),
					Username:  "My Username",
					Password:  ",kkjkjkjkjknnqfjkkjdnfsjklqblk",
					Lastname:  "My Lastname",
					Firstname: "My Firstname",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				if err := enc.Encode(user); err != nil {
					continue
				}

				if i < n-1 {
					w.WriteString(",")
				}

				w.Flush()
			}
			w.WriteString("]")
		})

		return nil
	}
}
