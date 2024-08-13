package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TaskRepository defines the interface for task repository operations.
type TaskRepository interface {
	GetAllTasks() ([]Task, error)
	GetTaskByID(id primitive.ObjectID) (*Task, error)
	AddTask(task *Task) error
	ReplaceTask(id primitive.ObjectID, taskData *Task) error
	UpdateTask(id primitive.ObjectID, taskData bson.M) error
	DeleteTask(id primitive.ObjectID) error
}

// UserRepository defines the interface for user repository operations.
type UserRepository interface {
	AddUser(user *User) error
	GetUsers() ([]User, error)
	GetUserByID(objectID primitive.ObjectID) (*User, error)
	GetUserByUsername(username string) (*User, error)
	UpdateUser(objectID primitive.ObjectID, userData bson.M) error
	DeleteUser(objectID primitive.ObjectID) error
}

// TaskUsecase defines the interface for task usecase operations.
type TaskUsecase interface {
	GetTasks() ([]Task, *Error)
	GetTaskByID(objectID primitive.ObjectID) (*Task, *Error)
	CreateTask(taskData *CreateTaskData, claims *Claims) (*TaskView, *Error)
	ReplaceTask(objectID primitive.ObjectID, taskData *ReplaceTaskData, claims *Claims) (*TaskView, *Error)
	UpdateTask(objectID primitive.ObjectID, taskData *UpdateTaskData, claims *Claims) (*TaskView, *Error)
	DeleteTask(objectID primitive.ObjectID, claims *Claims) *Error
}

// UserUsecase defines the interface for user usecase operations.
type UserUsecase interface {
	AddUser(userData *CreateUserData, claims *Claims) (*User, *Error)
	RegisterUser(userData *AuthUserData) (*User, *Error)
	LoginUser(userData *AuthUserData) (string, *Error)
	GetUsers() ([]User, *Error)
	GetUserByID(objectID primitive.ObjectID) (*User, *Error)
	UpdateUser(objectID primitive.ObjectID, userData *UpdateUserData, claims *Claims) (*User, *Error)
	DeleteUser(objectID primitive.ObjectID, claims *Claims) *Error
}

// Collection defines the interface for MongoDB collection operations.
type Collection interface {
	FindOne(context.Context, interface{}, ...*options.FindOneOptions) SingleResult
	InsertOne(context.Context, interface{}, ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	InsertMany(context.Context, []interface{}, ...*options.InsertManyOptions) (*mongo.InsertManyResult, error)
	DeleteOne(context.Context, interface{}, ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	Find(context.Context, interface{}, ...*options.FindOptions) (Cursor, error)
	FindOneAndReplace(context.Context, interface{}, interface{}, ...*options.FindOneAndReplaceOptions) SingleResult
	CountDocuments(context.Context, interface{}, ...*options.CountOptions) (int64, error)
	UpdateOne(context.Context, interface{}, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	UpdateMany(context.Context, interface{}, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error)
}

// Cursor defines the interface for MongoDB cursor operations.
type Cursor interface {
	All(context.Context, interface{}) error
	Next(context.Context) bool
	Decode(interface{}) error
	Close(context.Context) error
}

// SingleResult defines the interface for MongoDB single result operations.
type SingleResult interface {
	Decode(interface{}) error
	Err() error
}
