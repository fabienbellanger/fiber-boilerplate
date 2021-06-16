package db

import (
	"fmt"
	"time"

	"github.com/fabienbellanger/fiber-boilerplate/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DatabaseConfig represents the database configuration.
type DatabaseConfig struct {
	Driver          string
	Host            string
	Username        string
	Password        string
	Port            int
	Database        string
	MaxIdleConns    int           // Sets the maximum number of connections in the idle connection pool
	MaxOpenConns    int           // Sets the maximum number of open connections to the database
	ConnMaxLifetime time.Duration // Sets the maximum amount of time a connection may be reused
}

// DB represents the database.
type DB struct {
	*gorm.DB
}

// New makes the connection to the database.
func New(config *DatabaseConfig) (*DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=True&loc=UTC",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Options
	// -------
	db.Set("gorm:table_options", "ENGINE=InnoDB")

	// Connection Pool
	// ---------------
	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)

	return &DB{db}, err
}

func (db *DB) MakeMigrations() {
	db.AutoMigrate(&models.User{}, &models.Task{})
}
