package middleware

import (
	services "pg_sandbox/services/auth_services"
	"pg_sandbox/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthorizationMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			utils.RespondWithError(ctx, 401, "Authorization header is required")
			ctx.Abort() // Abort the request to prevent further execution
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := services.ValidateToken(tokenString)
		if err != nil {
			utils.RespondWithError(ctx, 401, "Invalid Token")
			ctx.Abort() // Abort the request to prevent further execution
			return
		}

		ctx.Set("userID", claims.Subject)
		ctx.Next()
	}
}
