package core

import (
	"github.com/gofiber/fiber/v2"
)

func ImplementMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return ctx.Next()
	}
}
