package cli

import (
	"log"
	"time"

	server "github.com/fabienbellanger/fiber-boilerplate"
	"github.com/fabienbellanger/fiber-boilerplate/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "run",
	Short: "Start server",
	Long:  `Start server`,
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

func startServer() {
	// Configuration initialization
	// ----------------------------
	if err := initConfig(); err != nil {
		log.Fatalln(err)
	}

	// Logger initialization
	// ---------------------
	logger, err := server.InitLogger()
	if err != nil {
		log.Fatalln(err)
	}
	defer logger.Sync()

	// Database connection
	// -------------------
	db, err := db.New(&db.DatabaseConfig{
		Driver:          viper.GetString("DB_DRIVER"),
		Host:            viper.GetString("DB_HOST"),
		Username:        viper.GetString("DB_USERNAME"),
		Password:        viper.GetString("DB_PASSWORD"),
		Port:            viper.GetInt("DB_PORT"),
		Database:        viper.GetString("DB_DATABASE"),
		Charset:         viper.GetString("DB_CHARSET"),
		Collation:       viper.GetString("DB_COLLATION"),
		Location:        viper.GetString("DB_LOCATION"),
		MaxIdleConns:    viper.GetInt("DB_MAX_IDLE_CONNS"),
		MaxOpenConns:    viper.GetInt("DB_MAX_OPEN_CONNS"),
		ConnMaxLifetime: viper.GetDuration("DB_CONN_MAX_LIFETIME") * time.Hour,
		SlowThreshold:   viper.GetDuration("GORM_SLOW_THRESHOLD"),
	})
	if err != nil {
		log.Fatalln(err)
	}

	// Database migrations
	// -------------------
	if viper.GetBool("DB_USE_AUTOMIGRATIONS") {
		db.MakeMigrations()
	}

	// Start server
	// ------------
	server.Run(db, logger)
}
