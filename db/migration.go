package db

import "github.com/fabienbellanger/fiber-boilerplate/entities"

// entitiesList lists all entities to automigrate.
var entitiesList = []interface{}{
	&entities.User{},
	&entities.PasswordResets{},
	&entities.Task{},
}
