package repo

import (
	"example.com/app/database"
	"example.com/app/domain"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"time"
)

type ConversationRepoImpl struct {
	Message          domain.Message
	Conversation     domain.Conversation
	Conversation2    domain.Conversation
	ConversationList []domain.Conversation
}

func (c ConversationRepoImpl) Create(message domain.Message) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	c.Conversation.Id = primitive.NewObjectID()
	c.Conversation.CreatedAt = time.Now()
	c.Conversation.Owner = message.From
	c.Conversation.From = message.From
	c.Conversation.To = message.To
	c.Conversation.Messages = append(c.Conversation.Messages, message)
	c.Conversation.UpdatedAt = time.Now()

	c.Conversation2.Id = primitive.NewObjectID()
	c.Conversation2.CreatedAt = time.Now()
	c.Conversation2.Owner = message.To
	c.Conversation2.From = message.From
	c.Conversation2.To = message.To
	c.Conversation2.Messages = append(c.Conversation2.Messages, message)
	c.Conversation2.UpdatedAt = time.Now()

	_, err := conn.ConversationCollection.InsertOne(context.TODO(), c.Conversation)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	_, err = conn.ConversationCollection.InsertOne(context.TODO(), c.Conversation2)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	return nil
}

func (c ConversationRepoImpl) FindByOwner(message domain.Message) (*domain.Conversation, error) {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	fmt.Println(message)
	filter := bson.D{{"owner", message.From}}

	err := conn.ConversationCollection.FindOne(context.TODO(),
		filter).Decode(&c.Conversation)

	if err != nil {
		return nil, err
	}

	return &c.Conversation, nil
}

func (c ConversationRepoImpl) FindConversation(owner, to string) (*domain.Conversation, error) {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	filter := bson.M{
		"owner": owner,
		"$or": []interface{}{
			bson.M{"to": to},
			bson.M{"from": to},
		},
	}

	err := conn.ConversationCollection.FindOne(context.TODO(),
		filter).Decode(&c.Conversation)

	if err != nil {
		return nil, err
	}

	return &c.Conversation, nil
}

func (c ConversationRepoImpl) UpdateConversation(conversation domain.Conversation) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.M{
		"owner": conversation.Owner,
		"$or": []interface{}{
			bson.M{"to": conversation.To},
			bson.M{"from": conversation.To},
		},
	}

	update := bson.D{{"$set", bson.D{{"messages", conversation.Messages}, {"updatedAt", time.Now()}}}}

	err := conn.ConversationCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&c.Conversation)

	if err != nil {
		return err
	}

	opts = options.FindOneAndUpdate().SetUpsert(true)
	filter = bson.M{
		"owner": conversation.To,
		"$or": []interface{}{
			bson.M{"to": conversation.Owner},
			bson.M{"from": conversation.Owner},
		},
	}
	update = bson.D{{"$set", bson.D{{"messages", conversation.Messages}, {"updatedAt", time.Now()}}}}

	err = conn.ConversationCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&c.Conversation)

	if err != nil {
		return err
	}

	return nil
}

func (c ConversationRepoImpl) DeleteByID(conversationId primitive.ObjectID, username string) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	res, err := conn.ConversationCollection.DeleteOne(context.TODO(), bson.D{{"_id", conversationId}, {"owner", username}})

	if err != nil {
		panic(err)
	}

	if res.DeletedCount == 0 {
		panic(fmt.Errorf("failed to delete story"))
	}

	return nil
}

func (c ConversationRepoImpl) DeleteAllByUsername(username string) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	res, err := conn.ConversationCollection.DeleteMany(context.TODO(), bson.D{{"owner", username}})

	if err != nil {
		panic(err)
	}

	if res.DeletedCount == 0 {
		panic(fmt.Errorf("failed to delete story"))
	}

	return nil
}

func NewConversationRepoImpl() ConversationRepoImpl {
	var conversationRepoImpl ConversationRepoImpl

	return conversationRepoImpl
}
