package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Blog struct {
	Id       primitive.ObjectID `json:"_id",bson:"_id",omitempty`
	AuthorId string             `json:"author_id",bson:"author_id"`
	Title    string             `json:"title",bson:"title"`
	Comment  string             `json:"comment",bson:"comment"`
}
