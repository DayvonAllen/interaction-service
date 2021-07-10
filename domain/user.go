package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id                   primitive.ObjectID   `bson:"_id" json:"id"`
	Email                string               `bson:"email" json:"email"`
	Username             string               `bson:"username" json:"username"`
	CurrentTagLine       string               `bson:"currentTagLine" json:"currentTagLine"`
	ProfilePictureUrl    string               `bson:"profilePictureUrl" json:"profilePictureUrl"`
	ProfileIsViewable    bool                 `bson:"profileIsViewable" json:"profileIsViewable"`
	AcceptMessages       bool                 `bson:"acceptMessages" json:"acceptMessages"`
	IsVerified           bool                 `bson:"isVerified" json:"-"`
}