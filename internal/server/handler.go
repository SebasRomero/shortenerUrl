package server

import "net/http"

func InitRoutes() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("GET /shrt/{url}", getUrlShortened)
	router.HandleFunc("POST /shrt", createUrlShortened)

	return router
}
