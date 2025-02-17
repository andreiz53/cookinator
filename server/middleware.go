package server

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/andreiz53/cookinator/token"
)

const (
	authHeaderKey        = "Authorization"
	authHeaderTypeBearer = "Bearer"
	ctxAuthPayloadKey    = "auth_key"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(authHeaderKey)
		if authHeader == "" {
			err := errors.New("please provide an authorization header")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, respondWithErorr(err))
			return
		}

		authFields := strings.Fields(authHeader)
		if len(authFields) != 2 {
			err := errors.New("invalid authorization header")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, respondWithErorr(err))
			return
		}

		authType := authFields[0]
		if authType != authHeaderTypeBearer {
			err := errors.New("invalid authorization type")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, respondWithErorr(err))
			return
		}

		token := authFields[1]
		payload, err := tokenMaker.VerifyToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, respondWithErorr(err))
			return
		}

		ctx.Set(ctxAuthPayloadKey, payload)
		ctx.Next()
	}
}
