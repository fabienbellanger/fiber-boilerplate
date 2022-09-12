package entities

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

// User represents a user in database.
type User struct {
	ID        string         `json:"id" xml:"id" form:"id" gorm:"primaryKey" validate:"required,uuid"`
	Username  string         `json:"username" xml:"username" form:"username" gorm:"unique;size:127" validate:"required,email"`
	Password  string         `json:"-" xml:"-" form:"password" gorm:"index;size=128" validate:"required,min=8"` // SHA512
	Lastname  string         `json:"lastname" xml:"lastname" form:"lastname" gorm:"size=63" validate:"required"`
	Firstname string         `json:"firstname" xml:"firstname" form:"firstname" gorm:"size=63" validate:"required"`
	CreatedAt time.Time      `json:"created_at" xml:"created_at" form:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" xml:"updated_at" form:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" xml:"-" form:"deleted_at" gorm:"index"`
}

// UserForm is used to create or update a user.
type UserForm struct {
	Username  string `json:"username" xml:"username" form:"username" validate:"required,email"`
	Password  string `json:"password" xml:"password" form:"password" validate:"required,min=8"`
	Lastname  string `json:"lastname" xml:"lastname" form:"lastname" validate:"required"`
	Firstname string `json:"firstname" xml:"firstname" form:"firstname" validate:"required"`
}

// GenerateJWT returns a token
// TODO: Add unit tests
func (u *User) GenerateJWT(lifetime time.Duration, algo, secret string) (string, time.Time, error) {
	// Create token
	if algo != "HS512" {
		return "", time.Now(), errors.New("unsupported JWT algo")
	}
	token := jwt.New(jwt.SigningMethodHS512)

	// Expiration time
	now := time.Now()
	expiresAt := now.Add(time.Hour * lifetime)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = u.ID
	claims["username"] = u.Username
	claims["lastname"] = u.Lastname
	claims["firstname"] = u.Firstname
	claims["createdAt"] = u.CreatedAt
	claims["exp"] = expiresAt.Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", expiresAt, err
	}

	return t, expiresAt, nil
}
