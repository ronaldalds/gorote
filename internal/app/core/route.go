package core

import (
	"github.com/ronaldalds/base-go-api/internal/config/access"
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

func (r *Router) Auth(router fiber.Router) {
	router.Post(
		"/login",
		r.Middleware.ValidationMiddleware(&Login{}, "json"),
		r.Controller.LoginHandler,
	)
}

func (r *Router) RegisterRouter(router fiber.Router) {
	router.Get("/health", r.Controller.HealthHandler)
	// Authentication
	authGroup := router.Group("/auth", r.Middleware.Limited(10))
	r.Auth(authGroup)

	// Users
	usersGroup := router.Group("/users")
	r.User(usersGroup)
	r.Role(usersGroup)
	r.Permission(usersGroup)
}

func (r *Router) User(router fiber.Router) {
	router.Get(
		"/",
		r.Middleware.ValidationMiddleware(&Paginate{}, "query"),
		r.Middleware.JWTProtected(access.Permissions.ViewUser),
		r.Controller.ListUserHandler,
	)
	router.Post(
		"/",
		r.Middleware.ValidationMiddleware(&CreateUser{}, "json"),
		r.Middleware.JWTProtected(access.Permissions.CreateUser),
		r.Controller.CreateUserHandler,
	)
	router.Put(
		"/:id",
		r.Middleware.ValidationMiddleware(&UserParam{}, "params"),
		r.Middleware.ValidationMiddleware(&User{}, "json"),
		r.Middleware.JWTProtected(),
		r.Controller.UpdateUserHandler,
	)
}

func (r *Router) Role(router fiber.Router) {
	router.Get(
		"/roles",
		r.Middleware.ValidationMiddleware(&Paginate{}, "query"),
		r.Middleware.JWTProtected(),
		r.Controller.ListRoleHandler,
	)
	router.Post(
		"/roles",
		r.Middleware.ValidationMiddleware(&CreateRole{}, "json"),
		r.Middleware.JWTProtected(access.Permissions.CreateRole),
		r.Controller.CreateRoleHandler,
	)
}

func (r *Router) Permission(router fiber.Router) {
	router.Get(
		"/permissions",
		r.Middleware.ValidationMiddleware(&Paginate{}, "query"),
		r.Middleware.JWTProtected(),
		r.Controller.ListPermissiontHandler,
	)
}
