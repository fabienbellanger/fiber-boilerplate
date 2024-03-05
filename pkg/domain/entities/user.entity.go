package entities

import (
	"github.com/fabienbellanger/fiber-boilerplate/utils"
	"github.com/spf13/viper"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// User represents a user in database.
type User struct {
	ID            string         `json:"id" xml:"id" form:"id" gorm:"primaryKey" validate:"required,uuid"`
	Username      string         `json:"username" xml:"username" form:"username" gorm:"not null;unique;size:127" validate:"required,email"`
	Password      string         `json:"-" xml:"-" form:"password" gorm:"not null;index;size:128" validate:"required,min=8"` // SHA512
	Lastname      string         `json:"lastname" xml:"lastname" form:"lastname" gorm:"size:63" validate:"required"`
	Firstname     string         `json:"firstname" xml:"firstname" form:"firstname" gorm:"size:63" validate:"required"`
	CreatedAt     time.Time      `json:"created_at" xml:"created_at" form:"created_at" gorm:"not null;autoCreateTime"`
	UpdatedAt     time.Time      `json:"updated_at" xml:"updated_at" form:"updated_at" gorm:"not null;autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"-" xml:"-" form:"deleted_at" gorm:"index"`
	PasswordReset PasswordResets `json:"-" xml:"-" form:"-" gorm:"constraint:OnDelete:CASCADE"`
}

// PasswordResets is used to reset user password.
type PasswordResets struct {
	UserID    string    `json:"user_id" xml:"user_id" form:"user_id" gorm:"primaryKey" validate:"required,uuid"`
	Token     string    `json:"token" xml:"token" form:"token" gorm:"size:36;not null" validate:"required,uuid"`
	ExpiredAt time.Time `json:"expired_at" xml:"expired_at" gorm:"not null" form:"expired_at"`
}

// GenerateJWT returns a token
func (u *User) GenerateJWT(lifetime time.Duration, algo, secret string) (string, time.Time, error) {
	// Create token and key
	token, key, err := utils.GetTokenAndKeyFromAlgo(algo, secret, viper.GetString("JWT_PRIVATE_KEY_PATH"))
	if err != nil {
		return "", time.Now(), err
	}

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

	// Generate encoded token and send it as response
	t, err := token.SignedString(key)
	if err != nil {
		return "", expiresAt, err
	}

	return t, expiresAt, nil
}
