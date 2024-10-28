package router

import (
	"log"
	"peer-store/handlers"
	"peer-store/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://uoj.uk.to", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Authorization"},
	}))

	r.Use(func(c *gin.Context) {
		log.Printf("Request from origin: %s", c.Request.Header.Get("Origin"))
		if c.Request.Method == "OPTIONS" {
			log.Println("CORS preflight request")
		}
		c.Next()
	})

	r.OPTIONS("/*path", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "https://uoj.uk.to")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.AbortWithStatus(204)
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	authApi := r.Group("/api/auth")
	authApi.Use(func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	{
		authApi.POST("/sign-up", handlers.CreateNewUser)
		authApi.POST("/sign-in", handlers.UserSignIn)
		authApi.POST("/private-session-sign-in", handlers.PrivateSessionSignIn)
		authApi.POST("/verify-account", handlers.VerifyAccount)
		authApi.POST("/reset-password-send-link", handlers.ResetPasswordSendMail)
		authApi.POST("/reset-password", handlers.ResetPasswordSetPassword)
		authApi.POST("/admin-account-setup", handlers.AdminAccountSetup)
		authApi.GET("/current-user", middlewares.UserAuthRequired(), handlers.GetCurrentUser)
		authApi.GET("/session-current-user", middlewares.UserSessionAuthRequired(), handlers.GetCurrentUser)
		authApi.GET("/get-user-pubkey/:userId", middlewares.UserAuthRequired(), handlers.GetUserPublicKey)
		authApi.POST("/logout", middlewares.UserAuthRequired(), handlers.Logout)
		authApi.POST("/session-logout", middlewares.UserSessionAuthRequired(), handlers.Logout)
		authApi.GET("/user-notifications", middlewares.UserAuthRequired(), handlers.GetUserNotifications)
		authApi.GET("/notifications-mark-as-read", middlewares.UserAuthRequired(), handlers.MarkNotificationAsRead)
		authApi.POST("/update-user-profile", middlewares.UserAuthRequired(), handlers.UpdateUserProfile)
		authApi.POST("/check-passphrase", middlewares.UserAuthRequired(), handlers.CheckPassPhrase)
	}

	fileApi := r.Group("/api/file")
	fileApi.Use(middlewares.UserAuthRequired())
	{
		fileApi.POST("/upload", handlers.UploadFile)
		fileApi.POST("/download", handlers.DownloadFile)
		fileApi.POST("/move-to-trash", handlers.MoveFileToTrash)
		fileApi.POST("/move-file", handlers.MoveFile)
		fileApi.DELETE("/delete", handlers.DeleteFile)
		fileApi.GET("/get-user-files", handlers.GetUserFiles)
		fileApi.POST("/change-file-fav-state", handlers.ChangeFileFavState)
	}

	sessionApi := r.Group("/api/session")
	sessionApi.Use(middlewares.UserSessionAuthRequired())
	{
		sessionApi.POST("/upload-session-file", handlers.UploadSessionFile)
		sessionApi.GET("/get-folder-items/:id", handlers.GetFolderItems)

	}

	shareApi := r.Group("/api/share")
	shareApi.Use(middlewares.UserAuthRequired())
	{
		shareApi.POST("/generate-link", handlers.GenerateLink)
		shareApi.POST("/revoke-link", handlers.RevokeLink)
		shareApi.GET("/file/:token", handlers.DownloadSharedFile)
		shareApi.GET("/search-user/:query", handlers.GetUsersToShare)
	}

	adminApi := r.Group("/api/admin")
	adminApi.Use(middlewares.AdminUserAuthRequired())
	{
		adminApi.POST("/change-account-status", handlers.ActivateAccount)
		adminApi.POST("/create-admin", handlers.CreateAdminUser)
		adminApi.GET("/fetch-users", handlers.GetAllUsers)
		adminApi.GET("/fetch-files", handlers.GetAllUserFiles)
	}

	folderApi := r.Group("/api/folder")
	folderApi.Use(middlewares.UserAuthRequired())
	{
		folderApi.POST("/create-folder", handlers.CreateFolder)
		folderApi.GET("/get-sub-folders/:parentId", handlers.GetSubFolders)
		folderApi.POST("/rename-folder", handlers.RenameFolder)
		folderApi.POST("/move-folder", handlers.MoveFolder)
		folderApi.DELETE("/delete-folder", handlers.DeleteFolder)
		folderApi.POST("/move-folder-trash", handlers.MoveToTrash)
		folderApi.GET("/get-folder-items/:id", handlers.GetFolderItems)
	}

	analyticsApi := r.Group("/api/analytics")
	analyticsApi.Use(middlewares.UserAuthRequired())
	{
		// folderApi.POST("/create-folder", handlers.CreateFolder)
		// folderApi.GET("/get-sub-folders/:parentId", handlers.GetSubFolders)
		// folderApi.POST("/rename-folder", handlers.RenameFolder)
		// folderApi.DELETE("/delete-folder", handlers.DeleteFolder)
		// folderApi.GET("/get-folder-items/:id", handlers.GetFolderItems)
	}

	return r
}
