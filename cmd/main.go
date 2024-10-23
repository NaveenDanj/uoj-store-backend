package main

import (
	"log"
	"peer-store/db"
	"peer-store/router"
)

func main() {
	db.Setup()
	db.SeedAdminAccount()

	r := router.SetupRouter()
	r.MaxMultipartMemory = 100 << 20

	log.Println("Starting server on :5001")
	err := r.Run(":5001")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
