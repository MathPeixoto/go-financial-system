package main

import (
	"bancario/controllers"
	"github.com/gin-gonic/gin"
)

func main() {

	c := controllers.Transaction{}

	router := gin.Default()
	router.GET("/transactions/sourceId/:sourceId", c.GetTransactionsBySourceId)
	router.GET("/transactions/destinationId/:destinationId", c.GetTransactionsByDestinationId)
	err := router.Run("localhost:8080")
	if err != nil {
		return
	}
}
