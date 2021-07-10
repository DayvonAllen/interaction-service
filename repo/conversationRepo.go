package repo

import (
	"example.com/app/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConversationRepo interface {
	Create(message domain.Message) error
	FindByOwner(message domain.Message) (*domain.Conversation, error)
	FindConversation(owner, to string) (*domain.Conversation, error)
	UpdateConversation(conversation domain.Conversation, message domain.Message) error
	DeleteByID(conversationId primitive.ObjectID, username string) error
	DeleteAllByUsername(username string) error
}
