package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleRaw, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Sem papel definido no token"})
			return
		}

		role, ok := roleRaw.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Formato inválido de papel"})
			return
		}

		role = strings.ToUpper(role) // garantir que esteja em maiúsculo

		for _, allowed := range allowedRoles {
			if role == strings.ToUpper(allowed) {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Permissão negada para esta rota"})
	}
}
