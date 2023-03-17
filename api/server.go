package api

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/hhow09/simple_bank/db/sqlc"
	docs "github.com/hhow09/simple_bank/docs"
	"github.com/hhow09/simple_bank/token"
	"github.com/hhow09/simple_bank/util"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

func NewServer(config util.Config, store db.Store, tokenMaker token.Maker) (*Server, error) {
	server := &Server{store: store, tokenMaker: tokenMaker, config: config}
	//binding custom validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		//registor validator to gin
		v.RegisterValidation("currency", validCurrency)
	}
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	docs.SwaggerInfo.BasePath = "/"
	router := gin.Default()
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/accounts", server.CreateAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccounts)

	authRoutes.POST("/transfers", server.CreateTransfer)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	//add routes to router
	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
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

var Module = fx.Options(
	fx.Provide(NewServer),
	fx.Invoke(registerHooks),
)
