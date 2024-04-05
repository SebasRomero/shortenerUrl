package main

import (
	"fmt"
	"net/http"

	"github.com/sebasromero/shortenerUrl/internal/server"
)

func main() {

	fmt.Println("Listen at: 8080")
	http.ListenAndServe(":8080", server.InitRoutes())
}
