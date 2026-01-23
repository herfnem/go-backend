package main

import (
	"net/http"
)

func registerRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handleHome)

	mux.HandleFunc("GET /post/{slug}", handlePostView)
	mux.HandleFunc("POST /post/create", handlePostCreate)
	mux.HandleFunc("POST /user/create", handleCreateUser)
	mux.HandleFunc("GET /user/{id}", handleGetUser)
	mux.HandleFunc("GET /users", handleGetUsers)

	return mux
}
