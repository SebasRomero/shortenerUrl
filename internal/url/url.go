package url

import (
	"fmt"
	"math"
	"time"

	"github.com/sebasromero/shortenerUrl/internal/cache"
	"github.com/sebasromero/shortenerUrl/internal/database"
	"github.com/sebasromero/shortenerUrl/internal/types"
)

var COUNTER = 100000000

const base62Digits = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func CreateShortenedUrl(url string) (*types.ShortUrlResponse, error) {
	cacheResponse, bool := cache.CreateUrlCache.Get(url)

	if bool {
		return &types.ShortUrlResponse{
			ShortUrl: cacheResponse.ShortUrl,
		}, nil
	}

	foundUrl := database.Connection.FindLongUrl(url)
	if foundUrl.LongUrl != "" {
		returnUrl := types.ShortUrlResponse{
			ShortUrl: foundUrl.ShortUrl,
		}
		fmt.Println(returnUrl)
		return &returnUrl, nil
	}

	encode := ConvertToBase62(COUNTER)
	foundEncode := database.Connection.FindEncode(encode).Encode
	if foundEncode != "" {
		lastOne := database.Connection.FindLastShortedUrl()
		decimal := ConvertToDecimal(lastOne)
		encode = ConvertToBase62(decimal + 1)
	} else {
		COUNTER++
	}

	shortUrl := types.Path + "/" + encode

	insertUrl := types.InsertUrl{
		Encode:   encode,
		ShortUrl: shortUrl,
		LongUrl:  url,
	}

	_, err := database.Connection.InsertShortenedUrl(insertUrl)
	if err != nil {
		return nil, err
	}

	cache.CreateUrlCache.Set(url, insertUrl, time.Minute*2)
	cache.GetUrlCache.Set(shortUrl, types.LongUrlResponse{
		LongUrl: url,
	}, time.Minute*5)

	return &types.ShortUrlResponse{
		ShortUrl: shortUrl,
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
