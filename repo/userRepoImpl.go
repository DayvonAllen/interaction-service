package repo

import (
	"context"
	"example.com/app/database"
	"example.com/app/domain"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"sync"
)

type UserRepoImpl struct {
	user        domain.User
}

func (u UserRepoImpl) Create(user *domain.User) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	cur, err := conn.UserCollection.Find(context.TODO(), bson.M{
		"$or": []interface{}{
			bson.M{"email": user.Email},
			bson.M{"username": user.Username},
		},
	})

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	if !cur.Next(context.TODO()) {
		_, err = conn.UserCollection.InsertOne(context.TODO(), &user)

		if err != nil {
			return fmt.Errorf("error processing data")
		}

		go func() {
			event := new(domain.Event)
			event.Action = "create user"
			event.Target = user.Id.String()
			event.ResourceId = user.Id
			event.ActorUsername = user.Username
			event.Message = "user was created, Id:" + user.Id.String()
			err = SendEventMessage(event, 0)
			if err != nil {
				fmt.Println("Error publishing...")
				return
			}
		}()

		return nil
	}

	return fmt.Errorf("user already exists")
}

func (u UserRepoImpl) UpdateByID(user *domain.User) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", user.Id}}
	update := bson.D{{"$set", user}}

	conn.UserCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts)

	return nil
}

func (u UserRepoImpl) DeleteByID(user *domain.User) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	// sets mongo's read and write concerns
	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	// set up for a transaction
	session, err := conn.StartSession()

	if err != nil {
		panic(err)
	}

	defer session.EndSession(context.Background())

	// execute this code in a logical transaction
	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {
		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			_, err := conn.UserCollection.DeleteOne(context.TODO(), bson.D{{"_id", user.Id}})

			if err != nil {
				panic(err)
			}

			return
		}()

		go func() {
			defer wg.Done()
			_, err = conn.StoryCollection.DeleteMany(context.TODO(), bson.D{{"authorUsername", user.Username}})

			if err != nil {
				panic(err)

			}
			return
		}()

		wg.Wait()

		return nil, err
	}

	_, err = session.WithTransaction(context.Background(), callback, txnOpts)

	if err != nil {
		return fmt.Errorf("failed to delete user")
	}

	go func() {
		event := new(domain.Event)
		event.Action = "delete user"
		event.Target = user.Id.String()
		event.ResourceId = user.Id
		event.ActorUsername = user.Username
		event.Message = "user was deleted, Id:" + user.Id.String()
		err = SendEventMessage(event, 0)
		if err != nil {
			fmt.Println("Error publishing...")
			return
		}
	}()
	return nil
}

func NewUserRepoImpl() UserRepoImpl {
	var userRepoImpl UserRepoImpl

	return userRepoImpl
}

