package routes

import (
	"github.com/hhow09/simple_bank/api/controllers"
	"github.com/hhow09/simple_bank/api/middlewares"
	"github.com/hhow09/simple_bank/lib"
)

type AccountRotes struct {
	controller     controllers.AccountController
	requestHandler lib.RequestHandler
	authMiddleware middlewares.AuthMiddleware
}

// Setup user routes
func (r AccountRotes) Setup() {
	accountRoutes := r.requestHandler.Gin.Group("/accounts").Use(r.authMiddleware.Handler())
	accountRoutes.POST("", r.controller.CreateAccount)
	accountRoutes.GET("/:id", r.controller.GetAccount)
	accountRoutes.GET("", r.controller.ListAccounts)
}

func NewAccountRoutes(
	controller controllers.AccountController,
	requestHandler lib.RequestHandler,
	authMiddleware middlewares.AuthMiddleware,
) AccountRotes {
	return AccountRotes{
		controller,
		requestHandler,
		authMiddleware,
	}
}
