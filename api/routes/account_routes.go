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
	authRoutes := r.requestHandler.Gin.Group("/accounts").Use(r.authMiddleware.Handler())
	authRoutes.POST("", r.controller.CreateAccount)
	authRoutes.GET("/:id", r.controller.GetAccount)
	authRoutes.GET("", r.controller.ListAccounts)
}

func NewAccountRotes(
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
