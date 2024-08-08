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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	taskTestCollection = "tasks_test"
)

// A test for the MongoTaskRepository.GetAllTasks method.
func TestMongoTaskRepository_GetAllTasks(t *testing.T) {
	client, collection := setupTaskCollection(t)
	defer client.Disconnect(context.Background())

	repo := repository.NewMongoTaskRepository(collection)
	task1 := getNewTask1()
	task2 := getNewTask2()

	// A test case for an empty collection
	t.Run("GetAllTasks_Empty", func(t *testing.T) {
		tasks, err := repo.GetAllTasks()
		require.NoError(t, err)
		assert.Len(t, tasks, 0)
	})

	// A test case for a non-empty collection
	t.Run("GetAllTasks_NotEmpty", func(t *testing.T) {
		_, err := collection.InsertMany(context.Background(), []interface{}{task1, task2})
		require.NoError(t, err)

		tasks, err := repo.GetAllTasks()
		require.NoError(t, err)
		assert.Len(t, tasks, 2)
	})
}

// A test for the MongoTaskRepository.GetTaskByUserID method.
func TestMongoTaskRepository_GetTaskByID(t *testing.T) {
	client, collection := setupTaskCollection(t)
	defer client.Disconnect(context.Background())

	repo := repository.NewMongoTaskRepository(collection)
	task := getNewTask1()

	// A test case for a task not found
	t.Run("GetTaskByID_NotFound", func(t *testing.T) {
		foundTask, err := repo.GetTaskByID(task.ID)
		assert.Error(t, err)

		assert.Equal(t, *foundTask, domain.Task{})
	})

	// A test case for a task found
	t.Run("GetTaskByID_Found", func(t *testing.T) {
		_, err := collection.InsertOne(context.Background(), task)
		require.NoError(t, err)

		foundTask, err := repo.GetTaskByID(task.ID)
		require.NoError(t, err)

		assert.Equal(t, *task, *foundTask)
	})
}

// A test for the MongoTaskRepository.GetTaskByUserID method.
func TestMongoTaskRepository_CreateTask(t *testing.T) {
	client, collection := setupTaskCollection(t)
	defer client.Disconnect(context.Background())

	repo := repository.NewMongoTaskRepository(collection)
	taskData := getNewTask1()

	// A test case for creating a task
	t.Run("CreateTask", func(t *testing.T) {
		// 1. Check the number of tasks in the collection
		count, err := collection.CountDocuments(context.Background(), bson.M{})
		require.NoError(t, err)
		assert.Equal(t, int64(0), count)

		// 2. Add a new task to the collection
		err = repo.AddTask(taskData)
		require.NoError(t, err)

		// 3. Check the number of tasks in the collection after adding a new task
		count, err = collection.CountDocuments(context.Background(), bson.M{})
		require.NoError(t, err)
		assert.Equal(t, int64(1), count)

		// 4. Check if the task was added correctly
		foundTask := &domain.Task{}
		result := collection.FindOne(context.Background(), taskData)
		require.NoError(t, result.Err())
		require.NoError(t, result.Decode(foundTask))
		assert.Equal(t, taskData, foundTask)
	})
}

// A test for the MongoTaskRepository.GetTaskByUserID method.
func TestMongoTaskRepository_ReplaceTask(t *testing.T) {
	client, collection := setupTaskCollection(t)
	defer client.Disconnect(context.Background())

	repo := repository.NewMongoTaskRepository(collection)
	taskData := getNewTask1()

	// A test case when the task is not found
	t.Run("ReplaceTask_NotFound", func(t *testing.T) {
		err := repo.ReplaceTask(taskData.ID, taskData)
		assert.Error(t, err)
	})

	_, err := collection.InsertOne(context.Background(), taskData)
	require.NoError(t, err)

	// A test case when the task is found
	t.Run("ReplaceTask_Found", func(t *testing.T) {
		newTaskData := getNewTask2()
		newTaskData.ID = taskData.ID
		err := repo.ReplaceTask(taskData.ID, newTaskData)
		require.NoError(t, err)

		taskFromDB := &domain.Task{}
		result := collection.FindOne(context.Background(), bson.M{"_id": taskData.ID})
		require.NoError(t, result.Decode(taskFromDB))
		assert.Equal(t, *newTaskData, *taskFromDB)
	})
}

// A test for the MongoTaskRepository.GetTaskByUserID method.
func TestMongoTaskRepository_UpdateTask(t *testing.T) {
	client, collection := setupTaskCollection(t)
	defer client.Disconnect(context.Background())

	repo := repository.NewMongoTaskRepository(collection)
	taskData := getNewTask1()

	_, err := collection.InsertOne(context.Background(), taskData)
	require.NoError(t, err)

	// A test case when the task is found
	t.Run("UpdateTask_Found", func(t *testing.T) {
		taskUpdate := bson.M{
			"title":  "New Task",
			"status": "Completed",
		}

		expectedTask := &domain.Task{
			ID:          taskData.ID,
			Title:       "New Task",
			Description: taskData.Description,
			DueDate:     taskData.DueDate,
			Status:      "Completed",
			UserID:      taskData.UserID,
		}

		err := repo.UpdateTask(taskData.ID, taskUpdate)
		require.NoError(t, err)

		taskFromDB := &domain.Task{}
		err = collection.FindOne(context.Background(), bson.M{"_id": taskData.ID}).Decode(taskFromDB)
		require.NoError(t, err)
		assert.Equal(t, *expectedTask, *taskFromDB)
	})

	// A second test case for updating a task
	t.Run("UpdateTask_Second", func(t *testing.T) {
		taskUpdate := bson.M{
			"title":       "Another New Task",
			"description": "New Description",
		}

		expectedTask := &domain.Task{
			ID:          taskData.ID,
			Title:       "Another New Task",
			Description: "New Description",
			DueDate:     taskData.DueDate,
			Status:      "Completed",
			UserID:      taskData.UserID,
		}

		err := repo.UpdateTask(taskData.ID, taskUpdate)
		require.NoError(t, err)

		taskFromDB := &domain.Task{}
		err = collection.FindOne(context.Background(), bson.M{"_id": taskData.ID}).Decode(taskFromDB)
		require.NoError(t, err)
		assert.Equal(t, *expectedTask, *taskFromDB)
	})

	// A test case for updating the _id field
	t.Run("UpdateTask_ID", func(t *testing.T) {
		taskUpdate := bson.M{
			"_id": primitive.NewObjectID(),
		}

		err := repo.UpdateTask(taskData.ID, taskUpdate)
		assert.Error(t, err)
	})
}

// A test for the MongoTaskRepository.DeleteTask method.
func TestMongoTaskRepository_DeleteTask(t *testing.T) {
	client, collection := setupTaskCollection(t)
	defer client.Disconnect(context.Background())

	repo := repository.NewMongoTaskRepository(collection)
	taskData := getNewTask1()

	_, err := collection.InsertOne(context.Background(), taskData)
	require.NoError(t, err)

	// A test case for deleting a task
	t.Run("DeleteTask_Found", func(t *testing.T) {
		err = repo.DeleteTask(taskData.ID)
		require.NoError(t, err)

		_, err = repo.GetTaskByID(taskData.ID)
		assert.Error(t, err)
	})
}

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

func getNewTask1() *domain.Task {
	return &domain.Task{
		ID:          primitive.NewObjectID(),
		Title:       "My First Task",
		Description: "This is some example description for the first task.",
		DueDate:     format(time.Now().AddDate(0, 0, 3)),
		Status:      "In Progress",
		UserID:      primitive.NewObjectID(),
	}
}

func getNewTask2() *domain.Task {
	return &domain.Task{
		ID:          primitive.NewObjectID(),
		Title:       "My Second Task",
		Description: "This is some example description for the second task.",
		DueDate:     format(time.Now().AddDate(0, 0, 1)),
		Status:      "Completed",
		UserID:      primitive.NewObjectID(),
	}
}

func format(d time.Time) time.Time {
	return d.UTC().Truncate(time.Millisecond)
}
