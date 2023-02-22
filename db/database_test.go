package db

import (
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm/logger"
)

func TestGetGormLogLevel(t *testing.T) {
	assert.Equal(t, logger.Silent, getGormLogLevel("silent", "development"))
	assert.Equal(t, logger.Info, getGormLogLevel("info", "development"))
	assert.Equal(t, logger.Warn, getGormLogLevel("warn", "development"))
	assert.Equal(t, logger.Error, getGormLogLevel("error", "development"))
	assert.Equal(t, logger.Warn, getGormLogLevel("", "development"))
	assert.Equal(t, logger.Error, getGormLogLevel("", "production"))
}

func TestGetGormLogOutput(t *testing.T) {
	output, err := getGormLogOutput("stdout", "", "production")
	assert.Equal(t, os.Stdout, output)
	assert.Nil(t, err)

	output, err = getGormLogOutput("stdout", "", "development")
	assert.Equal(t, os.Stdout, output)
	assert.Nil(t, err)

	output, err = getGormLogOutput("file", "test.log", "production")
	f, _ := os.OpenFile(path.Clean("test.log"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	defer f.Close()
	defer os.Remove("test.log")

	assert.IsType(t, f, output)
	assert.Nil(t, err)
}

func TestDsn(t *testing.T) {
	c := DatabaseConfig{
		Driver:          "mysql",
		Host:            "localhost",
		Username:        "root",
		Password:        "root",
		Port:            3306,
		Database:        "fiber",
		Charset:         "utf8mb4",
		Collation:       "utf8mb4_general_ci",
		Location:        "UTC",
		MaxIdleConns:    10,
		MaxOpenConns:    10,
		ConnMaxLifetime: time.Second,
		SlowThreshold:   200 * time.Millisecond,
	}
	expected := "root:root@tcp(localhost:3306)/fiber?parseTime=True&charset=utf8mb4&collation=utf8mb4_general_ci&loc=UTC"
	wanted, err := c.dsn()
	assert.Equal(t, expected, wanted)
	assert.Nil(t, err)

	c.Host = ""
	_, err = c.dsn()
	assert.NotNil(t, err)
	assert.EqualError(t, err, "error in database configuration")
}

func TestPaginateValues(t *testing.T) {
	type args struct {
		page  string
		limit string
	}

	type result struct {
		offset int
		limit  int
	}

	tests := []struct {
		name   string
		args   args
		wanted result
	}{
		{
			name:   "First page",
			args:   args{"1", "100"},
			wanted: result{offset: 0, limit: 100},
		},
		{
			name:   "Third page",
			args:   args{"3", "100"},
			wanted: result{offset: 200, limit: 100},
		},
		{
			name:   "Invalid page",
			args:   args{"-3", "100"},
			wanted: result{offset: 0, limit: 100},
		},
		{
			name:   "Too large limit",
			args:   args{"1", "1000"},
			wanted: result{offset: 0, limit: MaxLimit},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			offset, limit := paginateValues(tt.args.page, tt.args.limit)
			got := result{offset: offset, limit: limit}
			assert.Equal(t, got, tt.wanted)
		})
	}
}

func TestOrderValues(t *testing.T) {
	type args struct {
		list     string
		prefixes []string
	}

	tests := []struct {
		name   string
		args   args
		wanted map[string]string
	}{
		{
			name: "One field with no prefix",
			args: args{"+created_at", []string{}},
			wanted: map[string]string{
				"created_at": "ASC",
			},
		},
		{
			name: "One field with prefix",
			args: args{"+created_at", []string{"table"}},
			wanted: map[string]string{
				"table.created_at": "ASC",
			},
		},
		{
			name: "Many fields with no prefix",
			args: args{"+created_at,-id", []string{}},
			wanted: map[string]string{
				"created_at": "ASC",
				"id":         "DESC",
			},
		},
		{
			name: "Many fields with prefix",
			args: args{"+created_at,-id", []string{"table"}},
			wanted: map[string]string{
				"table.created_at": "ASC",
				"table.id":         "DESC",
			},
		},
		{
			name:   "No fields",
			args:   args{"", []string{"toto"}},
			wanted: map[string]string{},
		},
		{
			name:   "With invalid field",
			args:   args{"created_at", []string{}},
			wanted: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, orderValues(tt.args.list, tt.args.prefixes...), tt.wanted)
		})
	}
}
