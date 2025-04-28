package main

import (
	"github.com/gofiber/fiber/v2"

	"github.com/ronaldalds/gorote-core-rsa/core"
)

type Router struct {
	App             *fiber.App
	MiddlewareLocal *Middleware
	Core            *core.Router
}

func Config(app *fiber.App) *Router {
	coreGorm := &core.InitGorm{
		Host:     Env.Sql.Host,
		User:     Env.Sql.Username,
		Password: Env.Sql.Password,
		Database: Env.Sql.Database,
		Port:     Env.Sql.Port,
		Schema:   Env.Sql.Schema,
		TimeZone: Env.App.TimeZone,
	}
	jwt := &core.AppJwt{
		JwtExpireAccess:  Env.Jwt.ExpireAccess,
		JwtExpireRefresh: Env.Jwt.ExpireRefresh,
	}
	super := &core.AppSuper{
		SuperName:  Env.Super.Name,
		SuperUser:  Env.Super.Username,
		SuperPass:  Env.Super.Password,
		SuperEmail: Env.Super.Email,
		SuperPhone: Env.Super.Phone,
	}

	coreConfig := &core.AppConfig{
		App:         app,
		AppName:     Env.App.Name,
		AppTimeZone: Env.App.TimeZone,
		CoreGorm:    core.NewGormStore(coreGorm),
		Jwt:         jwt,
		Super:       super,
	}
	return &Router{
		App:             app,
		MiddlewareLocal: NewMiddleware(app),
		Core:            core.New(coreConfig),
	}
}
