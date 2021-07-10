package repo

import "example.com/app/domain"

type MessageRepo interface {
	Create(message *domain.Message) error
	//DeleteByID(message *domain.Message) error
}

