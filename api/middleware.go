package api

import (
	"errors"
	"github.com/MathPeixoto/go-financial-system/token"
	"github.com/gin-gonic/gin"
	"strings"
)

const (
	authTokenHeader = "Authorization"
	authTypeBearer  = "bearer"
	authPayloadKey  = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(authTokenHeader)

		if len(authHeader) == 0 {
			err := errors.New("authorization token is not provided")
			c.AbortWithStatusJSON(401, errorResponse(err))
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			err := errors.New("authorization token is invalid")
			c.AbortWithStatusJSON(401, errorResponse(err))
			return
		}

		authType := strings.ToLower(fields[0])
		if authType != authTypeBearer {
			err := errors.New("authorization type is not supported")
			c.AbortWithStatusJSON(401, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(401, errorResponse(err))
			return
		}

		c.Set(authPayloadKey, payload)
		c.Next()
	}
}
