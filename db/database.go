package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fabienbellanger/fiber-boilerplate/models"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/prometheus"
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

	// Logger
	// ------
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             200 * time.Millisecond, // Slow SQL threshold (Default: 200ms)
			LogLevel:                  logger.Warn,            // Log level (Silent, Error, Warn, Info) (Default: Warn)
			IgnoreRecordNotFoundError: false,                  // Ignore ErrRecordNotFound error for logger (Default: false)
			Colorful:                  true,                   // Disable color (Default: true)
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}

	// Options
	// -------
	db.Set("gorm:table_options", "ENGINE=InnoDB")

	// Prometeus
	// ---------
	db.Use(prometheus.New(prometheus.Config{
		DBName:          viper.GetString("DB_DATABASE"), // Use `DBName` as metrics label
		RefreshInterval: 15,                             // Refresh metrics interval (default 15 seconds)
		StartServer:     false,                          // Start http server to expose metrics
		MetricsCollector: []prometheus.MetricsCollector{
			&prometheus.MySQL{
				VariableNames: []string{"Threads_running"},
			},
		}, // user defined metrics
	}))

	// Connection Pool
	// ---------------
	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)

	return &DB{db}, err
}

// MakeMigrations runs GORM migrations.
func (db *DB) MakeMigrations() {
	db.AutoMigrate(&models.User{}, &models.Task{})
}
