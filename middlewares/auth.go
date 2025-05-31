package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/AntonioGuilhermeDev/InventoryHubApis/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token obrigatório"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Formato do token inválido"})
			return
		}

		tokenString := parts[1]

		secretKey, err := utils.GetSecretKey()
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Erro interno"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("método de assinatura inválido")
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token inválido"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token inválido"})
			return
		}

		ctx.Set("role", claims["role"])
		ctx.Set("userId", claims["userId"])

		ctx.Next()
	}
}
