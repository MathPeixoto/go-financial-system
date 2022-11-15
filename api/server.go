package api

import (
	db "github.com/MathPeixoto/go-financial-system/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := validate.RegisterValidation("currency", validCurrencies)
		if err != nil {
			return nil
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
	return server
}

// Start starts the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
