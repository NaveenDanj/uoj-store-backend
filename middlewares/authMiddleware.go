package middlewares

import (
	"net/http"
	"peer-store/service/auth"

	"github.com/gin-gonic/gin"
)

func UserAuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {

		authToken := c.GetHeader("Authorization")

		if authToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthenticated"})
			c.Abort()
			return
		}

		_, username, err := auth.VerifyJWT(authToken)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthenticated"})
			c.Abort()
			return
		}

		if auth.IsBlocked(authToken) {
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
