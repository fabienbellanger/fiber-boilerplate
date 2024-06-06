package db

import (
	"github.com/fabienbellanger/fiber-boilerplate/pkg/domain/entities"
)

// entitiesList lists all entities to auto migrate.
var entitiesList = []interface{}{
	&entities.User{},
	&entities.PasswordResets{},
	&entities.Task{},
}

var migrations = []func(db *DB) error{
	changeDescriptionTaskColumn,
}

// changeDescriptionTaskColumn, adds the state column to the tasks table.
func changeDescriptionTaskColumn(db *DB) error {
	if !db.Migrator().HasColumn(&entities.Task{}, "state") {
		tx := db.Exec(`
			ALTER TABLE tasks MODIFY COLUMN description VARCHAR(255);
		`)

		return tx.Error
	}

	return nil
}
