package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/sebasromero/shortenerUrl/internal/database"
	"github.com/sebasromero/shortenerUrl/internal/server"
)

var mainConnection = database.Connection

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	go func() {
		for range time.Tick(time.Minute * 10) {
			mainConnection.RemoveExpiredShortUrls()
		}
	}()

	fmt.Println("Listen at:", port)
	http.ListenAndServe(":"+port, server.InitRoutes())
}
