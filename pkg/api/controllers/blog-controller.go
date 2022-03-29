package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/yuldashev6267/blog-grpc/internals/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InsertBlog interface {
	InsertBlog(ctx *gin.Context) (*primitive.ObjectID, error)
}

type conn struct {
	mongoDb repository.Database
}

func (db *conn) InsertBlog(ctx *gin.Context)(*primitive.ObjectID, error){
	coll := db.mongoDb.BlogCollection()

	ctx.PostForm("")
}