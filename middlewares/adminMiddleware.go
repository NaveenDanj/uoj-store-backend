package middlewares

import (
	"net/http"
	"peer-store/service/auth"
	"strings"

	"github.com/gin-gonic/gin"
)

func AdminUserAuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {

		authToken := c.GetHeader("Authorization")

		if authToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthenticated"})
			c.Abort()
			return
		}

		if authToken == "" || !strings.HasPrefix(authToken, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		authToken = strings.TrimPrefix(authToken, "Bearer ")

		_, username, err := auth.VerifyJWT(authToken)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthenticated"})
			c.Abort()
			return
		}

		if auth.IsBlocked(authToken , false) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthenticated"})
			c.Abort()
			return
		}

		user, err := auth.GetUserByUsername(username)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthenticated"})
			c.Abort()
			return
		}

		if user.Role != "Admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to perform this action"})
			c.Abort()
			return
		}

		if !user.IsVerified {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Account not verified"})
			c.Abort()
			return
		}

		if !user.IsActive {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Account is not activated"})
			c.Abort()
			return
		}

		c.Set("currentUser", user)

		c.Next()
	}
}
