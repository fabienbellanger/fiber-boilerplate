package api

import (
	"strconv"

	"github.com/fabienbellanger/fiber-boilerplate/db"
	"github.com/fabienbellanger/fiber-boilerplate/models"
	"github.com/fabienbellanger/fiber-boilerplate/repositories"
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
		users, err := repositories.ListAllUsers(db)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    fiber.StatusInternalServerError,
				"message": "Error during users list",
			})
		}

		return c.JSON(users)
	}
}

// GetUser return a user.
func GetUser(db *db.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil || id <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Bad ID",
			})
		}

		user, err := repositories.GetUser(db, uint(id))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    fiber.StatusInternalServerError,
				"message": "Error when retrieving user",
			})
		}
		if user.ID == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"code":    fiber.StatusNotFound,
				"message": "No user found",
			})
		}

		return c.JSON(user)
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

		// Database insertion
		// ------------------
		newUser := models.User{
			Lastname:  user.Lastname,
			Firstname: user.Firstname,
			Password:  user.Password,
			Username:  user.Username,
		}

		if err := repositories.CreateUser(db, &newUser); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    fiber.StatusInternalServerError,
				"message": "Error during user creation",
			})
		}

		return c.JSON(newUser)
	}
}

// DeleteUser return a user.
func DeleteUser(db *db.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil || id <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Bad ID",
			})
		}

		err = repositories.DeleteUser(db, uint(id))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    fiber.StatusInternalServerError,
				"message": "Error when deleting user",
			})
		}

		return c.JSON(fiber.Map{
			"code":    fiber.StatusOK,
			"message": "User deleted",
		})
	}
}
