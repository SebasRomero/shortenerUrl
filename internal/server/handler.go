package server

import "net/http"

func InitRoutes() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("GET /v1/shortener", home)

	return router
}
