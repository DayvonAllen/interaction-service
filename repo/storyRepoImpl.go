package repo

import (
	"context"
	"example.com/app/database"
	"example.com/app/domain"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StoryRepoImpl struct {
	Story             domain.Story
}

func (s StoryRepoImpl) Create(story *domain.Story) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	story.Id = primitive.NewObjectID()

	_, err := conn.StoryCollection.InsertOne(context.TODO(), &story)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	return nil
}

func (s StoryRepoImpl) UpdateByID(story *domain.Story) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	filter := bson.D{{"_id", story.Id}, {"authorUsername", story.AuthorUsername}}
	update := bson.D{{"$set",
		bson.D{{"content", story.Content},
			{"title", story.Title},
		},
	}}

	_, err := conn.StoryCollection.UpdateOne(context.TODO(),
		filter, update)

	if err != nil {
		return fmt.Errorf("you can't update a story you didn't write")
	}

	return nil
}

func (s StoryRepoImpl) DeleteByID(story *domain.Story) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	res, err := conn.StoryCollection.DeleteOne(context.TODO(), bson.D{{"_id", story.Id}, {"authorUsername", story.AuthorUsername}})

	if err != nil {
		panic(err)
	}

	if res.DeletedCount == 0 {
		panic(fmt.Errorf("failed to delete story"))
	}

	return nil
}

func NewStoryRepoImpl() StoryRepoImpl {
	var storyRepoImpl StoryRepoImpl

	return storyRepoImpl
}

