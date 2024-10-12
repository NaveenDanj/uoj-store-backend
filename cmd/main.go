package main

import (
	"peer-store/db"
	"peer-store/router"
)

func main() {
	db.Setup()
	db.SeedAdminAccount()
	r := router.SetupRouter()
	r.Run(":5001")
}
