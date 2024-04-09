package types

type ShortUrlResponse struct {
	ShortUrl string `json:"shortUrl"`
}

type LongUrlResponse struct {
	LongUrl string `json:"longUrl"`
}

type FoundUrlResponse struct {
	ShortUrl string `json:"shortUrl"`
	LongUrl  string `json:"longUrl"`
}

type FoundEncodeResponse struct {
	Encode string `json:"encode"`
}

type FoundLastOne struct {
	Encode string `json:"encode"`
}

type InsertUrl struct {
	ShortUrl string `json:"shortUrl"`
	LongUrl  string `json:"longUrl"`
	Encode   string `json:"encode"`
}

var Path string = "http://localhost:8080"
