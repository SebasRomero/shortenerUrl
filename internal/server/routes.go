package server

import (
	"encoding/json"
	"net/http"

	"github.com/sebasromero/shortenerUrl/internal/database"
	"github.com/sebasromero/shortenerUrl/internal/types"
	"github.com/sebasromero/shortenerUrl/internal/url"
)

func getUrlShortened(w http.ResponseWriter, r *http.Request) {
	customUrl := r.URL.Path
	url, err := database.Connection.GetUrlShortened(types.Path + customUrl)
	if err != nil {
		http.Redirect(w, r, types.Path, http.StatusSeeOther)
		return
	}
	if url.LongUrl != "" {
		http.Redirect(w, r, url.LongUrl, http.StatusSeeOther)
		return
	}
}

func createUrlShortened(w http.ResponseWriter, r *http.Request) {
	longUrl := &types.LongUrlResponse{}
	err := json.NewDecoder(r.Body).Decode(longUrl)
	if err != nil {
		w.Write([]byte("Error decoding"))
		return
	}
	response, err := url.CreateShortenedUrl(longUrl.LongUrl)
	if err != nil {
		w.Write([]byte("Error creating the url"))
		return
	}
	json.NewEncoder(w).Encode(&response)
}
