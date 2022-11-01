package api

import (
	"database/sql"
	db "github.com/MathPeixoto/go-financial-system/db/sqlc"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CreateAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
	Amount   int64  `json:"amount"`
}

type IdAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type UpdateAccountBalanceRequest struct {
	Amount int64 `json:"amount" binding:"required,min=1"`
}

type ListAccountsRequest struct {
	Limit  int32 `form:"limit,default=5" binding:"min=5,max=10"`
	Offset int32 `form:"offset,default=1" binding:"min=1"`
}

func (server *Server) createAccount(c *gin.Context) {
	var account CreateAccountRequest
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

func (server *Server) getAccount(c *gin.Context) {
	var request IdAccountRequest
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

func (server *Server) listAccounts(c *gin.Context) {
	var request ListAccountsRequest
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
	var requestId IdAccountRequest
	if err := c.ShouldBindUri(&requestId); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var requestAccount UpdateAccountBalanceRequest
	if err := c.ShouldBindJSON(&requestAccount); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.AddAccountBalanceParams{
		ID:     requestId.ID,
		Amount: requestAccount.Amount,
	}

	accountUpdated, err := server.store.AddAccountBalance(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, accountUpdated)
}

func (server *Server) deleteAccount(c *gin.Context) {
	var request IdAccountRequest
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
