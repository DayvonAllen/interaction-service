package repo

import (
	"example.com/app/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageRepo interface {
	Create(message *domain.Message) error
	DeleteByID(owner string, id primitive.ObjectID) error
	DeleteAllByIDs(owner string, messages []domain.DeleteMessage) error
}

