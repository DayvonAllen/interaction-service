package services

import (
	"example.com/app/domain"
	"example.com/app/repo"
)

type ConversationService interface {
	FindConversation(owner, to string) (*domain.Conversation, error)
}

type DefaultConversationService struct {
	repo repo.ConversationRepo
}

func (c DefaultConversationService) FindConversation(owner, to string) (*domain.Conversation, error) {
	conversation, err := c.repo.FindConversation(owner, to)
	if err != nil {
		return nil, err
	}
	return conversation, nil
}

func NewConversationService(repository repo.ConversationRepo) DefaultConversationService {
	return DefaultConversationService{repository}
}