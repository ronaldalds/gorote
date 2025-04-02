package routes

import (
	"github.com/ronaldalds/gorote-core/core"
	"github.com/ronaldalds/base-go-api/internal/middlewares"

	"github.com/gofiber/fiber/v2"
)

type Router struct {
	App         *fiber.App
	MiddlewareLocal  *middlewares.Middleware
	Core        *core.Router
}

func NewRouter(app *fiber.App) *Router {
	return &Router{
		App:         app,
		MiddlewareLocal:  middlewares.NewMiddleware(app),
		Core:        core.New(app),
	}
}

func (r *Router) RegisterFiberRoutes() {
	r.MiddlewareLocal.CorsMiddleware()
	r.MiddlewareLocal.SecurityMiddleware()
	if envs.Env.LogsUrl != "" {
		r.MiddlewareLocal.Telemetry("auth/login")
	}

	apiV2 := r.App.Group("/api/v2")
	r.Core.RegisterRouter(apiV2)
}
