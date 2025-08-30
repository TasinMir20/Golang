package config

import (
	"context"
	"log"
	"os"

	"github.com/qiniu/qmgo"
)

var (
	Client   *qmgo.Client
	Database *qmgo.Database
)

func ConnectDB() {
	var err error

	DB_URI := os.Getenv("DB_URI")
	// Connect to MongoDB
	Client, err = qmgo.NewClient(context.Background(), &qmgo.Config{
		Uri: DB_URI,
	})
	if err != nil {
		log.Fatal("❌ Database connection failure", err)
	}

	// Get database and collections
	Database = Client.Database("go")

	log.Println("✅ Database connection successful.")
}

func DisconnectDB() {
	if Client != nil {
		Client.Close(context.Background())
		log.Println("MongoDB connection closed")
	}
}
