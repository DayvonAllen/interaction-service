package services

import (
	"example.com/app/domain"
	"example.com/app/repo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageService interface {
	Create(message *domain.Message) (*domain.Conversation, error)
	DeleteByID(owner string, id primitive.ObjectID) error
	DeleteAllByIDs(owner string, messageIds []domain.DeleteMessage) error
}

type DefaultMessageService struct {
	repo repo.MessageRepo
}

func (m DefaultMessageService) Create(message *domain.Message) (*domain.Conversation, error) {
	conversation, err := m.repo.Create(message)
	if err != nil {
		return nil, err
	}
	return conversation, nil
}

func (m DefaultMessageService) DeleteByID(owner string, id primitive.ObjectID) error {
	err := m.repo.DeleteByID(owner, id)
	if err != nil {
		return err
	}
	return nil
}

func (m DefaultMessageService) DeleteAllByIDs(owner string, messageIds []domain.DeleteMessage) error {
	err := m.repo.DeleteAllByIDs(owner, messageIds)
	if err != nil {
		return err
	}
	return nil
}

func NewMessageService(repository repo.MessageRepo) DefaultMessageService {
	return DefaultMessageService{repository}
}