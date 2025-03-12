package middleware

import (
	"net/http"
	"strings"

	"github.com/danisasmita/customer-search/pkg/message"
	"github.com/danisasmita/customer-search/pkg/utils"
	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": message.Unauthorized})
			c.Abort()
			return
		}

		token := strings.Split(authHeader, " ")[1]
		claims, err := utils.ValidateJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": message.Unauthorized})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}
