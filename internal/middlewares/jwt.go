package middlewares

import (
	"context"
	"fmt"
	"log"
	"slices"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/ronaldalds/base-go-api/internal/config/envs"
	"github.com/ronaldalds/base-go-api/internal/utils"
)

func (m *Middleware) GetKey(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	val, err := m.RedisStore.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key does not exist")
	} else if err != nil {
		return "", fmt.Errorf("failed to get key: %v", err)
	}
	return val, nil
}

func (m *Middleware) JWTProtected(permissions ...string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		token, err := utils.GetJwtHeaderPayload(ctx.Get("Authorization"), envs.Env.JwtSecret)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}

		// Check Redis for active session
		session, err := m.GetKey(fmt.Sprintf("%d", token.Claims.Sub))
		if err != nil {
			log.Println("Redis error:", err) // Log the Redis error
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}

		if session != token.Token {
			log.Println("Token does not match active session") // Log the error
			return fiber.NewError(fiber.StatusUnauthorized, "token does not match active session")
		}

		// Check permissions
		if token.Claims.IsSuperUser {
			return ctx.Next()
		}
		if len(permissions) == 0 {
			return ctx.Next()
		}

		// Check if any required permission exists in user's permissions
		for _, requiredPermission := range permissions {
			if slices.Contains(token.Claims.Permissions, requiredPermission) {
				log.Println("Permission validated, proceeding to next handler")
				return ctx.Next()
			}
		}

		// If no errors, log success and continue to the next handler
		log.Println("JWT validated and session matched, proceeding to next handler")
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}
}
