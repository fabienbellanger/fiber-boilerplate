package db

import (
	entities2 "github.com/fabienbellanger/fiber-boilerplate/pkg/domain/entities"
)

// entitiesList lists all entities to auto migrate.
var entitiesList = []interface{}{
	&entities2.User{},
	&entities2.PasswordResets{},
	&entities2.Task{},
}
