package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/sebasromero/shortenerUrl/internal/server"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Listen at:", port)
	http.ListenAndServe(":"+port, server.InitRoutes())
}
