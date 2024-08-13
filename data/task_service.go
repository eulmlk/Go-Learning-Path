package data

import (
	"context"
	"errors"
	"net/http"
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
func NewTaskService(collection *mongo.Collection) *TaskService {
	return &TaskService{
		collection: collection,
	}
}

// GetTasks returns a slice of all tasks.
func (ts *TaskService) GetTasks() ([]models.Task, *models.Error) {
	tasks := []models.Task{}

	// Query the database for all tasks.
	cursor, err := ts.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, &models.Error{
			Err:        err,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}

	// Iterate over the cursor and decode each task into a Task struct.
	for cursor.Next(context.Background()) {
		var task models.Task

		err := cursor.Decode(&task)
		if err != nil {
			return nil, &models.Error{
				Err:        err,
				StatusCode: http.StatusInternalServerError,
				Message:    "Internal server error",
			}
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

// GetTaskByID returns the task with the given ID.
func (ts *TaskService) GetTaskByID(id string, role string,
	userid primitive.ObjectID) (*models.Task, *models.Error) {
	task := &models.Task{}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, &models.Error{
			Err:        err,
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid ID",
		}
	}

	result := ts.collection.FindOne(context.Background(), bson.M{"_id": objectID})
	if result.Err() == mongo.ErrNoDocuments {
		return nil, &models.Error{
			Err:        err,
			StatusCode: http.StatusNotFound,
			Message:    "Task not found",
		}
	}

	err = result.Decode(task)
	if err != nil {
		return nil, &models.Error{
			Err:        err,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}

	if task.UserID != userid && role != "admin" {
		return nil, &models.Error{
			Err:        errors.New("unauthorized"),
			StatusCode: http.StatusForbidden,
			Message:    "Forbidden",
		}
	}

	return task, nil
}

// CreateTask creates a new task and adds it to the tasks map.
func (ts *TaskService) CreateTask(task *models.Task) *models.Error {
	_, err := ts.collection.InsertOne(context.Background(), task)
	if err != nil {
		return &models.Error{
			Err:        err,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}

	return nil
}

func (ts *TaskService) ReplaceTask(
	id string,
	taskData *models.Task,
	userId primitive.ObjectID,
	role string) (*models.Task, *models.Error) {

	_, _err := ts.GetTaskByID(id, role, userId)
	if _err != nil {
		return nil, _err
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, &models.Error{
			Err:        err,
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid ID",
		}
	}

	// Replace the task in the collection.
	taskData.ID = objectID
	taskData.UserID = userId
	result := ts.collection.FindOneAndReplace(context.Background(), bson.M{"_id": objectID}, taskData)
	if result.Err() == mongo.ErrNoDocuments {
		return nil, &models.Error{
			Err:        err,
			StatusCode: http.StatusNotFound,
			Message:    "Task not found",
		}
	}

	if err := result.Decode(taskData); err != nil {
		return nil, &models.Error{
			Err:        err,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}

	return taskData, nil
}

// UpdateTask updates the task with the given ID using the provided task data.
// UpdateTask updates the task with the given ID using the provided task data.
func (ts *TaskService) UpdateTask(id string, taskData *models.Task,
	userId primitive.ObjectID, role string) (*models.Task, *models.Error) {

	_, _err := ts.GetTaskByID(id, role, userId)
	if _err != nil {
		return nil, _err
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, &models.Error{
			Err:        err,
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid ID",
		}
	}

	// Update the task in the collection.
	update := bson.M{
		"id":      objectID,
		"user_id": userId,
	}

	if taskData.Title != "" {
		update["title"] = taskData.Title
	}

	if taskData.Description != "" {
		update["description"] = taskData.Description
	}

	if taskData.Status != "" {
		update["status"] = taskData.Status
	}

	if !taskData.DueDate.IsZero() {
		update["due_date"] = taskData.DueDate
	}

	result, err := ts.collection.UpdateOne(context.Background(), bson.M{"_id": objectID}, bson.M{"$set": update})
	if err != nil {
		return nil, &models.Error{
			Err:        err,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}

	if result.MatchedCount == 0 {
		return nil, &models.Error{
			Err:        errors.New("result MatchedCount is 0"),
			StatusCode: http.StatusNotFound,
			Message:    "Task not found",
		}
	}

	// Fetch the updated task from the collection.
	updatedTask := &models.Task{}
	err = ts.collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(updatedTask)
	if err != nil {
		return nil, &models.Error{
			Err:        err,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}

	return updatedTask, nil
}

// DeleteTask deletes the task with the given ID.
func (ts *TaskService) DeleteTask(id string, role string, userid primitive.ObjectID) *models.Error {
	_, _err := ts.GetTaskByID(id, role, userid)
	if _err != nil {
		return _err
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &models.Error{
			Err:        err,
			StatusCode: http.StatusNotFound,
			Message:    "Invalid ID",
		}
	}

	// Delete the task from the collection.
	result, err := ts.collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		return &models.Error{
			Err:        err,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}

	if result.DeletedCount == 0 {
		return &models.Error{
			Err:        errors.New("result's DeletedCount is 0"),
			StatusCode: http.StatusNotFound,
			Message:    "Task not found",
		}
	}

	return nil
}
