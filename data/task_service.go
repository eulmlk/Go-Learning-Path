package data

import (
	"context"
	"errors"
	"task_manager/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// TaskService is a struct that holds a map of tasks.
type TaskService struct {
	collection *mongo.Collection
}

// NewTaskService creates a new instance of TaskService.
// It initializes the tasks map and returns a pointer to the TaskService.
func NewTaskService(db *mongo.Database) *TaskService {
	return &TaskService{
		collection: db.Collection("tasks"),
	}
}

// GetTasks returns a slice of all tasks.
func (ts *TaskService) GetTasks() ([]models.Task, *models.Error) {
	tasks := []models.Task{}

	// Query the database for all tasks.
	cursor, err := ts.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, &models.Error{Err: err, StatusCode: 500}
	}

	// Iterate over the cursor and decode each task into a Task struct.
	for cursor.Next(context.Background()) {
		var task models.Task

		err := cursor.Decode(&task)
		if err != nil {
			return nil, &models.Error{Err: err, StatusCode: 500}
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

// GetTaskByID returns the task with the given ID.
func (ts *TaskService) GetTaskByID(id string) (*models.Task, *models.Error) {
	task := &models.Task{}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, &models.Error{Err: errors.New("invalid id"), StatusCode: 400}
	}

	result := ts.collection.FindOne(context.Background(), bson.M{"_id": objectID})
	if result.Err() == mongo.ErrNoDocuments {
		return nil, &models.Error{Err: errors.New("task not found"), StatusCode: 404}
	}

	err = result.Decode(task)
	if err != nil {
		return nil, &models.Error{Err: err, StatusCode: 500}
	}

	return task, nil
}

// CreateTask creates a new task and adds it to the tasks map.
func (ts *TaskService) CreateTask(task *models.Task) *models.Error {
	_, err := ts.collection.InsertOne(context.Background(), task)
	if err != nil {
		return &models.Error{Err: err, StatusCode: 500}
	}

	return nil
}

// UpdateTask updates the task with the given ID using the provided task data.
func (ts *TaskService) UpdateTask(id string, taskData *models.Task) (*models.Task, *models.Error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, &models.Error{Err: errors.New("invalid id"), StatusCode: 400}
	}

	// Update the task in the collection.
	result, err := ts.collection.UpdateOne(context.Background(), bson.M{"_id": objectID}, bson.M{"$set": bson.M{
		"title":       taskData.Title,
		"description": taskData.Description,
		"due_date":    taskData.DueDate,
		"status":      taskData.Status,
	}})

	if err != nil {
		return nil, &models.Error{Err: err, StatusCode: 500}
	}

	if result.MatchedCount == 0 {
		return nil, &models.Error{Err: errors.New("task not found"), StatusCode: 404}
	}

	taskData.ID = objectID
	return taskData, nil
}

// DeleteTask deletes the task with the given ID.
func (ts *TaskService) DeleteTask(id string) *models.Error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &models.Error{Err: errors.New("invalid id"), StatusCode: 400}
	}

	// Delete the task from the collection.
	result, err := ts.collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		return &models.Error{Err: err, StatusCode: 500}
	}

	if result.DeletedCount == 0 {
		return &models.Error{Err: errors.New("task not found"), StatusCode: 404}
	}

	return nil
}
