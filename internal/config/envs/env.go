package envs

import (
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type Envs struct {
	// SQL
	SqlUsername string
	SqlPassword string
	SqlHost     string
	SqlPort     int
	SqlDatabase string
	SqlSchema   string
	// NoSQL
	NoSqlUsername string
	NoSqlPassword string
	NoSqlHost     string
	NoSqlPort     int
	NoSqlDatabase string
	// Redis
	RedisDb       int
	RedisHost     string
	RedisPort     int
	RedisPassword string
	// JWT
	JwtSecret        string
	JwtExpireAcess   time.Duration
	JwtExpireRefresh time.Duration
	// SUPER USER
	SuperName     string
	SuperUsername string
	SuperPass     string
	SuperEmail    string
	SuperPhone    string
	// APP
	LogsPort           int
	LogsUrl            string
	AppName            string
	TimeUCT            time.Location
	TimeZone           string
	Port               int
}

var Env Envs

// Load reads and validates environment variables
func Load() {
	Env = Envs{
		// SQL
		SqlUsername: getEnv("SQL_USERNAME", true),
		SqlPassword: getEnv("SQL_PASSWORD", true),
		SqlHost:     getEnv("SQL_HOST", false, "localhost"),
		SqlPort:     getEnvAsInt("SQL_PORT", true),
		SqlDatabase: getEnv("SQL_DATABASE", true),
		SqlSchema:   getEnv("SQL_SCHEMA", true),
		// NOSQL
		NoSqlUsername: getEnv("NOSQL_USERNAME", false),
		NoSqlPassword: getEnv("NOSQL_PASSWORD", false),
		NoSqlHost:     getEnv("NOSQL_HOST", false, "localhost"),
		NoSqlPort:     getEnvAsInt("NOSQL_PORT", false),
		NoSqlDatabase: getEnv("NOSQL_DATABASE", false),
		// Redis
		RedisDb:       getEnvAsInt("REDIS_DB", false),
		RedisHost:     getEnv("REDIS_HOST", false, "localhost"),
		RedisPort:     getEnvAsInt("REDIS_PORT", false),
		RedisPassword: getEnv("REDIS_PASSWORD", false),
		// JWT
		JwtSecret:        getEnv("JWT_SECRET", true),
		JwtExpireAcess:   getEnvAsTime("JWT_EXPIRE_ACCESS", false, 5),
		JwtExpireRefresh: getEnvAsTime("JWT_EXPIRE_REFRESH", false, 10080),
		// SUPER USER
		SuperName:     getEnv("SUPER_NAME", false, "Admin"),
		SuperUsername: getEnv("SUPER_USERNAME", false, "admin"),
		SuperPass:     getEnv("SUPER_PASS", false, "admin"),
		SuperEmail:    getEnv("SUPER_EMAIL", false, "ronald.ralds@gmail.com"),
		SuperPhone:    getEnv("SUPER_PHONE", false, "+558892200365"),
		// APP
		LogsPort:           getEnvAsInt("LOG_PORT", true),
		LogsUrl:            getEnv("LOG_URL", false, "http://localhost"),
		AppName:            getEnv("APP_NAME", false, "app"),
		TimeUCT:            getUCT("TIMEZONE", false, "America/Fortaleza"),
		TimeZone:           getEnv("TIMEZONE", false, "America/Fortaleza"),
		Port:               getEnvAsInt("PORT", false, 3000),
	}
}

func getUCT(key string, required bool, defaultValue ...string) time.Location {
	value := os.Getenv(key)

	if value == "" {
		if required {
			panic(fmt.Sprintf("variable %s is required", key))
		}
		if len(defaultValue) > 0 {
			location, err := time.LoadLocation(value)
			if err != nil {
				panic(fmt.Sprintf("invalid timezone: %s", err.Error()))
			}
			return *location
		}
	}
	location, err := time.LoadLocation(value)
	if err != nil {
		panic(fmt.Sprintf("invalid timezone: %s", err.Error()))
	}
	return *location
}

func getEnv(key string, required bool, defaultValue ...string) string {
	value := os.Getenv(key)

	if value == "" {
		if required {
			panic(fmt.Sprintf("variable %s is required", key))
		}
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return ""
	}
	return value
}

func getEnvAsInt(key string, required bool, defaultValue ...int) int {
	valueStr := os.Getenv(key)

	if valueStr == "" {
		if required {
			panic(fmt.Sprintf("variable %s is required", key))
		}
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		panic(fmt.Sprintf("failed to convert %s to integer: %v", key, err))
	}
	return value
}

func getEnvAsTime(key string, required bool, defaultValue ...int) time.Duration {
	valueStr := os.Getenv(key)

	if valueStr == "" {
		if required {
			panic(fmt.Sprintf("variable %s is required", key))
		}
		if len(defaultValue) > 0 {
			return time.Duration(defaultValue[0]) * time.Minute
		}
		return 0
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		panic(fmt.Sprintf("failed to convert %s to integer: %v", key, err))
	}
	return time.Duration(value) * time.Minute
}
