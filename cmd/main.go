package main

import (
	"fmt"
	"peer-store/db"
	"peer-store/router"
)

func main() {

	// setup databases
	_, err := db.SetupDatabase()

	if err != nil {
		fmt.Errorf(err.Error())
	}

	// Set up the router using the separate router file
	r := router.SetupRouter()

	// Start the server on port 8080
	r.Run(":5001")
}
