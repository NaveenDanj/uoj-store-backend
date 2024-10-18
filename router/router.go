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
		AllowOrigins:     []string{"http://15.206.79.187", "http://localhost:5173"},
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
		c.Header("Access-Control-Allow-Origin", "http://15.206.79.187") // Match your frontend origin
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
	}

	return r
}
