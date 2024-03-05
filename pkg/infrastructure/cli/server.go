package cli

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	server "github.com/fabienbellanger/fiber-boilerplate/pkg/infrastructure/router"
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
	logger, db, err := initConfigLoggerDatabase(true, true)
	if err != nil {
		log.Fatalln(err)
	}

	// Database migrations
	// -------------------
	if viper.GetBool("DB_USE_AUTOMIGRATIONS") {
		err = db.MakeMigrations()
		if err != nil {
			log.Fatalln(err)
		}
	}

	// Start server
	// ------------
	err = server.Run(db, logger, "./templates")
	if err != nil {
		log.Fatalln(err)
	}
}
