package server

import "net/http"

func InitRoutes() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("GET /{url}", getUrlShortened)
	router.HandleFunc("POST /", createUrlShortened)

	return router
}
