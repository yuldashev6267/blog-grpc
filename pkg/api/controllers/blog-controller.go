package controllers

import (
	"context"
	"fmt"
	"github.com/yuldashev6267/blog-grpc/internals/repository"
	"github.com/yuldashev6267/blog-grpc/pkg/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogService interface {
	InsertBlog(blog *models.Blog) (*primitive.ObjectID, error)
}

type conn struct {
	mongoDb repository.Database
}

func (db *conn) InsertBlog(blog *models.Blog)(*primitive.ObjectID, error){
	coll := db.mongoDb.BlogCollection()

	res,err:=coll.InsertOne(context.Background(), blog)

	if err != nil {
		return  &primitive.NilObjectID,fmt.Errorf("Error has been occured inserting documnet %v", err)
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)

	if !ok  {
		return &primitive.NilObjectID, fmt.Errorf("Error mongo id %v", ok)
	}

	return &oid, nil
}

func New(db repository.Database) BlogService {
	return &conn{mongoDb: db}
}
