package db

import "github.com/fabienbellanger/fiber-boilerplate/models"

// modelsList lists all models to automigrate.
var modelsList = []interface{}{
	&models.User{},
	&models.Task{},
}
