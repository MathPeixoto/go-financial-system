package api

import (
	"fmt"
	db "github.com/MathPeixoto/go-financial-system/db/sqlc"
	"github.com/MathPeixoto/go-financial-system/token"
	"github.com/MathPeixoto/go-financial-system/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := validate.RegisterValidation("currency", validCurrencies)
		if err != nil {
			return nil, err
		}
	}

	server.setupRoutes()
	return server, nil
}

// setupRoutes sets up the routes for the server.
func (server *Server) setupRoutes() {
	router := gin.Default()

	// users routes
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/token/renew_access", server.renewAccessTokenUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	// account routes
	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccounts)
	authRoutes.PATCH("/accounts/:id", server.updateAccountBalance)
	authRoutes.DELETE("/accounts/:id", server.deleteAccount)

	// transfer routes
	authRoutes.POST("/transfers", server.createTransfer)
	authRoutes.GET("/transfers/:id", server.getTransfer)
	server.router = router
}

// Start starts the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
