package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type RenewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RenewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) renewAccessTokenUser(c *gin.Context) {
	var request RenewAccessTokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	refreshPayload, err := server.tokenMaker.VerifyToken(request.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	session, err := server.store.GetSession(c, refreshPayload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if session.IsBlocked {
		err := fmt.Errorf("session %s is blocked", session.ID)
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.Username != refreshPayload.Username {
		err := fmt.Errorf("session %s does not belong to user %s", session.ID, refreshPayload.Username)
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.RefreshToken != request.RefreshToken {
		err := fmt.Errorf("refresh token does not match")
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.ExpiresAt.Before(time.Now()) {
		err := fmt.Errorf("session is expired")
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	token, accessPayload, err := server.tokenMaker.CreateToken(refreshPayload.Username, server.config.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	userResponse := RenewAccessTokenResponse{
		AccessToken:          token,
		AccessTokenExpiresAt: accessPayload.ExpiresAt,
	}

	c.JSON(http.StatusOK, userResponse)
}
