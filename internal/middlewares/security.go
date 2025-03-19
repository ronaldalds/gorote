package middlewares

import (
	"github.com/gofiber/fiber/v2/middleware/helmet"
)

func (m *Middleware) SecurityMiddleware() {
	m.App.Use(helmet.New(helmet.Config{
		XSSProtection: "1; mode=block",
	}))
}
