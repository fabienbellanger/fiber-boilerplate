package responses

import "github.com/fabienbellanger/fiber-boilerplate/pkg/domain/entities"

// UserLogin response
type UserLogin struct {
	entities.User
	Token     string `json:"token" xml:"token" form:"token"`
	ExpiresAt string `json:"expires_at" xml:"expires_at" form:"expires_at"`
}
