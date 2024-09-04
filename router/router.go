package router

import (
	"peer-store/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRouter initializes and returns a new Gin router
func SetupRouter() *gin.Engine {
	// Create a new Gin router instance
	r := gin.Default()

	// Define routes
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// User-related routes
	api := r.Group("/api/auth")
	{
		api.POST("/sign-up", handlers.CreateNewUser)
		api.POST("/sign-in", handlers.UserSignIn)
		// api.GET("/users/:id", handlers.GetUser)
		// api.PUT("/users/:id", handlers.UpdateUser)
		// api.DELETE("/users/:id", handlers.DeleteUser)
	}

	print(api)

	return r
}
