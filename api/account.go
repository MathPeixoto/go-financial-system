package api

import (
	"database/sql"
	"errors"
	db "github.com/MathPeixoto/go-financial-system/db/sqlc"
	"github.com/MathPeixoto/go-financial-system/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"net/http"
)

type CreateAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
	Amount   int64  `json:"amount"`
}

type IDAccountRequest struct {
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

	authPayload := c.MustGet(authPayloadKey).(*token.Payload)

	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: account.Currency,
		Balance:  0,
	}

	createAccount, err := server.store.CreateAccount(c, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok { //nolint: errorlint
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				c.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}

		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, createAccount)
}

func (server *Server) getAccount(c *gin.Context) {
	var request IDAccountRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(c, request.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := c.MustGet(authPayloadKey).(*token.Payload)
	if account.Owner != authPayload.Username {
		err = errors.New("account does not belong to the authenticated user")
		c.JSON(http.StatusForbidden, errorResponse(err))
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

	authPayload := c.MustGet(authPayloadKey).(*token.Payload)
	arg := db.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  request.Limit,
		Offset: (request.Offset - 1) * request.Limit,
	}

	accounts, err := server.store.ListAccounts(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, accounts)
}

func (server *Server) updateAccountBalance(c *gin.Context) {
	var requestID IDAccountRequest
	if err := c.ShouldBindUri(&requestID); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := validateAccountID(c, server, requestID)
	if err != nil {
		return
	}

	var requestAccount UpdateAccountBalanceRequest
	if err := c.ShouldBindJSON(&requestAccount); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.AddAccountBalanceParams{
		ID:     requestID.ID,
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
	var request IDAccountRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := validateAccountID(c, server, request)
	if err != nil {
		return
	}

	err = server.store.DeleteAccount(c, request.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, nil)
}

func validateAccountID(c *gin.Context, server *Server, requestID IDAccountRequest) error {
	authPayload := c.MustGet(authPayloadKey).(*token.Payload)
	owner, err := server.store.GetAccountByOwner(c, authPayload.Username)

	if owner.ID != requestID.ID {
		err = errors.New("account does not belong to the authenticated user")
		c.JSON(http.StatusForbidden, errorResponse(err))
		return err
	}

	return nil
}
