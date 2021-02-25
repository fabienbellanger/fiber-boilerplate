package api

import (
	"crypto/sha512"
	"encoding/hex"
	"github.com/fabienbellanger/fiber-boilerplate/db"
	"github.com/fabienbellanger/fiber-boilerplate/models"
	"github.com/gofiber/fiber/v2"
)

type userForm struct {
	Username  string `json:"username" xml:"username" form:"username"`
	Password  string `json:"password" xml:"password" form:"password"`
	Lastname  string `json:"lastname" xml:"lastname" form:"lastname"`
	Firstname string `json:"firstname" xml:"firstname" form:"firstname"`
}

// GetAllUsers lists all users.
func GetAllUsers(db *db.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var users []models.User
		if response := db.Find(&users); response.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    fiber.StatusInternalServerError,
				"message": "Error during users list",
			})
		}

		return c.JSON(users)
	}
}

// CreateUser creates a new user.
func CreateUser(db *db.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := new(userForm)
		if err := c.BodyParser(user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Bad Request",
			})
		}

		// Data validation
		// ---------------
		if user.Firstname == "" || user.Lastname == "" || user.Username == "" || user.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Bad Parameters",
			})
		}

		// Hash password
		// -------------
		passwordBytes := sha512.Sum512([]byte(user.Password))
		hashPassword := hex.EncodeToString(passwordBytes[:])

		// Database insertion
		// ------------------
		newUser := models.User{
			Lastname:  user.Lastname,
			Firstname: user.Firstname,
			Password:  hashPassword,
			Username:  user.Username,
		}

		if result := db.Create(&newUser); result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    fiber.StatusInternalServerError,
				"message": "Error during user creation",
			})
		}

		return c.JSON(newUser)
	}
}
