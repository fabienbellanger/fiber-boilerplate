package api

import (
	"errors"
	"github.com/fabienbellanger/fiber-boilerplate/pkg/domain/requests"
	"github.com/fabienbellanger/fiber-boilerplate/pkg/domain/usecases"
	"github.com/fabienbellanger/fiber-boilerplate/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// User handler
type User struct {
	router      fiber.Router
	userUseCase usecases.User
	logger      *zap.Logger
}

// NewUser returns a new Handler
func NewUser(r fiber.Router, userUseCase usecases.User, logger *zap.Logger) User {
	return User{
		router:      r,
		userUseCase: userUseCase,
		logger:      logger,
	}
}

// UserProtectedRoutes adds users protected routes
func (u *User) UserProtectedRoutes() {
	u.router.Get("", u.getAll())
	u.router.Post("", u.create())
	u.router.Get("/:id", u.getByID())
	u.router.Put("/:id", u.update())
	u.router.Delete("/:id", u.delete())
}

// UserPublicRoutes adds users public routes
func (u *User) UserPublicRoutes() {
	u.router.Post("/login", u.login())
	u.router.Post("/forgotten-password/:email", u.forgottenPassword())
	u.router.Patch("/update-password/:token", u.updatePassword())
}

// login authenticates a user.
func (u *User) login() fiber.Handler {
	return func(c *fiber.Ctx) error {
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
}

// create creates a new user.
func (u *User) create() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := new(requests.UserCreation)
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
}

// getAll lists all users.
func (u *User) getAll() fiber.Handler {
	return func(c *fiber.Ctx) error {
		pagination := new(requests.Pagination)
		if err := c.QueryParser(pagination); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Bad Request",
			})
		}

		res, err := u.userUseCase.GetAll(*pagination)
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

		user := new(requests.UserCreation)
		if err := c.BodyParser(user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Bad Data",
			})
		}

		userUpdate := requests.UserUpdate{
			ID:        id,
			Username:  user.Username,
			Password:  user.Password,
			Lastname:  user.Lastname,
			Firstname: user.Firstname,
		}

		res, err := u.userUseCase.Update(userUpdate)
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

// updatePassword updates user password.
func (u *User) updatePassword() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Params("token")

		newPassword := struct {
			Password string `json:"password" xml:"password" form:"password"`
		}{}
		if err := c.BodyParser(&newPassword); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Bad Request",
			})
		}

		password := requests.UserPasswordUpdate{
			Token:    token,
			Password: newPassword.Password,
		}

		err := u.userUseCase.UpdatePassword(password)
		if err != nil {
			if errors.Is(err, utils.HTTPError{}) && err.Err != nil {
				if details, ok := err.Details.(string); ok {
					return utils.NewError(c, u.logger, err.Message, details, err.Err)
				}
			}
			return c.Status(err.Code).JSON(err)
		}

		return c.SendStatus(fiber.StatusOK)
	}
}

// forgottenPassword save a forgotten password request.
func (u *User) forgottenPassword() fiber.Handler {
	return func(c *fiber.Ctx) error {
		email := c.Params("email")
		req := requests.UserForgotPassword{Email: email}

		res, err := u.userUseCase.ForgottenPassword(req)
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
