package routes

import (
	"github.com/hhow09/simple_bank/api/controllers"
	"github.com/hhow09/simple_bank/api/middlewares"
	"github.com/hhow09/simple_bank/lib"
)

type TransferRoutes struct {
	controller     controllers.TransferController
	requestHandler lib.RequestHandler
	authMiddleware middlewares.AuthMiddleware
}

// Setup user routes
func (r TransferRoutes) Setup() {
	transferRoutes := r.requestHandler.Gin.Group("/transfers").Use(r.authMiddleware.Handler())
	transferRoutes.POST("", r.controller.CreateTransfer)
}

func NewTransferRoutes(
	controller controllers.TransferController,
	requestHandler lib.RequestHandler,
	authMiddleware middlewares.AuthMiddleware,
) TransferRoutes {
	return TransferRoutes{
		controller,
		requestHandler,
		authMiddleware,
	}
}
