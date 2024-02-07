package handlers

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/fabienbellanger/fiber-boilerplate/pkg/domain/entities"
	"github.com/fabienbellanger/fiber-boilerplate/pkg/domain/repositories"
	"github.com/fabienbellanger/fiber-boilerplate/pkg/domain/requests"
	"github.com/fabienbellanger/fiber-boilerplate/pkg/domain/usecases"
	"text/template"
	"time"

	"github.com/fabienbellanger/fiber-boilerplate/utils"
	"github.com/fabienbellanger/goutils/mail"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// User handler
type User struct {
	router      fiber.Router
	userUseCase usecases.User
	store       repositories.UserRepository // TODO: Remove this field
	logger      *zap.Logger
}

// NewUser returns a new Handler
func NewUser(r fiber.Router, userUseCase usecases.User, user repositories.UserRepository, logger *zap.Logger) User {
	return User{
		router:      r,
		userUseCase: userUseCase,
		store:       user, // TODO: Remove this field
		logger:      logger,
	}
}

// UserRoutes adds users routes
func (u *User) UserRoutes() {
	u.router.Get("", u.getAll())
	u.router.Get("/:id", u.getByID())
	u.router.Put("/:id", u.update())
	u.router.Delete("/:id", u.delete())
}

// Login authenticates a user.
func (u *User) Login(c *fiber.Ctx) error {
	req := new(requests.UserLogin)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid body",
		})
	}

	res, err := u.userUseCase.Login(*req)
	if err != nil {
		if errors.Is(err, utils.HTTPError{}) && err.Err != nil {
			if details, ok := err.Details.(string); ok {
				return utils.NewError(c, u.logger, err.Message, details, err.Err)
			}
		}
		return c.Status(err.Code).JSON(err)
	}

	return c.JSON(res)
}

// Create creates a new user.
func (u *User) Create(c *fiber.Ctx) error {
	user := new(requests.UserEdit)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: "Bad Request",
		})
	}

	res, err := u.userUseCase.Create(*user)
	if err != nil {
		if errors.Is(err, utils.HTTPError{}) && err.Err != nil {
			if details, ok := err.Details.(string); ok {
				return utils.NewError(c, u.logger, err.Message, details, err.Err)
			}
		}
		return c.Status(err.Code).JSON(err)
	}

	return c.JSON(res)
}

// getAll lists all users.
func (u *User) getAll() fiber.Handler {
	return func(c *fiber.Ctx) error {
		res, err := u.userUseCase.GetAll()
		if err != nil {
			if errors.Is(err, utils.HTTPError{}) && err.Err != nil {
				if details, ok := err.Details.(string); ok {
					return utils.NewError(c, u.logger, err.Message, details, err.Err)
				}
			}
			return c.Status(err.Code).JSON(err)
		}

		return c.JSON(res)
	}
}

// getByID return a user.
func (u *User) getByID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Bad ID",
			})
		}

		userID := requests.UserByID{ID: id}

		user, err := u.userUseCase.GetByID(userID)
		if err != nil {
			if errors.Is(err, utils.HTTPError{}) && err.Err != nil {
				if details, ok := err.Details.(string); ok {
					return utils.NewError(c, u.logger, err.Message, details, err.Err)
				}
			}
			return c.Status(err.Code).JSON(err)
		}

		return c.JSON(user)
	}
}

// delete return a user.
func (u *User) delete() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Bad ID",
			})
		}

		userID := requests.UserByID{ID: id}

		err := u.userUseCase.Delete(userID)
		if err != nil {
			if errors.Is(err, utils.HTTPError{}) && err.Err != nil {
				if details, ok := err.Details.(string); ok {
					return utils.NewError(c, u.logger, err.Message, details, err.Err)
				}
			}
			return c.Status(err.Code).JSON(err)
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

// update updates user information.
func (u *User) update() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Bad ID",
			})
		}

		user := new(requests.UserEdit)
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

		updatedUser, err := u.store.Update(id, *user)
		if err != nil {
			return utils.NewError(c, u.logger, "Database error", "Error when updating user", err)
		}

		return c.JSON(updatedUser)
	}
}

// UpdatePassword updates user password.
func (u *User) UpdatePassword(c *fiber.Ctx) error {
	token := c.Params("token")

	newPassword := new(entities.UserUpdatePassword)
	if err := c.BodyParser(newPassword); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: "Bad Request",
		})
	}

	// Data validation
	// ---------------
	createErrors := utils.ValidateStruct(*newPassword)
	if createErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: "Bad Request",
			Details: createErrors,
		})
	}

	// Update user password
	// --------------------
	userID, currentPassword, err := u.store.GetIDFromPasswordReset(token, newPassword.Password)
	if err != nil {
		return utils.NewError(c, u.logger, "Database error", "Error when searching user", err)
	}
	if userID == "" {
		return fiber.NewError(fiber.StatusNotFound, "no user found")
	}

	// Change by the same password is forbidden
	hashedPassword := sha512.Sum512([]byte(newPassword.Password))
	if hex.EncodeToString(hashedPassword[:]) == currentPassword {
		return fiber.NewError(fiber.StatusBadRequest, "new password cannot be the same as the current one")
	}

	err = u.store.UpdatePassword(userID, currentPassword, newPassword.Password)
	if err != nil {
		return utils.NewError(c, u.logger, "Database error", "Error when updating user password", err)
	}

	// Delete password reset
	// ---------------------
	err = u.store.DeletePasswordReset(userID)
	if err != nil {
		return utils.NewError(c, u.logger, "Database error", "Error when deleting user password reset", err)
	}

	return c.SendStatus(fiber.StatusOK)
}

// ForgottenPassword save a forgotten password request.
func (u *User) ForgottenPassword(c *fiber.Ctx) error {
	// Find user
	user, err := u.store.GetByUsername(c.Params("email"))
	if err != nil {
		return utils.NewError(c, u.logger, "Database error", "Error when retrieving user", err)
	}
	if user.ID == "" {
		return c.Status(fiber.StatusNotFound).JSON(utils.HTTPError{
			Code:    fiber.StatusNotFound,
			Message: "No user found",
		})
	}

	// Sale line in database
	passwordReset := entities.PasswordResets{
		UserID:    user.ID,
		Token:     uuid.NewString(),
		ExpiredAt: time.Now().Add(viper.GetDuration("FORGOTTEN_PASSWORD_EXPIRATION_DURATION") * time.Hour).UTC(),
	}
	err = u.store.CreateOrUpdatePasswordReset(passwordReset)
	if err != nil {
		return utils.NewError(c, u.logger, "Database error", "Error when requesting new password", err)
	}

	// Send email with link
	to := make([]string, 1)
	to[0] = user.Username
	subject := fmt.Sprintf("[%s] Forgotten password", viper.GetString("APP_NAME"))
	var body bytes.Buffer

	tp, err := template.ParseFiles("templates/forgotten_password.gohtml")
	if err != nil {
		return utils.NewError(c, u.logger, "Database error", "Error when creating password reset email", err)
	}
	err = tp.Execute(&body, struct {
		Title string
		Link  string
	}{
		Title: fmt.Sprintf("%s - Forgotten password", viper.GetString("APP_NAME")),
		Link:  fmt.Sprintf("%s/%s", viper.GetString("FORGOTTEN_PASSWORD_BASE_URL"), passwordReset.Token),
	})
	if err != nil {
		return utils.NewError(c, u.logger, "Database error", "Error when creating password reset email", err)
	}

	err = mail.Send(
		viper.GetString("FORGOTTEN_PASSWORD_EMAIL_FROM"),
		to,
		nil,
		nil,
		subject,
		body.String(),
		"",
		"",
		viper.GetString("SMTP_HOST"),
		viper.GetInt("SMTP_PORT"))
	if err != nil {
		return utils.NewError(c, u.logger, "Database error", "Error when sending password reset email", err)
	}

	return c.JSON(passwordReset)
}
