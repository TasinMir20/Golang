package main

import (
	"log"
	"new-project-go/config"
	"new-project-go/router"
	"os"

	"github.com/gofiber/fiber/v3"

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

	router.SetupRouter(app)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "4000"
	}

	if err := app.Listen(":" + PORT); err != nil {
		panic(err)
	}
}
