package server

import (
	"encoding/json"
	"net/http"

	dabatase "github.com/sebasromero/shortenerUrl/internal/database"
	"github.com/sebasromero/shortenerUrl/internal/types"
)

var db = dabatase.Connect()

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}

func getUrlShortened(w http.ResponseWriter, r *http.Request) {
	customUrl := r.URL.Path
	url := db.GetUrlShortened(types.Path + customUrl)
	if url != nil {
		http.Redirect(w, r, url.LongUrl, http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, types.Path, http.StatusSeeOther)
}

func createUrlShortened(w http.ResponseWriter, r *http.Request) {
	url := &types.InputLongUrl{}
	err := json.NewDecoder(r.Body).Decode(url)
	if err != nil {
		w.Write([]byte("Error decoding"))
		return
	}
	response := db.CreateShortenerUrl(url.LongUrl)
	json.NewEncoder(w).Encode(&response)
}
