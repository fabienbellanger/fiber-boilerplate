package db

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"time"

	"github.com/fabienbellanger/fiber-boilerplate/models"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/prometheus"
)

// TODO: Add a custom logger for GORM : https://www.soberkoder.com/go-gorm-logging/

// DatabaseConfig represents the database configuration.
type DatabaseConfig struct {
	Driver          string // Not used
	Host            string
	Username        string
	Password        string
	Port            int
	Database        string
	Charset         string
	Collation       string
	Location        string
	MaxIdleConns    int           // Sets the maximum number of connections in the idle connection pool
	MaxOpenConns    int           // Sets the maximum number of open connections to the database
	ConnMaxLifetime time.Duration // Sets the maximum amount of time a connection may be reused
}

// DB represents the database.
type DB struct {
	*gorm.DB
}

// New makes the connection to the database.
// TODO:
// - Mettre à jour la doc
// - logger ORM sur une sortie au choix (.env)
func New(config *DatabaseConfig) (*DB, error) {
	dsn, err := config.dsn()
	if err != nil {
		return nil, err
	}

	// GORM log configuration
	// ----------------------
	level := getGormLogLevel(viper.GetString("GORM_LEVEL"), viper.GetString("APP_ENV"))
	output := setGormLogOutput(viper.GetString("GORM_OUTPUT"), viper.GetString("GORM_LOG_FILE_PATH"))

	// Logger
	// ------
	newLogger := logger.New(
		log.New(output, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             200 * time.Millisecond, // Slow SQL threshold (Default: 200ms)
			LogLevel:                  level,                  // Log level (Silent, Error, Warn, Info) (Default: Warn)
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
		RefreshInterval: 60,                             // Refresh metrics interval (default 15 seconds)
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

// getGormLogLevel returns the log level for GORM.
// If APP_ENV is development, the default log level is info,
// warn in other case.
// TODO: Add test.
func getGormLogLevel(level, env string) logger.LogLevel {
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
		if env == "development" {
			return logger.Warn
		}
		return logger.Error
	}
}

// TODO: Ne fonctionne pas avec le type file.
func setGormLogOutput(output, filePath string) (file io.Writer) {
	switch output {
	case "file":
		f, _ := os.OpenFile(path.Clean(filePath), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		defer f.Close()
		file = f
	case "stdout":
		file = os.Stdout
	default:
		file = os.Stderr
	}
	return
}

// dsn returns the DSN if the configuration is OK or an error in other case.
// TODO: Add test.
func (c *DatabaseConfig) dsn() (dsn string, err error) {
	if c.Driver == "" || c.Host == "" || c.Database == "" || c.Port == 0 || c.Username == "" || c.Password == "" {
		return dsn, errors.New("error in database configuration")
	}

	dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database)
	if c.Charset != "" {
		dsn += fmt.Sprintf("&charset=%s", c.Charset)
	}
	if c.Collation != "" {
		dsn += fmt.Sprintf("&collation=%s", c.Collation)
	}
	if c.Location != "" {
		dsn += fmt.Sprintf("&loc=%s", c.Location)
	}
	return
}
