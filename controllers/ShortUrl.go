package controllers

import (
	"ScrabShortener/db"
	"ScrabShortener/helpers"
	"context"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
)

type UrlInfo struct {
	FullUrl      string
	ShortUrl     string
	ClickAmmount int64
}

var Address = helpers.GetEnv("REDIS_URI")
var REDIS_PASS = helpers.GetEnv("REDIS_PASS")

var collection = db.ConnectionString()
var ctx = context.Background()
var threshold = 1000
var rdb = db.RbdConn()

// Client declaration, it  calls  db.Connect that is a connection to the mongodb db

func PostUrlShort(c fiber.Ctx) error {
	// Logic to create a record in MongoDB
	var data map[string]interface{}

	id := xid.New().String()
	shortedUrl := strings.Replace(id[:9], "-", "", -1)

	if !helpers.IsValidUrl(c.FormValue("FullUrl")) {
		data = map[string]interface{}{
			"Message": "Url's format is incorrect",
			"Status":  500,
		}
		return c.Status(500).JSON(data)
	}

	register := UrlInfo{
		FullUrl:      c.FormValue("FullUrl"),
		ShortUrl:     shortedUrl,
		ClickAmmount: 0,
	}

	// Insert the record into the database
	insertResult, err := collection.InsertOne(context.Background(), register)
	if err != nil {
		data = map[string]interface{}{
			"Message": "Error creating record",
			"Error":   err.Error(),
			"Status":  500,
		}
		return c.Status(500).JSON(data)
	}

	err = rdb.HSet(ctx, shortedUrl, map[string]interface{}{
		"full_url":        register.FullUrl,
		"short_url":       shortedUrl,
		"view_count":      0,
		"prev_view_count": 0, // Set previous view count initially to the same as the current view count
	}).Err()
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = rdb.Expire(ctx, shortedUrl, 100*time.Minute).Err()
	if err != nil {
		log.Fatal(err)
		return err
	}

	data = map[string]interface{}{
		"Message":     "URL Shortened successfully",
		"short_url":   shortedUrl,
		"full_url":    c.FormValue("FullUrl"),
		"document_id": insertResult.InsertedID, // Get the ID of the inserted document
		"Status":      200,
	}
	return c.Status(200).JSON(data)
}

func GetShortUrls(c fiber.Ctx) error {

	var data = make(map[string]interface{})

	filter := bson.M{}

	retrievedData, err := collection.Find(ctx, filter)
	if err != nil {
		data = map[string]interface{}{
			"Message": "Error retrieving data",
			"Error":   err.Error(),
			"Status":  fiber.StatusBadRequest,
		}
		return c.Status(fiber.StatusBadRequest).JSON(data)
	}
	defer retrievedData.Close(ctx)

	var shortUrls []UrlInfo
	if err := retrievedData.All(ctx, &shortUrls); err != nil {
		data = map[string]interface{}{
			"Message": "Error decoding data",
			"Error":   err.Error(),
			"Status":  fiber.StatusInternalServerError,
		}
		return c.Status(fiber.StatusInternalServerError).JSON(data)
	}

	responseData := map[string]interface{}{
		"Message": "Data retrieved correctly",
		"Data":    shortUrls,
		"Status":  fiber.StatusOK,
	}

	return c.Status(fiber.StatusOK).JSON(responseData)
}

// Redirection
// Returns a redirection to the URL page. If the page is concurred it does via redis cache.
func RedirectUrl(c fiber.Ctx) error {

	shortedUrl := c.Params("shortUrl")

	filter := bson.M{"shorturl": shortedUrl}
	var urlInfo UrlInfo

	if err := collection.FindOne(ctx, filter).Decode(&urlInfo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{
			"Message": "Error while finding document",
			"Error":   err.Error(),
		})
	}

	urlInfo.ClickAmmount++

	update := bson.M{"$set": bson.M{"clickammount": urlInfo.ClickAmmount}}

	//Update MongoDB count

	if _, err := collection.UpdateOne(ctx, filter, update); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{
			"Message": "Error updating the document",
			"Status":  fiber.StatusBadRequest,
		})
	}

	//Update Redis Count
	err := rdb.HIncrBy(ctx, shortedUrl, "view_count", 1).Err()
	if err != nil {
		log.Fatal("Error updating in redis")
	}

	if int(urlInfo.ClickAmmount) > threshold {
		cachedResponse, err := rdb.HGet(ctx, urlInfo.ShortUrl, "full_url").Result()
		if err != nil {
			log.Fatal("Error redirectioning with redis")
		}
		return c.Redirect().To(cachedResponse)
	}

	return c.Redirect().To(urlInfo.FullUrl)

}
