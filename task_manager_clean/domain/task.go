package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	TaskCollection = "tasks"
)

// A struct that defines the task model.
type Task struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	DueDate     time.Time          `json:"due_date"`
	Status      string             `json:"status"`
	UserID      primitive.ObjectID `json:"user_id"`
}

// A struct that defines the data required to create a task.
type CreateTaskData struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date" binding:"required"`
	Status      string    `json:"status"`
}

// A struct that defines the data required to fully update a task.
type ReplaceTaskData struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	DueDate     time.Time `json:"due_date" binding:"required"`
	Status      string    `json:"status" binding:"required"`
}

// A struct that defines the data required to partially update a task.
type UpdateTaskData struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status"`
}

// A struct that defines the data that is returned when a task is manipulated.
type TaskView struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status"`
}

// TaskRepository defines the interface for task repository operations.
type TaskRepository interface {
	GetAllTasks() ([]Task, error)
	GetTaskByID(id primitive.ObjectID) (*Task, error)
	AddTask(task *Task) error
	ReplaceTask(id primitive.ObjectID, taskData *Task) error
	UpdateTask(id primitive.ObjectID, taskData bson.M) error
	DeleteTask(id primitive.ObjectID) error
}
