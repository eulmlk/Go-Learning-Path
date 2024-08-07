package repository

import (
	"task_manager/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	AddUser(user *domain.User) error
	GetUsers() ([]domain.User, error)
	GetUserByID(objectID primitive.ObjectID) (*domain.User, error)
	GetUserByUsername(username string) (*domain.User, error)
	UpdateUser(objectID primitive.ObjectID, userData bson.M) (*domain.User, error)
	DeleteUser(objectID primitive.ObjectID) error
}
