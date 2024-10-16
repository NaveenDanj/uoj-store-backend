package router

import (
	"peer-store/handlers"
	"peer-store/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

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
		authApi.POST("/admin-account-setup", handlers.AdminAccountSetup)
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

	adminApi := r.Group("/api/admin")
	adminApi.Use(middlewares.AdminUserAuthRequired())
	{
		adminApi.POST("/change-account-status", handlers.ActivateAccount)
		adminApi.POST("/create-admin", handlers.CreateAdminUser)
	}

	folderApi := r.Group("/api/folder")
	folderApi.Use(middlewares.UserAuthRequired())
	{
		folderApi.POST("/create-folder", handlers.CreateFolder)
		folderApi.GET("/get-sub-folders/:parentId", handlers.GetSubFolders)
		folderApi.POST("/rename-folder", handlers.RenameFolder)
		folderApi.DELETE("/delete-folder", handlers.DeleteFolder)
		folderApi.GET("/get-folder-items/:id", handlers.GetFolderItems)
		// folderApi.GET("/get-user-folders", handlers.GetUserFolders)
		// folderApi.GET("/get-folder/:folderId", handlers.GetFolder)
		// folderApi.PUT("/update-folder/:folderId", handlers.UpdateFolder)
		// folderApi.DELETE("/delete-folder/:folderId", handlers.DeleteFolder)
	}

	return r
}
