package router

import (
	"peer-store/handlers"
	"peer-store/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	api := r.Group("/api/auth")
	{
		api.POST("/sign-up", handlers.CreateNewUser)
		api.POST("/sign-in", handlers.UserSignIn)
		api.GET("/current-user", middlewares.UserAuthRequired(), handlers.GetCurrentUser)
	}

	return r
}
