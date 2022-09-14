package user

import (
	"bufio"
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"text/template"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"

	entities "github.com/fabienbellanger/fiber-boilerplate/entities"
	"github.com/fabienbellanger/fiber-boilerplate/stores"
	"github.com/fabienbellanger/fiber-boilerplate/utils"
	"github.com/fabienbellanger/goutils/mail"
)

type UserHandler struct {
	router fiber.Router
	store  stores.UserStorer
	logger *zap.Logger
}

// New returns a new UserHandler
func New(r fiber.Router, user stores.UserStorer, logger *zap.Logger) UserHandler {
	return UserHandler{
		router: r,
		store:  user,
		logger: logger,
	}
}

type userLogin struct {
	entities.User
	Token     string `json:"token" xml:"token" form:"token"`
	ExpiresAt string `json:"expires_at" xml:"expires_at" form:"expires_at"`
}

type userAuth struct {
	Username string `json:"username" xml:"username" form:"username" validate:"required,email"`
	Password string `json:"password" xml:"password" form:"password" validate:"required,min=8"`
}

// Routes adds users routes
func (u *UserHandler) Routes() {
	u.router.Get("", u.getAll())
	u.router.Get("/stream", u.stream())
	u.router.Get("/:id", u.getOne())
	u.router.Put("/:id", u.update())
	u.router.Delete("/:id", u.delete())
}

// Login authenticates a user.
func (u *UserHandler) Login(c *fiber.Ctx) error {
	ua := new(userAuth)
	if err := c.BodyParser(ua); err != nil {
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

	user, err := u.store.Login(ua.Username, ua.Password)
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

// Create creates a new user.
func (u *UserHandler) Create(c *fiber.Ctx) error {
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

	if err := u.store.Create(&newUser); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error during user creation")
	}
	return c.JSON(newUser)
}

// getAll lists all users.
func (u *UserHandler) getAll() fiber.Handler {
	return func(c *fiber.Ctx) error {
		users, err := u.store.GetAll()
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.JSON(users)
	}
}

// getOne return a user.
func (u *UserHandler) getOne() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Bad ID",
			})
		}

		user, err := u.store.GetOne(id)
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

// delete return a user.
func (u *UserHandler) delete() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Bad ID",
			})
		}

		err := u.store.Delete(id)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error when deleting user")
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

// update updates user information.
func (u *UserHandler) update() fiber.Handler {
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

		updatedUser, err := u.store.Update(id, user)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error when updating user")
		}

		return c.JSON(updatedUser)
	}
}

// UpdatePassword updates user password.
func (u *UserHandler) UpdatePassword(c *fiber.Ctx) error {
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
		return fiber.NewError(fiber.StatusInternalServerError, "Error when searching user")
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
		return fiber.NewError(fiber.StatusInternalServerError, "Error when updating user password")
	}

	// Delete password reset
	// ---------------------
	err = u.store.DeletePasswordReset(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error when deleting user password reset")
	}

	return c.SendStatus(fiber.StatusOK)
}

// ForgottenPassword save a forgotten password request.
func (u *UserHandler) ForgottenPassword(c *fiber.Ctx) error {
	// Find user
	user, err := u.store.GetByUsername(c.Params("email"))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error when retrieving user")
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
		Token:     uuid.New().String(),
		ExpiredAt: time.Now().Add(viper.GetDuration("FORGOTTEN_PASSWORD_EXPIRATION_DURATION") * time.Hour).UTC(),
	}
	err = u.store.CreateOrUpdatePasswordReset(&passwordReset)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error when requesting new password")
	}

	// TODO: Send email with link
	to := make([]string, 1)
	to[0] = user.Username
	subject := fmt.Sprintf("[%s] Forgotten password", viper.GetString("APP_NAME"))
	var body bytes.Buffer

	tp, err := template.ParseFiles("templates/forgotten_password.gohtml")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error when creating password reset email")
	}
	err = tp.Execute(&body, struct {
		Title string
		Link  string
	}{
		Title: fmt.Sprintf("%s - Forgotten password", viper.GetString("APP_NAME")),
		Link:  fmt.Sprintf("%s/%s", viper.GetString("FORGOTTEN_PASSWORD_BASE_URL"), passwordReset.Token),
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error when creating password reset email")
	}

	err = mail.Send(viper.GetString("FORGOTTEN_PASSWORD_EMAIL_FROM"), to, subject, body.String(), "", "", viper.GetString("SMTP_HOST"), viper.GetInt("SMTP_PORT"))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error when sending password reset email")
	}

	return c.JSON(passwordReset)
}

// stream returns users list in a stream.
func (u *UserHandler) stream() fiber.Handler {
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
