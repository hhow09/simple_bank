package middlewares

import "go.uber.org/fx"

// Module Middleware exported
var Module = fx.Options(
	fx.Provide(NewAuthMiddleware),
	fx.Provide(NewMiddlewares),
)

// IMiddleware middleware interface
type IMiddleware interface {
	Setup()
}

// Middlewares contains multiple middleware
type Middlewares []IMiddleware

func NewMiddlewares(
	authMiddleware AuthMiddleware,
) Middlewares {
	return Middlewares{
		authMiddleware,
	}
}

// Setup sets up middlewares
func (m Middlewares) Setup() {
	for _, middleware := range m {
		middleware.Setup()
	}
}
