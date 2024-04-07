package types

type UrlShortened struct {
	UrlShortened string `json:"urlShortened"`
}

type InputLongUrl struct {
	LongUrl string `json:"longUrl"`
}

type ShortUrlResponse struct {
	LongUrl string `json:"longUrl"`
}

var Path string = "http://localhost:8080"
