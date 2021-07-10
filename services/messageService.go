package services

import (
	"example.com/app/domain"
	"example.com/app/repo"
)

type MessageService interface {
	Create(message *domain.Message) error
}

type DefaultMessageService struct {
	repo repo.MessageRepo
}

func (m DefaultMessageService) Create(message *domain.Message) error {
	err := m.repo.Create(message)
	if err != nil {
		return err
	}
	return nil
}

func NewMessageService(repository repo.MessageRepo) DefaultMessageService {
	return DefaultMessageService{repository}
}