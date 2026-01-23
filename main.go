package main

import (
	"fmt"
	"net/http"
)

func main() {
	initDB()
	defer DB.Close() // Ensures DB closes when server stops

	router := registerRoutes()

	wrappedRoute := loggingMiddleware(router)

	server := &http.Server{
		Addr:    ":8000",
		Handler: wrappedRoute,
	}

	fmt.Println("Server running on http://localhost:8000")

	server.ListenAndServe()

}
