package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in database.
type User struct {
	ID        string         `json:"id" xml:"id" form:"id" gorm:"primaryKey"`
	Username  string         `json:"username" xml:"username" form:"username" gorm:"unique;size:127"`
	Password  string         `json:"-" xml:"-" form:"password" gorm:"index;size=128"` // SHA512
	Lastname  string         `json:"lastname" xml:"lastname" form:"lastname" gorm:"size=63"`
	Firstname string         `json:"firstname" xml:"firstname" form:"firstname" gorm:"size=63"`
	CreatedAt time.Time      `json:"created_at" xml:"created_at" form:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" xml:"updated_at" form:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" xml:"-" form:"deleted_at" gorm:"index"`
}
