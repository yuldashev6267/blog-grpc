package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DatabaseConn() {

	err := godotenv.Load()

	if err != nil {
		log.Fatalf("There is an error reading env file %v", err)
	}

	client, databaseErr := mongo.NewClient(options.Client().ApplyURI(os.Getenv("DATABASE_CONNECTION")))

	if err != nil {
		log.Fatalf("Error connecting database %v", databaseErr)
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	databaseErr = client.Connect(ctx)

	if databaseErr != nil {
		log.Fatalf("Mongo driver does not respond %v", err)
	}

	defer client.Disconnect(ctx)
	fmt.Println("Database Connected")

}
