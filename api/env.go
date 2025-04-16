package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type NoSQL struct {
	Username string
	Password string
	Host     string
	Port     int
	Database string
}

type Redis struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type SQL struct {
	Username string
	Password string
	Host     string
	Port     int
	Database string
	Schema   string
}

type Super struct {
	Name     string
	Username string
	Password string
	Email    string
	Phone    string
}

type JWT struct {
	Secret        string
	ExpireAccess  time.Duration
	ExpireRefresh time.Duration
}

type Log struct {
	Url  string
	Port int
}

type APP struct {
	Name     string
	TimeZone string
	Port     int
}

type Envs struct {
	Sql   SQL
	NoSql *NoSQL
	Redis *Redis
	Super Super
	Jwt   JWT
	Logs  *Log
	App   APP
}

var Env Envs

// Load reads and validates environment variables
func Load() {
	app := APP{
		Name:     getEnv("APP_NAME", true),
		TimeZone: getEnv("APP_TIMEZONE", false, "America/Sao_Paulo"),
		Port:     getEnvAsInt("APP_PORT", true),
	}
	sql := SQL{
		Username: getEnv("SQL_USERNAME", true),
		Password: getEnv("SQL_PASSWORD", true),
		Host:     getEnv("SQL_HOST", false, "localhost"),
		Port:     getEnvAsInt("SQL_PORT", true),
		Database: getEnv("SQL_DATABASE", true),
		Schema:   getEnv("SQL_SCHEMA", true),
	}
	jwt := JWT{
		Secret:        getEnv("JWT_SECRET", true),
		ExpireAccess:  getEnvAsTime("JWT_EXPIRE_ACCESS", false, 5),
		ExpireRefresh: getEnvAsTime("JWT_EXPIRE_REFRESH", false, 10080),
	}
	super := Super{
		Name:     getEnv("SUPER_NAME", false, "Admin"),
		Username: getEnv("SUPER_USERNAME", false, "admin"),
		Password: getEnv("SUPER_PASS", false, "admin"),
		Email:    getEnv("SUPER_EMAIL", false, "ronald.ralds@gmail.com"),
		Phone:    getEnv("SUPER_PHONE", false, "+558892200365"),
	}

	noSql := &NoSQL{
		Username: getEnv("NOSQL_USERNAME", false),
		Password: getEnv("NOSQL_PASSWORD", false),
		Host:     getEnv("NOSQL_HOST", false, "localhost"),
		Port:     getEnvAsInt("NOSQL_PORT", false),
		Database: getEnv("NOSQL_DATABASE", false),
	}
	if noSql.Username == "" || noSql.Password == "" || noSql.Host == "" || noSql.Port == 0 || noSql.Database == "" {
		fmt.Println("NoSQL disabled")
		noSql = nil
	}

	redis := &Redis{
		Host:     getEnv("REDIS_HOST", false, "localhost"),
		Port:     getEnvAsInt("REDIS_PORT", false),
		Password: getEnv("REDIS_PASSWORD", false),
		DB:       getEnvAsInt("REDIS_DB", false),
	}
	if redis.Host == "" || redis.Port == 0 || redis.Password == "" || redis.DB == 0 {
		fmt.Println("Redis disabled")
		redis = nil
	}

	log := &Log{
		Url:  getEnv("LOG_URL", false),
		Port: getEnvAsInt("LOG_PORT", false),
	}
	if log.Url == "" || log.Port == 0 {
		fmt.Println("Logs disabled")
		log = nil
	}

	Env = Envs{
		Sql:   sql,
		NoSql: noSql,
		Redis: redis,
		Super: super,
		Jwt:   jwt,
		Logs:  log,
		App:   app,
	}
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
