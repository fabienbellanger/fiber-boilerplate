package main

import (
	"log"
	"time"

	server "github.com/fabienbellanger/fiber-boilerplate"
	"github.com/fabienbellanger/fiber-boilerplate/db"
	"github.com/fabienbellanger/fiber-boilerplate/models"
	"github.com/spf13/viper"
)

func main() {
	// Configuration initialization
	// ----------------------------
	if err := initConfig(); err != nil {
		log.Fatalln(err)
	}

	// Database connection
	// -------------------
	db, err := db.New(&db.DatabaseConfig{
		Driver:          viper.GetString("DB_DRIVER"),
		Host:            viper.GetString("DB_HOST"),
		Username:        viper.GetString("DB_USERNAME"),
		Password:        viper.GetString("DB_PASSWORD"),
		Port:            viper.GetInt("DB_PORT"),
		Database:        viper.GetString("DB_DATABASE"),
		MaxIdleConns:    viper.GetInt("DB_MAX_IDLE_CONNS"),
		MaxOpenConns:    viper.GetInt("DB_MAX_OPEN_CONNS"),
		ConnMaxLifetime: viper.GetDuration("DB_CONN_MAX_LIFETIME") * time.Hour,
	})
	if err != nil {
		log.Fatalln(err)
	}

	// Database migrations
	// -------------------
	db.AutoMigrate(&models.User{})

	// Start server
	// ------------
	server.Run(db)
}

// initConfig initializes configuration from config file.
func initConfig() error {
	viper.SetConfigFile(".env")
	return viper.ReadInConfig()
}
