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

func (m MessageRepoImpl) DeleteByID(owner string, id primitive.ObjectID) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	filter := bson.D{{"owner", owner}, {"messages._id", id}}

	err := conn.ConversationCollection.FindOne(context.TODO(), filter).Decode(&m.Conversation)

	if err != nil {
		return err
	}

	messages := make([]domain.Message,0, len(m.Conversation.Messages))

	for _, v := range m.Conversation.Messages {
		fmt.Println(v)
		if v.Id != id {
			messages = append(messages, v)
		}
	}

	m.Conversation.Messages = messages

	update := bson.M{"$set": bson.M{"messages": messages}}

	_, err  = conn.ConversationCollection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		return err
	}

	return nil
}


func NewMessageRepoImpl() MessageRepoImpl {
	var messageRepoImpl MessageRepoImpl

	return messageRepoImpl
}