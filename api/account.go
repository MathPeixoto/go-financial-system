package api

import (
	"database/sql"
	db "github.com/MathPeixoto/go-financial-system/db/sqlc"
	"github.com/gin-gonic/gin"
	"net/http"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR BRL"`
	Amount   int64  `json:"amount"`
}

func (server *Server) createAccount(c *gin.Context) {
	var account createAccountRequest
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateAccountParams{
		Owner:    account.Owner,
		Currency: account.Currency,
		Balance:  0,
	}

	createAccount, err := server.store.CreateAccount(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, createAccount)
}

type idAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(c *gin.Context) {
	var request idAccountRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(c, request.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, account)
}

type listAccountsRequest struct {
	Limit  int32 `form:"limit,default=5" binding:"min=5,max=10"`
	Offset int32 `form:"offset,default=1" binding:"min=1"`
}

func (server *Server) listAccounts(c *gin.Context) {
	var request listAccountsRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListAccountsParams{
		Limit: request.Limit, Offset: (request.Offset - 1) * request.Limit,
	}

	accounts, err := server.store.ListAccounts(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, accounts)
}

func (server *Server) updateAccountBalance(c *gin.Context) {
	var request idAccountRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var account createAccountRequest
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.AddAccountBalanceParams{
		ID:     request.ID,
		Amount: account.Amount,
	}

	_, err := server.store.AddAccountBalance(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (server *Server) deleteAccount(c *gin.Context) {
	var request idAccountRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteAccount(c, request.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, nil)
}
