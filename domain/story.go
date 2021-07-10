package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Story todo validate struct
type Story struct {
	Id             primitive.ObjectID `bson:"_id" json:"id"`
	Title          string             `bson:"title" json:"title"`
	Content        string             `bson:"content" json:"content"`
	Preview        string             `bson:"preview" json:"preview"`
	AuthorUsername string             `bson:"authorUsername" json:"authorUsername"`
	LikeCount      int                `bson:"likeCount" json:"likeCount"`
	DislikeCount   int                `bson:"dislikeCount" json:"dislikeCount"`
}
