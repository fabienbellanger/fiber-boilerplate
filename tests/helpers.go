package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"testing"
	"time"

	server "github.com/fabienbellanger/fiber-boilerplate"
	"github.com/fabienbellanger/fiber-boilerplate/db"
	"github.com/fabienbellanger/fiber-boilerplate/entities"
	storeUser "github.com/fabienbellanger/fiber-boilerplate/stores/user"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const (
	UserUsername = "test@test.com"
	UserPassword = "00000000"
)

// Test defines a structure for specifying input and output data of a single test case.
type Test struct {
	Description string

	// Test input
	Route   string
	Method  string
	Body    io.Reader
	Headers []Header

	// Check
	CheckError bool
	CheckBody  bool
	CheckCode  bool

	// Expected output
	ExpectedError bool
	ExpectedCode  int
	ExpectedBody  string
}

// Header represents an header value.
type Header struct {
	Key   string
	Value string
}

// Init initializes configuration from .env path and returns TestDB.
func Init(p string) TestDB {
	viper.SetConfigFile(p)
	viper.ReadInConfig()

	viper.Set("SERVER_MONITOR", false)
	viper.Set("SERVER_PROMETHEUS", false)
	viper.Set("SERVER_PPROF", false)
	viper.Set("GORM_LOG_OUTPUT", "stdout")
	viper.Set("LIMITER_ENABLE", false)

	tdb, err := newTestDB()
	if err != nil {
		log.Panicf("%v\n", err)
	}
	return tdb
}

// TestDB is used to create and use a database for tests.
type TestDB struct {
	name  string
	DB    *db.DB
	Token string
}

// newTestDB returns a TestDB instance.
func newTestDB() (TestDB, error) {
	rand.Seed(time.Now().UnixNano())
	dbName := viper.GetString("DB_DATABASE") + "__" + fmt.Sprintf("%08d", rand.Int63n(1e8))

	config := db.DatabaseConfig{
		Driver:          viper.GetString("DB_DRIVER"),
		Host:            viper.GetString("DB_HOST"),
		Username:        viper.GetString("DB_USERNAME"),
		Password:        viper.GetString("DB_PASSWORD"),
		Port:            viper.GetInt("DB_PORT"),
		Database:        "",
		Charset:         viper.GetString("DB_CHARSET"),
		Collation:       viper.GetString("DB_COLLATION"),
		Location:        viper.GetString("DB_LOCATION"),
		MaxIdleConns:    viper.GetInt("DB_MAX_IDLE_CONNS"),
		MaxOpenConns:    viper.GetInt("DB_MAX_OPEN_CONNS"),
		ConnMaxLifetime: viper.GetDuration("DB_CONN_MAX_LIFETIME") * time.Hour,
	}

	db, err := db.New(&config)
	if err != nil {
		return TestDB{}, err
	}

	// Create databse for test, use it and run migrations
	db.Exec("CREATE DATABASE IF NOT EXISTS `" + dbName + "`;")
	db.Exec("USE `" + dbName + "`;")
	db.MakeMigrations()

	// Create first user and get token
	token, err := createUserAndAuthenticate(db)
	if err != nil {
		return TestDB{}, err
	}

	return TestDB{DB: db, name: dbName, Token: token}, nil
}

// Drop database after the test.
func (tdb *TestDB) Drop() error {
	result := tdb.DB.Exec("DROP DATABASE IF EXISTS `" + tdb.name + "`;")

	return result.Error
}

// Create a first user, authenticate him and return JWT.
func createUserAndAuthenticate(db *db.DB) (token string, err error) {
	// Create first user
	userStore := storeUser.New(db)
	err = userStore.Create(&entities.User{
		Lastname:  "User",
		Firstname: "Test",
		Password:  UserPassword,
		Username:  UserUsername,
	})
	if err != nil {
		return
	}

	// Get User
	user, err := userStore.Login(UserUsername, UserPassword)
	if err != nil {
		return
	}

	// Get token
	token, _, err = user.GenerateJWT(viper.GetDuration("JWT_LIFETIME"), viper.GetString("JWT_ALGO"), viper.GetString("JWT_SECRET"))
	if err != nil {
		return
	}

	return token, err
}

// Execute runs all tests.
func Execute(t *testing.T, db *db.DB, tests []Test) {
	// Setup the app as it is done in the main function
	app := server.Setup(db, nil)

	// Iterate through test single test cases
	for _, test := range tests {
		// Create a new http request with the route from the test case
		req, _ := http.NewRequest(test.Method, test.Route, test.Body)
		for _, h := range test.Headers {
			req.Header.Add(h.Key, h.Value)
		}

		// Perform the request plain with the app.
		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		if test.CheckError {
			// Verify that no error occured, that is not expected
			assert.Equalf(t, test.ExpectedError, err != nil, test.Description)

			// As expected errors lead to broken responses, the next test case needs to be processed
			if test.ExpectedError {
				continue
			}
		}

		// Verify if the status code is as expected
		if test.CheckCode {
			assert.Equalf(t, test.ExpectedCode, res.StatusCode, test.Description)
		}

		// Verify if the body is as expected
		if test.CheckBody {
			// Read the response body
			body, err := io.ReadAll(res.Body)

			// Reading the response body should work everytime, such that
			// the err variable should be nil
			assert.Nilf(t, err, test.Description)

			// Verify, that the reponse body equals the expected body
			assert.Equalf(t, test.ExpectedBody, string(body), test.Description)
		}
	}
}

// JsonToString converts a JSON to a string.
func JsonToString(d interface{}) string {
	b, err := json.Marshal(d)
	if err != nil {
		log.Panicf("%v\n", err)
	}
	return string(b)
}
