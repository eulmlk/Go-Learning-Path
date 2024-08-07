package repository

import (
	"task_manager/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TaskRepository defines the interface for task repository operations.
type TaskRepository interface {
	GetAllTasks() ([]domain.Task, error)
	GetTaskByID(id primitive.ObjectID) (*domain.Task, error)
	AddTask(task *domain.Task) error
	ReplaceTask(id primitive.ObjectID, taskData *domain.Task) (*domain.Task, error)
	UpdateTask(id primitive.ObjectID, taskData bson.M) (*domain.Task, error)
	DeleteTask(id primitive.ObjectID) error
}
