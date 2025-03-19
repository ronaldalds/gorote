package middlewares

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func (m *Middleware) CorsMiddleware() {
	m.App.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: false,
		MaxAge:           300,
	}))
}
