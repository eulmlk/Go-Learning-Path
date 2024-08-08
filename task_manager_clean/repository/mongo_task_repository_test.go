package repository_test

import (
	"context"
	"task_manager/database"
	"task_manager/domain"
	"task_manager/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	taskTestCollection = "tasks_test"
)

func setupTaskCollection(t *testing.T) (*mongo.Client, *mongo.Collection) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	require.NoError(t, err)
	require.NoError(t, client.Ping(context.Background(), nil))

	db := client.Database(database.DatabaseName)
	collection := db.Collection(taskTestCollection)
	require.NoError(t, collection.Drop(context.Background()))

	return client, collection
}

func TestMongoTaskRepository_GetAllTasks(t *testing.T) {
	client, collection := setupTaskCollection(t)
	defer client.Disconnect(context.Background())

	repo := repository.NewMongoTaskRepository(collection)

	task1 := &domain.Task{
		ID:          primitive.NewObjectID(),
		Title:       "Task 1",
		Description: "Description 1",
		DueDate:     time.Now().AddDate(0, 0, 3),
		Status:      "In Progress",
		UserID:      primitive.NewObjectID(),
	}

	task2 := &domain.Task{
		ID:          primitive.NewObjectID(),
		Title:       "Task 2",
		Description: "Description 2",
		DueDate:     time.Now().AddDate(0, 0, 1),
		Status:      "Completed",
		UserID:      primitive.NewObjectID(),
	}

	tasks, err := repo.GetAllTasks()
	require.NoError(t, err)
	assert.Len(t, tasks, 0)

	_, err = collection.InsertMany(context.Background(), []interface{}{task1, task2})
	require.NoError(t, err)

	tasks, err = repo.GetAllTasks()
	require.NoError(t, err)
	assert.Len(t, tasks, 2)
}

func TestMongoTaskRepository_GetTaskByID(t *testing.T) {
	client, collection := setupTaskCollection(t)
	defer client.Disconnect(context.Background())

	repo := repository.NewMongoTaskRepository(collection)

	// Set the DueDate and normalize to UTC
	dueDate := time.Now().AddDate(0, 0, 3).UTC().Truncate(time.Millisecond)

	task := &domain.Task{
		ID:          primitive.NewObjectID(),
		Title:       "Task",
		Description: "Description",
		DueDate:     dueDate,
		Status:      "In Progress",
		UserID:      primitive.NewObjectID(),
	}

	_, err := collection.InsertOne(context.Background(), task)
	require.NoError(t, err)

	foundTask, err := repo.GetTaskByID(task.ID)
	require.NoError(t, err)

	// Normalize foundTask.DueDate to UTC and truncate for comparison
	expectedDueDate := task.DueDate.UTC().Truncate(time.Millisecond)
	actualDueDate := foundTask.DueDate.UTC().Truncate(time.Millisecond)

	assert.Equal(t, task.ID, foundTask.ID)
	assert.Equal(t, task.Title, foundTask.Title)
	assert.Equal(t, task.Description, foundTask.Description)
	assert.Equal(t, expectedDueDate, actualDueDate)
	assert.Equal(t, task.Status, foundTask.Status)
	assert.Equal(t, task.UserID, foundTask.UserID)
}

func TestMongoTaskRepository_CreateTask(t *testing.T) {
	client, collection := setupTaskCollection(t)
	defer client.Disconnect(context.Background())

	repo := repository.NewMongoTaskRepository(collection)

	taskData := &domain.Task{
		ID:          primitive.NewObjectID(),
		Title:       "Task",
		Description: "Description",
		DueDate:     time.Now().AddDate(0, 0, 3),
		Status:      "In Progress",
		UserID:      primitive.NewObjectID(),
	}

	err := repo.AddTask(taskData)
	require.NoError(t, err)

	foundTask, err := repo.GetTaskByID(taskData.ID)
	require.NoError(t, err)

	// Normalize foundTask.DueDate to UTC and truncate for comparison
	expectedDueDate := taskData.DueDate.UTC().Truncate(time.Millisecond)
	actualDueDate := foundTask.DueDate.UTC().Truncate(time.Millisecond)

	assert.Equal(t, taskData.ID, foundTask.ID)
	assert.Equal(t, taskData.Title, foundTask.Title)
	assert.Equal(t, taskData.Description, foundTask.Description)
	assert.Equal(t, expectedDueDate, actualDueDate)
	assert.Equal(t, taskData.Status, foundTask.Status)
	assert.Equal(t, taskData.UserID, foundTask.UserID)
}

func TestMongoTaskRepository_ReplaceTask(t *testing.T) {
	client, collection := setupTaskCollection(t)
	defer client.Disconnect(context.Background())

	repo := repository.NewMongoTaskRepository(collection)

	taskData := &domain.Task{
		ID:          primitive.NewObjectID(),
		Title:       "Task",
		Description: "Description",
		DueDate:     time.Now().AddDate(0, 0, 3),
		Status:      "In Progress",
		UserID:      primitive.NewObjectID(),
	}

	err := repo.AddTask(taskData)
	require.NoError(t, err)

	newTaskData := &domain.Task{
		ID:          taskData.ID,
		Title:       "New Task",
		Description: "New Description",
		DueDate:     time.Now().AddDate(0, 0, 1),
		Status:      "Completed",
		UserID:      taskData.UserID,
	}

	updatedTask, err := repo.ReplaceTask(taskData.ID, newTaskData)
	require.NoError(t, err)

	// Normalize updatedTask.DueDate to UTC and truncate for comparison
	expectedDueDate := newTaskData.DueDate.UTC().Truncate(time.Millisecond)
	actualDueDate := updatedTask.DueDate.UTC().Truncate(time.Millisecond)

	assert.Equal(t, newTaskData.ID, updatedTask.ID)
	assert.Equal(t, newTaskData.Title, updatedTask.Title)
	assert.Equal(t, newTaskData.Description, updatedTask.Description)
	assert.Equal(t, expectedDueDate, actualDueDate)
	assert.Equal(t, newTaskData.Status, updatedTask.Status)
	assert.Equal(t, newTaskData.UserID, updatedTask.UserID)
}

func TestMongoTaskRepository_UpdateTask(t *testing.T) {
	client, collection := setupTaskCollection(t)
	defer client.Disconnect(context.Background())

	repo := repository.NewMongoTaskRepository(collection)

	taskData := &domain.Task{
		ID:          primitive.NewObjectID(),
		Title:       "Task",
		Description: "Description",
		DueDate:     time.Now().AddDate(0, 0, 3),
		Status:      "In Progress",
		UserID:      primitive.NewObjectID(),
	}

	err := repo.AddTask(taskData)
	require.NoError(t, err)

	taskUpdate1 := map[string]interface{}{
		"title":  "New Task",
		"status": "Completed",
	}

	updatedTask1, err := repo.UpdateTask(taskData.ID, taskUpdate1)
	require.NoError(t, err)

	// Normalize updatedTask1.DueDate to UTC and truncate for comparison
	expectedDueDate := taskData.DueDate.UTC().Truncate(time.Millisecond)
	actualDueDate := updatedTask1.DueDate.UTC().Truncate(time.Millisecond)

	assert.Equal(t, taskData.ID, updatedTask1.ID)
	assert.Equal(t, "New Task", updatedTask1.Title)
	assert.Equal(t, taskData.Description, updatedTask1.Description)
	assert.Equal(t, expectedDueDate, actualDueDate)
	assert.Equal(t, "Completed", updatedTask1.Status)
	assert.Equal(t, taskData.UserID, updatedTask1.UserID)

	taskUpdate2 := map[string]interface{}{
		"title":       "Another New Task",
		"description": "New Description",
	}

	updatedTask2, err := repo.UpdateTask(taskData.ID, taskUpdate2)
	require.NoError(t, err)

	assert.Equal(t, taskData.ID, updatedTask2.ID)
	assert.Equal(t, "Another New Task", updatedTask2.Title)
	assert.Equal(t, "New Description", updatedTask2.Description)
	assert.Equal(t, "Completed", updatedTask2.Status)
	assert.Equal(t, taskData.UserID, updatedTask2.UserID)
}

func TestMongoTaskRepository_DeleteTask(t *testing.T) {
	client, collection := setupTaskCollection(t)
	defer client.Disconnect(context.Background())

	repo := repository.NewMongoTaskRepository(collection)

	taskData := &domain.Task{
		ID:          primitive.NewObjectID(),
		Title:       "Task",
		Description: "Description",
		DueDate:     time.Now().AddDate(0, 0, 3),
		Status:      "In Progress",
		UserID:      primitive.NewObjectID(),
	}

	err := repo.AddTask(taskData)
	require.NoError(t, err)

	err = repo.DeleteTask(taskData.ID)
	require.NoError(t, err)

	_, err = repo.GetTaskByID(taskData.ID)
	assert.Error(t, err)
}
