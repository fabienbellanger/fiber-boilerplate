package api

import (
	"bufio"
	"encoding/json"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/spf13/viper"

	"github.com/fabienbellanger/fiber-boilerplate/db"
	models "github.com/fabienbellanger/fiber-boilerplate/models"
	"github.com/fabienbellanger/fiber-boilerplate/repositories"
	"github.com/fabienbellanger/fiber-boilerplate/utils"
)

type userLogin struct {
	models.User
	Token     string `json:"token" xml:"token" form:"token"`
	ExpiresAt string `json:"expires_at" xml:"expires_at" form:"expires_at"`
}

type userAuth struct {
	Username string `json:"username" xml:"username" form:"username"`
	Password string `json:"password" xml:"password" form:"password"`
}

// Login authenticates a user.
func Login(db *db.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		u := new(userAuth)
		if err := c.BodyParser(u); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Bad Request",
			})
		}

		user, err := repositories.Login(db, u.Username, u.Password)
		if err != nil {
			log.Printf("Login error: %v\n", err)
			// TODO: Gérer les cas d'erreur (par exemple si la connexion à la BDD est rompue)
			// Différence entre "record not found" et "driver: bad connection".
			return c.Status(fiber.StatusUnauthorized).JSON(utils.HTTPError{
				Code:    fiber.StatusUnauthorized,
				Message: "Unauthorized",
			})
		}

		// Create token
		token := jwt.New(jwt.SigningMethodHS512)

		// Expiration time
		now := time.Now()
		expiresAt := now.Add(time.Hour * viper.GetDuration("JWT_LIFETIME"))

		// Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["id"] = user.ID
		claims["username"] = user.Username
		claims["lastname"] = user.Lastname
		claims["firstname"] = user.Firstname
		claims["createdAt"] = user.CreatedAt
		claims["exp"] = expiresAt.Unix()
		claims["iat"] = now.Unix()
		claims["nbf"] = now.Unix()

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte(viper.GetString("JWT_SECRET")))
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error during token generation")
		}

		return c.JSON(userLogin{
			User:      user,
			Token:     t,
			ExpiresAt: expiresAt.Format("2006-01-02T15:04:05.000Z"),
		})
	}
}

// GetAllUsers lists all users.
func GetAllUsers(db *db.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		users, err := repositories.ListAllUsers(db)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error during users list")
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
		user := new(models.UserForm)
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
		newUser := models.User{
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

		user := new(models.UserForm)
		if err := c.BodyParser(user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Bad Data",
			})
		}

		updatedUser, err := repositories.UpdateUser(db, id, user)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error when deleting user")
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
				user := models.User{
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
