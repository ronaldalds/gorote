package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ronaldalds/base-go-api/internal/config/databases"
)

type Middleware struct {
	App        *fiber.App
	RedisStore *databases.RedisStore
}

func NewMiddleware(app *fiber.App) *Middleware {
	return &Middleware{
		App:        app,
		RedisStore: databases.DB.RedisStore,
	}
}
