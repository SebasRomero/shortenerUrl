package url

import (
	"fmt"
	"math"

	"github.com/sebasromero/shortenerUrl/internal/database"
	"github.com/sebasromero/shortenerUrl/internal/types"
)

var COUNTER = 100000000

const base62Digits = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func CreateShortenedUrl(url string) (*types.UrlShortened, error) {
	foundUrl := database.Connect().FindLongUrl(url)
	if foundUrl.LongUrl != "" {
		returnUrl := types.UrlShortened{
			UrlShortened: foundUrl.ShortUrl,
		}
		fmt.Println(returnUrl)
		return &returnUrl, nil
	}

	encode := ConvertToBase62(COUNTER)
	COUNTER++
	shortUrl := types.Path + "/" + encode

	_, err := database.Connect().InsertShortenedUrl(shortUrl, url)
	if err != nil {
		return nil, err
	}

	return &types.UrlShortened{
		UrlShortened: shortUrl,
	}, nil
}

func ConvertToBase62(number int) string {
	base62 := ""
	for number > 0 {
		remainder := number % 62
		base62 = string(base62Digits[remainder]) + base62
		number /= 62
	}
	return base62
}

func ConvertToDecimal(encode string) int {
	result := 0.0
	for i, value_i := range encode {
		for j, value_j := range base62Digits {
			if string(value_i) == string(value_j) {
				value := len(encode) - i - 1
				result += (math.Pow(62, float64(value)) * float64(j))
			}
		}
	}
	return int(result)
}
