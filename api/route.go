package main

import (
	"github.com/ronaldalds/gorote-core/core"

	"github.com/gofiber/fiber/v2"
)

type Router struct {
	App             *fiber.App
	MiddlewareLocal *Middleware
	Core            *core.Router
}

func New(app *fiber.App) *Router {
	coreConfig := &core.AppConfig{
		App:       app,
		GormStore: DB.GormStore.DB,
		Jwt: core.AppJwt{
			AppName:          Env.AppName,
			TimeZone:         Env.TimeZone,
			JwtSecret:        Env.JwtSecret,
			JwtExpireAccess:  Env.JwtExpireAccess,
			JwtExpireRefresh: Env.JwtExpireRefresh,
		},
		Super: &core.AppSuper{
			SuperName:  Env.SuperName,
			SuperUser:  Env.SuperUsername,
			SuperPass:  Env.SuperPass,
			SuperEmail: Env.SuperEmail,
			SuperPhone: Env.SuperPhone,
		},
	}

	return &Router{
		App:             app,
		MiddlewareLocal: NewMiddleware(app),
		Core:            core.New(coreConfig),
	}
}

func (r *Router) RegisterFiberRoutes() {
	r.MiddlewareLocal.CorsMiddleware()
	r.MiddlewareLocal.SecurityMiddleware()
	if Env.LogsUrl != "" {
		r.MiddlewareLocal.Telemetry("auth/login")
	}

	apiV1 := r.App.Group("/api/v1")
	r.Core.RegisterRouter(apiV1.Group("/core"))
}
