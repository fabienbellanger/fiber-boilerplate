package cli

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(logRotationCmd)
}

var logRotationCmd = &cobra.Command{
	Use:   "log-rotate",
	Short: "Log rotation",
	Long:  `Log rotation`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := initConfig(); err != nil {
			log.Fatalln(err)
		}
		fmt.Println("Log rotation...")
	},
}
