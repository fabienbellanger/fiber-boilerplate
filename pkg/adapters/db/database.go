package db

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/prometheus"
)

const (
	// MaxLimit represents the max number of items for pagination
	MaxLimit = 100

	// DefaultSlowThreshold represents the default slow threshold value
	DefaultSlowThreshold time.Duration = 200 * time.Millisecond
)

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
	SlowThreshold   time.Duration // Slow SQL threshold (Default: 200ms)
}

// DB represents the database.
type DB struct {
	*gorm.DB
}

// New makes the connection to the database.
func New(config *DatabaseConfig) (*DB, error) {
	dsn, err := config.dsn()
	if err != nil {
		return nil, err
	}

	if config.SlowThreshold == 0 {
		config.SlowThreshold = DefaultSlowThreshold
	}

	// GORM logger configuration
	// -------------------------
	env := viper.GetString("APP_ENV")
	level := getGormLogLevel(viper.GetString("GORM_LOG_LEVEL"), env)
	output, err := getGormLogOutput(viper.GetString("GORM_LOG_OUTPUT"),
		viper.GetString("GORM_LOG_FILE_PATH"),
		env)
	if err != nil {
		return nil, err
	}

	// Logger
	// ------
	// TODO: Add a custom logger for GORM like https://www.soberkoder.com/go-gorm-logging/
	// Or try something like this: https://github.com/moul/zapgorm2
	customLogger := logger.New(
		log.New(output, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             config.SlowThreshold, // Slow SQL threshold (Default: 200ms)
			LogLevel:                  level,                // Log level (Silent, Error, Warn, Info) (Default: Warn)
			IgnoreRecordNotFoundError: true,                 // Ignore ErrRecordNotFound error for logger (Default: false)
			Colorful:                  true,                 // Disable color (Default: true)
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: customLogger,
	})
	if err != nil {
		return nil, err
	}

	// Options
	// -------
	db.Set("gorm:table_options", "ENGINE=InnoDB")

	// Prometheus
	// ----------
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
func (db *DB) MakeMigrations() error {
	// Auto migrations
	if err := db.AutoMigrate(entitiesList...); err != nil {
		return err
	}

	// Custom migrations
	for _, m := range migrations {
		if err := m(db); err != nil {
			return err
		}
	}

	return nil
}

// getGormLogLevel returns the log level for GORM.
// If APP_ENV is development, the default log level is info,
// warn in other case.
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

// getGormLogOutput returns GORM log output.
// The default value is os.Stdout.
// In development mode, the ouput is set to os.Stdout.
func getGormLogOutput(output, filePath, env string) (file io.Writer, err error) {
	if env == "development" {
		return os.Stdout, nil
	}

	switch output {
	case "file":
		f, err := os.OpenFile(path.Clean(filePath), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
		return f, nil
	default:
		return os.Stdout, nil
	}
}

// dsn returns the DSN if the configuration is OK or an error in other case.
func (c *DatabaseConfig) dsn() (dsn string, err error) {
	if c.Host == "" || c.Port == 0 || c.Username == "" || c.Password == "" || c.Database == "" {
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

// paginateValues transforms page and limit into offset and limit.
func paginateValues(p, l string) (offset int, limit int) {
	page, err := strconv.Atoi(p)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err = strconv.Atoi(l)
	if err != nil || limit > MaxLimit || limit < 1 {
		limit = MaxLimit
	}

	offset = (page - 1) * limit

	return
}

// Paginate creates a GORM scope to paginate queries.
func Paginate(p, l string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset, limit := paginateValues(p, l)

		return db.Offset(offset).Limit(limit)
	}
}

// orderValues transforms list of fields to sort into a map.
func orderValues(list string, prefixes ...string) map[string]string {
	r := make(map[string]string)

	if len(list) <= 0 {
		return r
	}

	prefix := ""
	if len(prefixes) == 1 {
		prefix = prefixes[0] + "."
	}

	sorts := strings.Split(list, ",")
	for _, s := range sorts {
		key := fmt.Sprintf("%s%s", prefix, s[1:])
		if strings.HasPrefix(s, "+") && len(s[1:]) > 1 {
			r[key] = "ASC"
		} else if strings.HasPrefix(s, "-") && len(s[1:]) > 1 {
			r[key] = "DESC"
		}
	}

	return r
}

// Order creates a GORM scope to sort query attributes.
// Example: "+created_at,-id" will produce "ORDER BY created_at ASC, id DESC".
func Order(list string, prefixes ...string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		values := orderValues(list, prefixes...)

		for f, s := range values {
			db.Order(f + " " + s)
		}

		return db
	}
}
