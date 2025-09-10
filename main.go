package main

import (
	"log"
	"new-project-go/config"
	"new-project-go/router"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/static"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	// Initialize database
	config.ConnectDB()
	defer config.DisconnectDB()

	app := fiber.New()

	// logging middleware - logs every request to console
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${method} ${path} - ${ip} - ${latency}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	}))

	app.Get("/static/*", static.New("./public"))

	router.SetupRouter(app)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "4000"
	}

	if err := app.Listen(":" + PORT); err != nil {
		panic(err)
	}
}
