package main

import (
	"peer-store/router"
)

func main() {
	// Set up the router using the separate router file
	r := router.SetupRouter()

	// Start the server on port 8080
	r.Run(":5000")
}
