package domain

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	UserCollection = "users"
)

// A struct that defines the user model.
type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Username string             `json:"username"`
	Password string             `json:"password"`
	Role     string             `json:"role"`
}

// A struct that defines the data required to register/login a user.
type AuthUserData struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// A struct that defines the data required to create a user.
type CreateUserData struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

// A struct that defines the data required to update a user.
type UpdateUserData struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// UserRepository defines the interface for user repository operations.
type UserRepository interface {
	AddUser(user *User) error
	GetUsers() ([]User, error)
	GetUserByID(objectID primitive.ObjectID) (*User, error)
	GetUserByUsername(username string) (*User, error)
	UpdateUser(objectID primitive.ObjectID, userData bson.M) (*User, error)
	DeleteUser(objectID primitive.ObjectID) error
}
