package databases

import (
	"github.com/go-redis/redis/v8"
	"github.com/ronaldalds/base-go-api/internal/config/envs"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var DB Database

type InitRedis struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type InitGorm struct {
	Host     string
	User     string
	Password string
	Database string
	Port     int
	TimeZone string
	Schema   string
	Models   []any
}

type InitMongo struct {
	Username string
	Password string
	Host     string
	Port     int
	Database string
}

type RedisStore struct {
	Client *redis.Client
}

type GormStore struct {
	*gorm.DB
}

type MongoStore struct {
	Client   *mongo.Client
	Database *mongo.Database
}

type Database struct {
	GormStore  *GormStore
	RedisStore *RedisStore
	MongoStore *MongoStore
}

func LoadNOSQL() {
	// Inicializando Mongo
	dbMongo := &InitMongo{
		Username: envs.Env.NoSqlUsername,
		Password: envs.Env.NoSqlPassword,
		Host:     envs.Env.NoSqlHost,
		Port:     envs.Env.NoSqlPort,
		Database: envs.Env.NoSqlDatabase,
	}

	// Configuração final do banco de dados
	DB.MongoStore = dbMongo.newMongoStore()
}

func LoadRedis() {
	// Inicializando Redis
	dbRedis := &InitRedis{
		Host:     envs.Env.RedisHost,
		Port:     envs.Env.RedisPort,
		Password: envs.Env.RedisPassword,
		DB:       envs.Env.RedisDb,
	}
	DB.RedisStore = dbRedis.newRedisStore()
}

func LoadSQL(extraModels ...any) {
	// Inicializando Gorm
	dbGorm := &InitGorm{
		Host:     envs.Env.SqlHost,
		User:     envs.Env.SqlUsername,
		Password: envs.Env.SqlPassword,
		Database: envs.Env.SqlDatabase,
		Port:     envs.Env.SqlPort,
		TimeZone: envs.Env.TimeZone,
		Schema:   envs.Env.SqlSchema,
		Models:   append([]any{}, extraModels...), // Concatena os modelos recebidos como parâmetro
	}

	// Configuração final do banco de dados
	DB.GormStore = dbGorm.newGormStore()
}
