package main

import (
	"log"
	"net/http"
	"peer-store/db"
	"peer-store/router"
)

func main() {
	db.Setup()
	db.SeedAdminAccount()

	go func() {
		log.Fatal(http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "https://"+r.Host+r.URL.String(), http.StatusMovedPermanently)
		})))
	}()

	r := router.SetupRouter()
	r.Run(":5001")
	// log.Fatal(http.ListenAndServeTLS(":443", "/etc/ssl/certs/selfsigned.crt", "/etc/ssl/private/selfsigned.key", r))
}
