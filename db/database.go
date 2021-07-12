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
	Charset         string
	Collation       string
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
// TODO: Mettre à jour la doc
func New(config *DatabaseConfig) (*DB, error) {
	// TODO: Vérifier la config et mettre des valeurs par défaut si c'est possible.
	// Sinon retourner une erreur.
	// TODO: Mettre loc en var d'environnement (checker sur le site de GORM les valeurs possibles).
	// Attention à l'encodage !
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&collation=%s&parseTime=True&loc=UTC",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
		config.Charset,
		config.Collation)

	// Logger
	// ------
	logLevel := viper.GetString("GORM_LOG_LEVEL")
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,    // Slow SQL threshold (Default: 200ms)
			LogLevel:                  getGORMLogLevel(logLevel), // Log level (Silent, Error, Warn, Info) (Default: Warn)
			IgnoreRecordNotFoundError: false,                     // Ignore ErrRecordNotFound error for logger (Default: false)
			Colorful:                  true,                      // Disable color (Default: true)
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

// getGORMLogLevel returns the log level for GORM.
// TODO: Use APP_ENV for default case.
func getGORMLogLevel(level string) logger.LogLevel {
	switch level {
	case "silent":
		return logger.Silent
	case "info":
		return logger.Info
	case "warn":
		return logger.Warn
	case "error":
		return logger.Error
	default:
		return logger.Warn
	}
}
