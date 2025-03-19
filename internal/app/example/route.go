package example

import (
	"github.com/ronaldalds/base-go-api/internal/middlewares"

	"github.com/gofiber/fiber/v2"
)

type Router struct {
	Middleware *middlewares.Middleware
	Controller *Controller
}

func NewRouter(app *fiber.App) *Router {
	return &Router{
		Middleware: middlewares.NewMiddleware(app),
		Controller: NewController(),
	}
}