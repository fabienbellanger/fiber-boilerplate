package tests

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fabienbellanger/fiber-boilerplate/db"
	"github.com/spf13/viper"
)

// TestDB is used to create and use a database for tests.
type TestDB struct {
	name string
	db   *db.DB
}

// NewTestDB return a TestDB instance.
func NewTestDB() (TestDB, error) {
	rand.Seed(time.Now().UnixNano())

	config := db.DatabaseConfig{
		Driver:          viper.GetString("DB_DRIVER"),
		Host:            viper.GetString("DB_HOST"),
		Username:        viper.GetString("DB_USERNAME"),
		Password:        viper.GetString("DB_PASSWORD"),
		Port:            viper.GetInt("DB_PORT"),
		Database:        viper.GetString("DB_DATABASE") + "_" + fmt.Sprintf("%08d", rand.Int63n(1e8)),
		Charset:         viper.GetString("DB_CHARSET"),
		Collation:       viper.GetString("DB_COLLATION"),
		Location:        viper.GetString("DB_LOCATION"),
		MaxIdleConns:    viper.GetInt("DB_MAX_IDLE_CONNS"),
		MaxOpenConns:    viper.GetInt("DB_MAX_OPEN_CONNS"),
		ConnMaxLifetime: viper.GetDuration("DB_CONN_MAX_LIFETIME") * time.Hour,
	}

	db, err := db.New(&config)
	if err != nil {
		return TestDB{}, err
	}
	return TestDB{db: db, name: config.Database}, nil
}

// Drop database after the test.
func (d *TestDB) Drop() error {

	return nil
}
