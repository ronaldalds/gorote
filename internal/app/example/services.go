package example

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"slices"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/ronaldalds/base-go-api/internal/config/access"
	"github.com/ronaldalds/base-go-api/internal/config/databases"
	"github.com/ronaldalds/base-go-api/internal/config/envs"
	"gorm.io/gorm"
)

type Service struct {
	GormStore  *databases.GormStore
	RedisStore *databases.RedisStore
	MongoStore *databases.MongoStore
}

func NewService() *Service {
	service := &Service{
		GormStore:  databases.DB.GormStore,
		RedisStore: databases.DB.RedisStore,
		MongoStore: databases.DB.MongoStore,
	}
	// Executar as Migrations
	service.GormStore.DB.AutoMigrate(&Example{})
	// Executar as Seeds
	return service
}