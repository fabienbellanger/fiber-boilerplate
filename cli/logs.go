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
	logReaderCmd.Flags().BoolVarP(&serverLogsFlag, "server", "s", false, "server access/error logs")
	logReaderCmd.Flags().BoolVarP(&dbLogsFlag, "db", "d", false, "database logs")

	rootCmd.AddCommand(logReaderCmd)
}

var logReaderCmd = &cobra.Command{
	Use:   "logs",
	Short: "Logs reader",
	Long:  `Logs reader`,
	Run: func(cmd *cobra.Command, args []string) {
		if serverLogsFlag == dbLogsFlag {
			fmt.Println(cmd.UsageString())
			os.Exit(1)
		}

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line, _ := parseLine(scanner.Bytes(), serverLogsFlag, dbLogsFlag)
			fmt.Println(line)
		}

		if err := scanner.Err(); err != nil {
			log.Println(err)
		}
	},
}

type errorLog struct {
	Level       string    `json:"level"`
	Time        time.Time `json:"time"`
	Caller      string    `json:"caller"`
	Message     string    `json:"message"`
	Description string    `json:"description"`
	Error       string    `json:"error"`
	Code        uint      `json:"code"`
	Method      string    `json:"method"`
	Path        string    `json:"path"`
	Body        string    `json:"body"`
	URL         string    `json:"url"`
	Host        string    `json:"host"`
	IP          string    `json:"ip"`
	RequestID   string    `json:"requestId"`
	Latency     string    `json:"latency"`
	UserAgent   string    `json:"userAgent"`
}

func parseLine(line []byte, serverLogs, dbLogs bool) (string, error) {
	if serverLogs {
		return parseLineServer(line)
	} else if dbLogs {
		return string(line), nil
	}
	return "", errors.New("invalid flag")
}

// TODO: Improve parser
func parseLineServer(line []byte) (string, error) {
	var errLog errorLog
	err := json.Unmarshal(line, &errLog)
	if err != nil {
		return string(line), err
	}

	code := ""
	if errLog.Code != 0 {
		code = fmt.Sprintf(" | %d", displayLogStatusCode(errLog.Code))
	}
	message := ""
	if errLog.Message != "" {
		message = fmt.Sprintf(" | Message: %s", errLog.Message)
	}
	description := ""
	if errLog.Description != "" {
		description = fmt.Sprintf(" | Description: %s", errLog.Description)
	}
	errorLog := ""
	if errLog.Error != "" && errLog.Error != "<nil>" {
		errorLog = fmt.Sprintf(" | Error: %s", errLog.Error)
	}
	method := ""
	if errLog.Method != "" {
		method = fmt.Sprintf(" | %6s", displayLogMethod(errLog.Method))
	}
	url := ""
	if errLog.URL != "" {
		url = fmt.Sprintf(" | %s", errLog.URL)
	}
	path := ""
	if errLog.Path != "" {
		path = fmt.Sprintf(" | %s", errLog.Path)
	}
	host := ""
	if errLog.Host != "" {
		host = fmt.Sprintf(" | %s", errLog.Host)
	}
	ip := ""
	if errLog.IP != "" {
		ip = fmt.Sprintf(" | IP: %s", errLog.IP)
	}
	requestID := ""
	if errLog.RequestID != "" {
		requestID = fmt.Sprintf(" | RequestID: %s", errLog.RequestID)
	}
	userAgent := ""
	if errLog.UserAgent != "" {
		userAgent = fmt.Sprintf(" | UserAgent: %s", errLog.UserAgent)
	}
	latency := ""
	if errLog.Latency != "" {
		latency = fmt.Sprintf(" | %s", errLog.Latency)
	}

	result := fmt.Sprintf("%s | %7s %s%s%s%s%s%s%s%s%s%s%s%s%s",
		errLog.Time.Format(time.RFC3339),
		displayLogLevel(errLog.Level),
		code,
		method,
		message,
		description,
		errorLog,
		path,
		url,
		host,
		ip,
		requestID,
		userAgent,
		latency,
		" | "+errLog.Caller,
	)

	return result, nil
}
