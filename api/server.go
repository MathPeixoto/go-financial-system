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
	router := gin.Default()

	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := validate.RegisterValidation("currency", validCurrencies)
		if err != nil {
			return nil, err
		}
	}

	// users routes
	router.POST("/users", server.createUser)

	// account routes
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	router.PATCH("/accounts/:id", server.updateAccountBalance)
	router.DELETE("/accounts/:id", server.deleteAccount)

	// transfer routes
	router.POST("/transfers", server.createTransfer)
	router.GET("/transfers/:id", server.getTransfer)

	server.router = router
	return server, nil
}

// Start starts the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
