package main

import (
	"fmt"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/ronaldalds/gorote-core-rsa/core"
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
	App   APP
	Sql   SQL
	NoSql *NoSQL
	Redis *Redis
	Super Super
	Jwt   JWT
	Logs  *Log
}

var Env Envs

// Load reads and validates environment variables
func Load() {
	app := APP{
		Name:     core.GetEnv("APP_NAME", true),
		TimeZone: core.GetEnv("APP_TIMEZONE", false, "America/Sao_Paulo"),
		Port:     core.GetEnvAsInt("APP_PORT", true),
	}
	sql := SQL{
		Username: core.GetEnv("SQL_USERNAME", true),
		Password: core.GetEnv("SQL_PASSWORD", true),
		Host:     core.GetEnv("SQL_HOST", true),
		Port:     core.GetEnvAsInt("SQL_PORT", true),
		Database: core.GetEnv("SQL_DATABASE", true),
		Schema:   core.GetEnv("SQL_SCHEMA", true),
	}
	jwt := JWT{
		ExpireAccess:  core.GetEnvAsTime("JWT_EXPIRE_ACCESS", false, 5),
		ExpireRefresh: core.GetEnvAsTime("JWT_EXPIRE_REFRESH", false, 10080),
	}
	super := Super{
		Name:     core.GetEnv("SUPER_NAME", false, "Admin"),
		Username: core.GetEnv("SUPER_USERNAME", false, "admin"),
		Password: core.GetEnv("SUPER_PASS", false, "admin"),
		Email:    core.GetEnv("SUPER_EMAIL", false, "ronald.ralds@gmail.com"),
		Phone:    core.GetEnv("SUPER_PHONE", false, "+558892200365"),
	}

	noSql := &NoSQL{
		Username: core.GetEnv("NOSQL_USERNAME", false),
		Password: core.GetEnv("NOSQL_PASSWORD", false),
		Host:     core.GetEnv("NOSQL_HOST", true),
		Port:     core.GetEnvAsInt("NOSQL_PORT", false),
		Database: core.GetEnv("NOSQL_DATABASE", false),
	}
	if noSql.Username == "" || noSql.Password == "" || noSql.Host == "" || noSql.Port == 0 || noSql.Database == "" {
		fmt.Println("NoSQL disabled")
		noSql = nil
	}

	redis := &Redis{
		Host:     core.GetEnv("REDIS_HOST", false, "localhost"),
		Port:     core.GetEnvAsInt("REDIS_PORT", false),
		Password: core.GetEnv("REDIS_PASSWORD", false),
		DB:       core.GetEnvAsInt("REDIS_DB", false),
	}
	if redis.Host == "" || redis.Port == 0 || redis.Password == "" || redis.DB == 0 {
		fmt.Println("Redis disabled")
		redis = nil
	}

	log := &Log{
		Url:  core.GetEnv("LOG_URL", false),
		Port: core.GetEnvAsInt("LOG_PORT", false),
	}
	if log.Url == "" || log.Port == 0 {
		fmt.Println("Logs disabled")
		log = nil
	}

	Env = Envs{
		App:   app,
		Sql:   sql,
		NoSql: noSql,
		Redis: redis,
		Super: super,
		Jwt:   jwt,
		Logs:  log,
	}
}
