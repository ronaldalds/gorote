package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
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
		Username: Env.NoSqlUsername,
		Password: Env.NoSqlPassword,
		Host:     Env.NoSqlHost,
		Port:     Env.NoSqlPort,
		Database: Env.NoSqlDatabase,
	}

	// Configuração final do banco de dados
	DB.MongoStore = dbMongo.newMongoStore()
}

func LoadRedis() {
	// Inicializando Redis
	dbRedis := &InitRedis{
		Host:     Env.RedisHost,
		Port:     Env.RedisPort,
		Password: Env.RedisPassword,
		DB:       Env.RedisDb,
	}
	DB.RedisStore = dbRedis.newRedisStore()
}

func LoadSQL(extraModels ...any) {
	// Inicializando Gorm
	dbGorm := &InitGorm{
		Host:     Env.SqlHost,
		User:     Env.SqlUsername,
		Password: Env.SqlPassword,
		Database: Env.SqlDatabase,
		Port:     Env.SqlPort,
		TimeZone: Env.TimeZone,
		Schema:   Env.SqlSchema,
		Models:   extraModels,
	}

	// Configuração final do banco de dados
	DB.GormStore = dbGorm.newGormStore()
}

func (m *InitMongo) newMongoStore() *MongoStore {
	// Cria um contexto com timeout para a conexão
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=admin&ssl=false",
		m.Username,
		m.Password,
		m.Host,
		m.Port,
		m.Database,
	)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Conecta ao MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("failed to connect to MongoDB: %w", err)
	}
	// Verifica a conexão
	if err = client.Ping(ctx, nil); err != nil {
		log.Fatal("failed to ping MongoDB: %w", err)
	}
	fmt.Println("Connected to MongoDB")
	// Configura o banco de dados
	db := client.Database(m.Database)
	return &MongoStore{
		Client:   client,
		Database: db,
	}
}

func (r *InitRedis) newRedisStore() *RedisStore {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Configura o cliente Redis
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", r.Host, r.Port),
		Password: r.Password, // "" para Redis sem autenticação
		DB:       r.DB,       // Número do banco
	})

	// Verifica a conexão com o Redis
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	fmt.Println("Connected to Redis")
	return &RedisStore{client}
}

func (g *InitGorm) newGormStore() *GormStore {
	// Cria um contexto com timeout de 10 segundos
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Cria o Data Source Name (DSN)
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=%s search_path=%s",
		g.Host,
		g.Port,
		g.User,
		g.Password,
		g.Database,
		g.TimeZone,
		g.Schema,
	)

	// Conecta ao banco de dados usando GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to the database:", err)
	}

	// Obtém o banco de dados subjacente (*sql.DB)
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("failed to get *sql.DB from GORM:", err)
	}

	// Tenta fazer um ping para garantir que a conexão está ativa
	err = sqlDB.PingContext(ctx)
	if err != nil {
		log.Fatal("failed to ping database:", err)
	}

	// Cria o schema, caso não exista
	_, err = sqlDB.ExecContext(ctx, fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s;", g.Schema))
	if err != nil {
		log.Fatal("failed to create schema:", err)
	}

	if len(g.Models) > 0 {
		// Faz a migração do modelo
		err_db := db.AutoMigrate(
			g.Models...,
		)
		if err_db != nil {
			log.Fatal("failed to auto migrate:", err_db)
		}
	}

	fmt.Println("Connected to Gorm and schema is set up.")

	// Retorna o GormStore
	return &GormStore{db}

}
