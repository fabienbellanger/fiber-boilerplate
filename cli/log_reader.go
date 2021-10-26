package cli

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var serverLogsFlag bool
var dbLogsFlag bool

func init() {
	logReaderCmd.Flags().BoolVarP(&serverLogsFlag, "server", "s", false, "server error logs")
	logReaderCmd.Flags().BoolVarP(&dbLogsFlag, "db", "d", false, "database logs")

	rootCmd.AddCommand(logReaderCmd)
}

var logReaderCmd = &cobra.Command{
	Use:   "log-reader",
	Short: "Log reader",
	Long:  `Log reader`,
	Run: func(cmd *cobra.Command, args []string) {
		if serverLogsFlag == dbLogsFlag {
			fmt.Println(cmd.UsageString())
			os.Exit(1)
		}

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			if line, err := parseLine(scanner.Bytes(), serverLogsFlag, dbLogsFlag); err == nil {
				fmt.Println(line)
			}
		}

		if err := scanner.Err(); err != nil {
			log.Println(err)
		}
	},
}

type errorLog struct {
	Level     string    `json:"level"`
	Time      time.Time `json:"time"`
	Caller    string    `json:"caller"`
	Message   string    `json:"message"`
	Error     string    `json:"error"`
	Code      uint      `json:"code"`
	Method    string    `json:"method"`
	Path      string    `json:"path"`
	Body      string    `json:"body"`
	URL       string    `json:"url"`
	Host      string    `json:"host"`
	IP        string    `json:"ip"`
	RequestID string    `json:"requestId"`
}

func parseLine(line []byte, serverLogs, dbLogs bool) (string, error) {
	if serverLogs {
		return parseLineServer(line)
	} else if dbLogs {
		return string(line), nil
	}
	return "", errors.New("invalid flag")
}

func parseLineServer(line []byte) (string, error) {
	var errLog errorLog
	err := json.Unmarshal(line, &errLog)
	if err != nil {
		return "", err
	}

	code := ""
	if errLog.Code != 0 {
		code = fmt.Sprintf(" | Code: %d", errLog.Code)
	}
	message := ""
	if errLog.Message != "" {
		message = fmt.Sprintf(" | Message: %s", errLog.Message)
	}
	errorLog := ""
	if errLog.Error != "" && errLog.Error != "<nil>" {
		errorLog = fmt.Sprintf(" | Error: %s", errLog.Error)
	}
	method := ""
	if errLog.Method != "" {
		method = fmt.Sprintf(" | Method: %s", errLog.Method)
	}
	url := ""
	if errLog.URL != "" {
		url = fmt.Sprintf(" | URL: %s", errLog.URL)
	}
	host := ""
	if errLog.Host != "" {
		host = fmt.Sprintf(" | Host: %s", errLog.Host)
	}
	ip := ""
	if errLog.IP != "" {
		ip = fmt.Sprintf(" | IP: %s", errLog.IP)
	}
	requestId := ""
	if errLog.RequestID != "" {
		requestId = fmt.Sprintf(" | RequestID: %s", errLog.RequestID)
	}

	result := fmt.Sprintf("%s %7s | %s%s%s%s%s%s%s%s%s",
		errLog.Time.Format(time.RFC3339),
		displayLevel(errLog.Level),
		errLog.Caller,
		code,
		message,
		errorLog,
		method,
		host,
		url,
		ip,
		requestId,
	)

	return result, nil
}
