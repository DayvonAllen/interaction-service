package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Message struct {
	Id        primitive.ObjectID `bson:"_id" json:"-"`
	Content   string             `bson:"content" json:"content"`
	From      string             `json:"-" bson:"from"`
	To        string             `json:"to" bson:"to"`
	CreatedAt time.Time          `bson:"createdAt" json:"sentAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"-"`
}
