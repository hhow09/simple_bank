package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hhow09/simple_bank/token"
)

const (
	authHeaderKey  = "authorization"
	authTypeBearer = "bearer"
	authPayloadKey = "auth_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(authHeaderKey)
		if len(authHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		fields := strings.Fields(authHeader) //split authHeader
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authType := strings.ToLower(fields[0])
		if authType != authTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		ctx.Set(authPayloadKey, payload)
		ctx.Next()
	}
}
