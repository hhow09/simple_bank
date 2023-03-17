package routes

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewUserRoutes),
	fx.Provide(NewAccountRoutes),
	fx.Provide(NewTransferRoutes),
	// add more here
	fx.Provide(NewSwaggerRoutes),
	fx.Provide(NewRoutes),
)

type Routes []Route

// Route interface
type Route interface {
	Setup()
}

// NewRoutes sets up routes
func NewRoutes(
	userRoutes UserRoutes,
	swaggerRoutes SwaggerRoutes,
	accountRoutes AccountRotes,
	transferRoutes TransferRoutes,
) Routes {
	return Routes{
		userRoutes,
		accountRoutes,
		transferRoutes,
		swaggerRoutes,
	}
}

// Setup all the route
func (r Routes) Setup() {
	for _, route := range r {
		route.Setup()
	}
}
