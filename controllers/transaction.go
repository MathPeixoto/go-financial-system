package controllers

import (
	"bancario/dao"
	"bancario/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type transactionModel model.TransactionDb
type Transaction struct {
	transactionDao dao.TransactionDb
}

func (t *Transaction) GetTransactionsBySourceId(c *gin.Context) {
	sourceId := c.Param("sourceId")

	transactionsResponse := t.transactionDao.GetTransactionsBySourceId(sourceId)
	if transactionsResponse != nil {
		c.IndentedJSON(http.StatusOK, transactionsResponse)
		return
	}

	c.IndentedJSON(http.StatusNoContent, []transactionModel{})
}

func (t *Transaction) GetTransactionsByDestinationId(c *gin.Context) {
	destinationId := c.Param("destinationId")

	transactionsResponse := t.transactionDao.GetTransactionsByDestinationId(destinationId)
	if transactionsResponse != nil {
		c.IndentedJSON(http.StatusOK, transactionsResponse)
		return
	}

	c.IndentedJSON(http.StatusNoContent, []transactionModel{})
}
