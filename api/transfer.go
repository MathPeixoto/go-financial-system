package api

import (
	"database/sql"
	"fmt"
	db "github.com/MathPeixoto/go-financial-system/db/sqlc"
	"github.com/gin-gonic/gin"
	"net/http"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

type idTransferRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) createTransfer(c *gin.Context) {
	var request transferRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !server.validAccount(c, request.FromAccountID, request.Currency) {
		return
	}

	if !server.validAccount(c, request.ToAccountID, request.Currency) {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: request.FromAccountID,
		ToAccountID:   request.ToAccountID,
		Amount:        request.Amount,
	}

	result, err := server.store.TransferTx(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, result)
}

func (server *Server) getTransfer(c *gin.Context) {
	var request idTransferRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	transfer, err := server.store.GetTransfer(c, request.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, transfer)
}

func (server *Server) validAccount(c *gin.Context, accountID int64, currency string) bool {
	account, err := server.store.GetAccount(c, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, errorResponse(err))
			return false
		}

		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account %d has currency %s, but transfer currency is %s", accountID, account.Currency, currency)
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}

	return true
}
