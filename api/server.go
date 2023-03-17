package api

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/hhow09/simple_bank/api/controllers"
	"github.com/hhow09/simple_bank/api/middlewares"
	"github.com/hhow09/simple_bank/api/routes"
	db "github.com/hhow09/simple_bank/db/sqlc"
	"github.com/hhow09/simple_bank/lib"
	"github.com/hhow09/simple_bank/token"
	"github.com/hhow09/simple_bank/util"
	"go.uber.org/fx"
)

// @title Simple Bank API
// @version 1.0
// @description A simple bank service.

// @contact.name API Support
// @contact.url https://github.com/hhow09/simple_bank/issues
// @contact.email hhow09@gmail.com

// @host localhost:8080
// @BasePath /

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey authorization
// @in header
// @name Authorization

type Server struct {
	config     util.Config
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.Store, tokenMaker token.Maker, requestHandler lib.RequestHandler) (*Server, error) {
	server := &Server{store: store, tokenMaker: tokenMaker, config: config, router: requestHandler.Gin}
	//binding custom validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		//registor validator to gin
		v.RegisterValidation("currency", validCurrency)
	}
	// server.setupRouter()
	return server, nil
}

// func (server *Server) setupRouter() {

// 	authRoutes.POST("/transfers", server.CreateTransfer)
// }

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func registerHooks(lc fx.Lifecycle, server *Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() error {
				err := server.Start(server.config.ServerAddress)
				if err != nil {
					log.Fatal("error starting server: ", err)
				}
				return nil
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("Stopping server")
			return nil
		},
	})
}

func setupRoutes(lc fx.Lifecycle, r routes.Routes) {
	r.Setup()
}
func setupMiddleware(lc fx.Lifecycle, m middlewares.Middlewares) {
	m.Setup()
}

var Module = fx.Options(
	controllers.Module,
	routes.Module,
	middlewares.Module,
	fx.Provide(NewServer),
	fx.Invoke(setupRoutes),
	fx.Invoke(setupMiddleware),
	fx.Invoke(registerHooks),
)
