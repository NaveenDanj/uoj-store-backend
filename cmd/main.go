package main

import (
	"peer-store/db"
	"peer-store/router"
)

func main() {
	db.Setup()
	r := router.SetupRouter()
	r.Run(":5001")
}
