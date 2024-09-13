package server

import (
	update "ScrabShortener/Update"
	"ScrabShortener/controllers"
	"ScrabShortener/db"
	"ScrabShortener/helpers"
	"context"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func Server() {

	defer func() {
		db.Connect().Disconnect(context.TODO())
	}()

	app := fiber.New()

	port := helpers.GetEnv("PORT")

	if port == "" {
		port = "5000"
	}

	go update.TickerFunc()

	//Middleware
	db.TestDb()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://SCShort.dev", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))
	//Routes
	app.Post("/PostShortUrl", controllers.PostUrlShort) //Creation of URLS
	app.Get("/GetShortUrls", controllers.GetShortUrls)  // Getting the URLS from the db
	app.Get("/:shortUrl", controllers.RedirectUrl)      // Redirect to original links

	log.Fatal(app.Listen(":" + port))

}
