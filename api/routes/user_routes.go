package routes

import (
	"github.com/hhow09/simple_bank/api/controllers"
	"github.com/hhow09/simple_bank/lib"
)

type UserRoutes struct {
	controller     controllers.UserController
	requestHandler lib.RequestHandler
}

// Setup user routes
func (r UserRoutes) Setup() {
	users := r.requestHandler.Gin.Group("/users")
	users.POST("", r.controller.CreateUser)
	users.POST("/login", r.controller.LoginUser)
}

func NewUserRoutes(
	controller controllers.UserController,
	requestHandler lib.RequestHandler,
) UserRoutes {
	return UserRoutes{
		controller,
		requestHandler,
	}
}
