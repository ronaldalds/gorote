package main

import (
	"github.com/gofiber/fiber/v2"

	"github.com/ronaldalds/gorote-core/core"
)

type Router struct {
	App             *fiber.App
	MiddlewareLocal *Middleware
	Core            *core.Router
}

func Ready() error {
	return nil
}

func Config(app *fiber.App) *Router {
	coreGorm := &InitGorm{
		Host:     Env.Sql.Host,
		User:     Env.Sql.Username,
		Password: Env.Sql.Password,
		Database: Env.Sql.Database,
		Port:     Env.Sql.Port,
		Schema:   Env.Sql.Schema,
		TimeZone: Env.App.TimeZone,
	}
	jwt := core.AppJwt{
		AppName:          Env.App.Name,
		TimeZone:         Env.App.TimeZone,
		JwtSecret:        Env.Jwt.Secret,
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
		App:       app,
		GormStore: newGormStore(coreGorm),
		Jwt:       jwt,
		Super:     super,
	}
	return &Router{
		App:             app,
		MiddlewareLocal: NewMiddleware(app),
		Core:            core.New(coreConfig),
	}
}
