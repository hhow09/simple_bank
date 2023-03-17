package routes

import (
	"github.com/hhow09/simple_bank/docs"
	"github.com/hhow09/simple_bank/lib"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type SwaggerRoutes struct {
	requestHandler lib.RequestHandler
}

// Setup swagger routes
func (r SwaggerRoutes) Setup() {
	docs.SwaggerInfo.BasePath = "/"
	r.requestHandler.Gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

func NewSwaggerRoutes(
	requestHandler lib.RequestHandler,
) SwaggerRoutes {
	return SwaggerRoutes{requestHandler: requestHandler}
}
