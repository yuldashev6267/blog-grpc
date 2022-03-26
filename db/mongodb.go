package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	Blogs = "blogs"
)

type Database interface {
	Disconnect() error
	BlogCollection() *mongo.Collection
}

type mongoDatabase struct {
	client       *mongo.Client
	database     *mongo.Database
	_collections map[string]*mongo.Collection
}

func (db *mongoDatabase) Disconnect() error {
	return db.client.Disconnect(context.Background())
}

func (db *mongoDatabase) BlogCollection() *mongo.Collection {
	return db.collection(Blogs)
}

func (db *mongoDatabase) collection(name string) *mongo.Collection {
	if k, v := db._collections[name]; v {
		return k
	}

	collections, _ := db.database.ListCollectionNames(context.Background(), bson.M{})

	for _, v := range collections {
		if v == name {
			col := db.database.Collection(name)
			db._collections[name] = col
			return col
		}
	}

	db.database.CreateCollection(context.TODO(), name)
	col := db.database.Collection(name)
	db._collections[name] = col
	return col
}

func DatabaseConn(connectiotStr string, databaseName string) Database {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectiotStr))
	if err != nil {
		return nil
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil
	}

	database := client.Database(databaseName)
	fmt.Println("Database connected")
	return &mongoDatabase{
		database:     database,
		client:       client,
		_collections: make(map[string]*mongo.Collection, 0),
	}
}
