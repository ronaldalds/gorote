package main

func (r *Router) RegisterFiberRoutes() {
	r.MiddlewareLocal.CorsMiddleware()
	r.MiddlewareLocal.SecurityMiddleware()
	if Env.Logs != nil {
		r.MiddlewareLocal.Telemetry("auth/login")
	}
	apiV1 := r.App.Group("/api/v1")
	r.Core.RegisterRouter(apiV1.Group("/core"))
}
