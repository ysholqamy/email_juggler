package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ysholqamy/email_juggler/email"
)

func main() {
	emailService := email.DefaultService
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	http.Handle("/messages", emailService)
	err := http.ListenAndServe(":"+port, logRequest(http.DefaultServeMux))
	if err != nil {
		log.Fatal(err)
	}
}

// A wrapper around an http.Handler to log incoming requests
func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s %s\n", r.RemoteAddr, r.Method, r.URL, r.Header.Get("Content-Type"))
		handler.ServeHTTP(w, r)
	})
}
