package server

import (
	"ScrabShortener/controllers"
	"ScrabShortener/db"
	"ScrabShortener/helpers"
	"context"
	"log"

	"github.com/gofiber/fiber/v3"
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

	go helpers.TickerFunc()

	db.TestDb()

	app.Post("/PostShortUrl", controllers.PostUrlShort)
	app.Get("/GetShortUrls", controllers.GetShortUrls)
	app.Get("/:shortUrl", controllers.RedirectUrl)

	log.Fatal(app.Listen(":" + port))

}
