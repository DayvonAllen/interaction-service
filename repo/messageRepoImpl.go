package repo

import (
	"context"
	"example.com/app/database"
	"example.com/app/domain"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type MessageRepoImpl struct {
	Message             domain.Message
	MessageList             []domain.Message
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

	err = ConversationRepoImpl{}.UpdateConversation(*conversation, m.Message)

	if err != nil {
		return err
	}

	return nil
}

func (m MessageRepoImpl) DeleteByID(owner string, id primitive.ObjectID) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	filter := bson.D{{"_id", id}}

	err := conn.MessageCollection.FindOne(context.TODO(), filter).Decode(&m.Message)

	if err != nil {
		return err
	}

	filter = bson.D{{"owner", owner}, {"messages._id", id}}

	err = conn.ConversationCollection.FindOne(context.TODO(), filter).Decode(&m.Conversation)

	if err != nil {
		return err
	}

	update := bson.M{"$pull": bson.M{"messages": m.Message}}

	_, err  = conn.ConversationCollection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		return err
	}

	return nil
}

func (m MessageRepoImpl) DeleteAllByIDs(owner string, ids []primitive.ObjectID) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	filter := bson.D{{"owner", owner}, {"_id",bson.D{{"$in", ids}}}}

	cur, err := conn.MessageCollection.Find(context.TODO(), filter)

	if err != nil {
		return err
	}

	if err = cur.All(context.TODO(), &m.Message); err != nil {
		log.Fatal(err)
	}

	update := bson.M{"$pull": bson.M{"messages": m.MessageList}}


	_, err = conn.MessageCollection.UpdateMany(context.TODO(), filter, update)

	if err != nil {
		return err
	}

	return nil
}

func NewMessageRepoImpl() MessageRepoImpl {
	var messageRepoImpl MessageRepoImpl

	return messageRepoImpl
}