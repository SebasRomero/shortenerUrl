package dabatase

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sebasromero/shortenerUrl/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DB struct {
	client *mongo.Client
}

var connectionString string = os.Getenv("SRV_MONGO")

func (db *DB) urlShortenerCollection() *mongo.Collection {
	return db.client.Database("url-shortener").Collection("urlShortener")
}

func Connect() *DB {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	return &DB{
		client: client,
	}
}

func (db *DB) GetUrlShortened(url string) *types.ShortUrlResponse {
	urlShortenerCollection := db.urlShortenerCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var urlShortened types.ShortUrlResponse
	filter := bson.D{{"shortUrl", url}}

	err := urlShortenerCollection.FindOne(ctx, filter).Decode(&urlShortened)
	if err != nil {
		fmt.Println(err)
	}
	return &urlShortened
}

func (db *DB) CreateShortenerUrl(url string) *types.UrlShortened {
	urlShortenerCollection := db.urlShortenerCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	shortUrl := types.Path + "/randomCode"

	_, err := urlShortenerCollection.InsertOne(ctx, bson.M{
		"shortUrl": shortUrl,
		"longUrl":  url,
		"clicked":  0,
	})
	if err != nil {
		log.Fatal(err)
	}

	returnShortUrlResponse := types.UrlShortened{
		UrlShortened: shortUrl,
	}

	return &returnShortUrlResponse
}
