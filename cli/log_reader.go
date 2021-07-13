package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

type errorLog struct {
	Level   string    `json:"level"`
	Time    time.Time `json:"time"`
	Caller  string    `json:"caller"`
	Message string    `json:"message"`
}

func init() {
	rootCmd.AddCommand(logReaderCmd)
}

var logReaderCmd = &cobra.Command{
	Use:   "log-reader",
	Short: "Log reader",
	Long:  `Log reader`,
	Run: func(cmd *cobra.Command, args []string) {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			// Line to parse and display
			fmt.Println(parseLine(scanner.Bytes()))
		}

		if err := scanner.Err(); err != nil {
			log.Println(err)
		}
	},
}

func parseLine(line []byte) (string, error) {
	var errLog errorLog
	err := json.Unmarshal(line, &errLog)
	if err != nil {
		return "", err
	}
	log.Printf("%#v\n", errLog)
	return "> " + string(line), nil
}
