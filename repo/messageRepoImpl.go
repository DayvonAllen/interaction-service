package repo

import (
	"context"
	"example.com/app/database"
	"example.com/app/domain"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MessageRepoImpl struct {
	Message             domain.Message
	Conversation        domain.Conversation
}

func (m MessageRepoImpl) Create(message *domain.Message) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	message.Id = primitive.NewObjectID()

	_, err := conn.MessageCollection.InsertOne(context.TODO(), &message)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	err = conn.MessageCollection.FindOne(context.TODO(),
		bson.D{{"_id", message.Id}}).Decode(&m.Message)

	if err != nil {
		return nil
	}

	conversation, err := ConversationRepoImpl{}.FindByOwner(m.Message)

	if err != nil  {
		if err == mongo.ErrNoDocuments {
			err := ConversationRepoImpl{}.Create(m.Message)
			fmt.Println(m.Message)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	conversation.Messages = append(conversation.Messages, m.Message)

	err = ConversationRepoImpl{}.UpdateConversation(*conversation)

	if err != nil {
		return err
	}

	return nil
}

func NewMessageRepoImpl() MessageRepoImpl {
	var messageRepoImpl MessageRepoImpl

	return messageRepoImpl
}