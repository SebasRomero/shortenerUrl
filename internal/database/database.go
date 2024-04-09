package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sebasromero/shortenerUrl/internal/cache"
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

var Connection = Connect()

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

func (db *DB) GetUrlShortened(shortUrl string) (*types.LongUrlResponse, error) {
	response, bool := cache.GetUrlCache.Get(shortUrl)
	if bool {
		return &types.LongUrlResponse{
			LongUrl: response.LongUrl,
		}, nil
	}
	urlShortenerCollection := db.urlShortenerCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var longUrl types.LongUrlResponse
	filter := bson.D{{Key: "shortUrl", Value: shortUrl}}
	update := bson.D{{Key: "$inc", Value: bson.D{{Key: "clicked", Value: 1}}}}

	err := urlShortenerCollection.FindOneAndUpdate(ctx, filter, update).Decode(&longUrl)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	cache.GetUrlCache.Set(shortUrl, longUrl, time.Minute*5)
	return &longUrl, nil
}

func (db *DB) InsertShortenedUrl(insertUrl types.InsertUrl) (*types.ShortUrlResponse, error) {
	urlShortenerCollection := db.urlShortenerCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	now := time.Now()
	expires := time.Now().Add(time.Hour)
	_, err := urlShortenerCollection.InsertOne(ctx, bson.M{
		"shortUrl":  insertUrl.ShortUrl,
		"longUrl":   insertUrl.LongUrl,
		"encode":    insertUrl.Encode,
		"clicked":   0,
		"expiresAt": expires,
		"createdAt": now,
		"updatedAt": now,
	})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	returnShortUrlResponse := types.ShortUrlResponse{
		ShortUrl: insertUrl.ShortUrl,
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

func (db *DB) FindEncode(encode string) *types.FoundEncodeResponse {
	urlShortenerCollection := db.urlShortenerCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var foundEncode types.FoundEncodeResponse
	filter := bson.D{{Key: "encode", Value: encode}}
	err := urlShortenerCollection.FindOne(ctx, filter).Decode(&foundEncode)

	if err != nil {
		fmt.Println(err)
	}
	return &foundEncode
}

func (db *DB) FindLastShortedUrl() string {
	urlShortenerCollection := db.urlShortenerCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var foundLastOne types.FoundLastOne
	sort := options.FindOne().SetSort(bson.D{{Key: "createdAt", Value: -1}})

	err := urlShortenerCollection.FindOne(ctx, bson.D{}, sort).Decode(&foundLastOne)

	if err != nil {
		fmt.Println(err)
	}

	return foundLastOne.Encode
}

func (db *DB) RemoveExpiredShortUrls() error {
	urlShortenerCollection := db.urlShortenerCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var foundUrl types.FoundUrlResponse
	filter := bson.D{{Key: "expiresAt", Value: bson.D{{Key: "$lte", Value: time.Now()}}}}

	err := urlShortenerCollection.FindOneAndDelete(ctx, filter).Decode(&foundUrl)
	if err != nil {
		fmt.Println("Could not remove")
		return err
	}
	cache.CreateUrlCache.Remove(foundUrl.LongUrl)
	cache.GetUrlCache.Remove(foundUrl.ShortUrl)
	fmt.Println("Removed")
	return nil
}
