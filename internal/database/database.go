package database

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

func (db *DB) GetUrlShortened(url string) (*types.ShortUrlResponse, error) {
	urlShortenerCollection := db.urlShortenerCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var urlShortened types.ShortUrlResponse
	filter := bson.D{{Key: "shortUrl", Value: url}}

	err := urlShortenerCollection.FindOne(ctx, filter).Decode(&urlShortened)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &urlShortened, nil
}

func (db *DB) InsertShortenedUrl(shortUrl string, longUrl string) (*types.UrlShortened, error) {
	urlShortenerCollection := db.urlShortenerCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := urlShortenerCollection.InsertOne(ctx, bson.M{
		"shortUrl": shortUrl,
		"longUrl":  longUrl,
		"clicked":  0,
	})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	returnShortUrlResponse := types.UrlShortened{
		UrlShortened: shortUrl,
	}

	return &returnShortUrlResponse, err
}

func (db *DB) FindLongUrl(url string) *types.FoundUrlResponse {
	urlShortenerCollection := db.urlShortenerCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var foundUrl types.FoundUrlResponse
	filter := bson.D{{Key: "longUrl", Value: url}}
	err := urlShortenerCollection.FindOne(ctx, filter).Decode(&foundUrl)

	if err != nil {
		fmt.Println(err)
	}
	return &foundUrl
}
