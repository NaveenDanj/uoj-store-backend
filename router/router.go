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

	authApi := r.Group("/api/auth")
	{
		authApi.POST("/sign-up", handlers.CreateNewUser)
		authApi.POST("/sign-in", handlers.UserSignIn)
		authApi.POST("/verify-account", handlers.VerifyAccount)
		authApi.POST("/reset-password-send-link", handlers.ResetPasswordSendMail)
		authApi.POST("/reset-password", handlers.ResetPasswordSetPassword)
		authApi.GET("/current-user", middlewares.UserAuthRequired(), handlers.GetCurrentUser)
		authApi.GET("/get-user-pubkey/:userId", middlewares.UserAuthRequired(), handlers.GetUserPublicKey)
		authApi.POST("/logout", middlewares.UserAuthRequired(), handlers.Logout)
	}

	fileApi := r.Group("/api/file")
	fileApi.Use(middlewares.UserAuthRequired())
	{
		fileApi.POST("/upload", handlers.UploadFile)
		fileApi.POST("/download/", handlers.DownloadFile)
		fileApi.DELETE("/delete/", handlers.DeleteFile)
		fileApi.GET("/get-user-files", handlers.GetUserFiles)
	}

	shareApi := r.Group("/api/share")
	shareApi.Use(middlewares.UserAuthRequired())
	{
		shareApi.POST("/generate-link", handlers.GenerateLink)
		shareApi.POST("/revoke-link", handlers.RevokeLink)
		shareApi.GET("/file/:token", handlers.DownloadSharedFile)
	}

	return r
}
