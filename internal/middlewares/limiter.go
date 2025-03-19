package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func (m *Middleware) Limited(max int) func(c *fiber.Ctx) error {
	config := limiter.Config{
		Max: max,
		LimitReached: func(c *fiber.Ctx) error {
			return fiber.ErrTooManyRequests
		},
	}
	return limiter.New(config)
}
