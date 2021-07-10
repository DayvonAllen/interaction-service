package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Conversation struct {
	Id        primitive.ObjectID `bson:"_id" json:"-"`
	Owner     string `bson:"owner" json:"-"`
	From      string             `json:"-" bson:"from"`
	To        string             `json:"to" bson:"to"`
	Messages  []Message `bson:"messages" json:"messages"`
	CreatedAt time.Time          `bson:"createdAt" json:"-"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"-"`
}
