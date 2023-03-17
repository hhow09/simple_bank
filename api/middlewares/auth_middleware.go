package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hhow09/simple_bank/constants"
	"github.com/hhow09/simple_bank/token"
	"github.com/hhow09/simple_bank/util"
)

type AuthMiddleware struct {
	tokenMaker token.Maker
}

// Setup sets up jwt auth middleware
func (m AuthMiddleware) Setup() {}

// Handler handles middleware functionality
func (m AuthMiddleware) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(constants.AuthHeaderKey)
		if len(authHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ErrorResponse(err))
			return
		}
		fields := strings.Fields(authHeader) //split authHeader
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ErrorResponse(err))
			return
		}

		authType := strings.ToLower(fields[0])
		if authType != constants.AuthTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ErrorResponse(err))
			return
		}
		accessToken := fields[1]
		payload, err := m.tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ErrorResponse(err))
			return
		}
		ctx.Set(constants.AuthPayloadKey, payload)
		ctx.Next()
	}
}

func NewAuthMiddleware(
	tokenMaker token.Maker,
) AuthMiddleware {
	return AuthMiddleware{
		tokenMaker: tokenMaker,
	}
}
