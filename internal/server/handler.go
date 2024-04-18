package server

import "net/http"

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		next.ServeHTTP(w, r)
	})
}

func InitRoutes() *http.ServeMux {
	router := http.NewServeMux()
	router.Handle("GET /{url}", cors(http.HandlerFunc(getUrlShortened)))
	router.Handle("POST /", cors(http.HandlerFunc(createUrlShortened)))

	return router
}
