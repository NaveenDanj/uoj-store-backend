package main

import (
	"log"
	"peer-store/db"
	"peer-store/router"

	"github.com/gin-gonic/gin"
)

func main() {
	db.Setup()
	db.SeedAdminAccount()

	r := router.SetupRouter()
	r.MaxMultipartMemory = 100 << 20

	r.Use(func(c *gin.Context) {
		log.Printf("Incoming request: %s %s from %s", c.Request.Method, c.Request.URL.Path, c.ClientIP())
		c.Next()
	})

	log.Println("Starting server on :5001")
	err := r.Run(":5001")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
