package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	TaskCollection = "tasks"
)

// A struct that defines the task model.
type Task struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	DueDate     time.Time          `json:"due_date" bson:"due_date"`
	Status      string             `json:"status" bson:"status"`
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`
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
