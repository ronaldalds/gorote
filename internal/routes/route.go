package routes

import (
	"github.com/ronaldalds/base-go-api/internal/app/core"
	"github.com/ronaldalds/base-go-api/internal/app/papai"
	"github.com/ronaldalds/base-go-api/internal/app/teletubbies"
	"github.com/ronaldalds/base-go-api/internal/middlewares"

	"github.com/gofiber/fiber/v2"
)

type Router struct {
	App         *fiber.App
	Middleware  *middlewares.Middleware
	Core        *core.Router
}

func NewRouter(app *fiber.App) *Router {
	return &Router{
		App:         app,
		Middleware:  middlewares.NewMiddleware(app),
		Core:        core.NewRouter(app),
	}
}

func (r *Router) RegisterFiberRoutes() {
	r.Middleware.CorsMiddleware()
	r.Middleware.SecurityMiddleware()
	r.Middleware.Telemetry("auth/login")

	apiV2 := r.App.Group("/api/v2")
	r.Core.RegisterRouter(apiV2)
}
