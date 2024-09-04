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
			return
		}

		_, email, err := auth.VerifyJWT(authToken)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := auth.GetUserByEmail(email)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.Set("currentUser", user)

		c.Next()
	}
}
