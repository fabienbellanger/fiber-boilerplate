package models

import (
	"time"

	"gorm.io/gorm"
)

// Task represents a task in database.
type Task struct {
	ID          string         `json:"id" xml:"id" form:"id" gorm:"primaryKey"`
	Name        string         `json:"name" xml:"name" form:"name" gorm:"size:127"`
	Description string         `json:"description" xml:"description" form:"description" gorm:"size=127"`
	CreatedAt   time.Time      `json:"created_at" xml:"created_at" form:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" xml:"updated_at" form:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" xml:"-" form:"deleted_at" gorm:"index"`
}

// TaskForm represents structure to create or update a task.
type TaskForm struct {
	Name        string `json:"name" xml:"name" form:"name" validate:"required,min=3,max=127"`
	Description string `json:"description" xml:"description" form:"description"`
}
