package repository_test

import (
	"context"
	"task_manager/database"
	"task_manager/domain"
	"task_manager/mocks"
	"task_manager/repository"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	taskTestCollection = "tasks_test"
)

type MongoTaskRepositorySuite struct {
	suite.Suite
	repo       *repository.MongoTaskRepository
	collection *mongo.Collection
	client     *mongo.Client
}

func (suite *MongoTaskRepositorySuite) SetupSuite() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	suite.NoError(err)
	suite.NoError(client.Ping(context.Background(), nil))
	suite.client = client

	db := client.Database(database.DatabaseName)
	suite.collection = db.Collection(taskTestCollection)
	suite.NoError(suite.collection.Drop(context.Background()))

	suite.repo = repository.NewMongoTaskRepository(suite.collection)
}

func (suite *MongoTaskRepositorySuite) SetupTest() {
	suite.NoError(suite.collection.Drop(context.Background()))
}

func (suite *MongoTaskRepositorySuite) TearDownSuite() {
	suite.NoError(suite.client.Disconnect(context.Background()))
}

func (suite *MongoTaskRepositorySuite) TestGetAllTasks_Empty() {
	tasks, err := suite.repo.GetAllTasks()
	suite.NoError(err)
	suite.Len(tasks, 0)
}

func (suite *MongoTaskRepositorySuite) TestGetAllTasks_NotEmpty() {
	task1 := mocks.GetNewTask()
	task2 := mocks.GetNewTask2()

	_, err := suite.collection.InsertMany(context.Background(), []interface{}{task1, task2})
	suite.NoError(err)

	tasks, err := suite.repo.GetAllTasks()
	suite.NoError(err)
	suite.Len(tasks, 2)
}

func (suite *MongoTaskRepositorySuite) TestGetTaskByID_NotFound() {
	task := mocks.GetNewTask()

	foundTask, err := suite.repo.GetTaskByID(task.ID)
	suite.Error(err)

	suite.Equal(domain.Task{}, *foundTask)
}

func (suite *MongoTaskRepositorySuite) TestGetTaskByID_Found() {
	task := mocks.GetNewTask()

	_, err := suite.collection.InsertOne(context.Background(), task)
	suite.NoError(err)

	foundTask, err := suite.repo.GetTaskByID(task.ID)
	suite.NoError(err)

	suite.Equal(*task, *foundTask)
}

func (suite *MongoTaskRepositorySuite) TestCreateTask() {
	taskData := mocks.GetNewTask()

	count, err := suite.collection.CountDocuments(context.Background(), bson.M{})
	suite.NoError(err)
	suite.Equal(int64(0), count)

	err = suite.repo.AddTask(taskData)
	suite.NoError(err)

	count, err = suite.collection.CountDocuments(context.Background(), bson.M{})
	suite.NoError(err)
	suite.Equal(int64(1), count)

	foundTask := &domain.Task{}
	result := suite.collection.FindOne(context.Background(), taskData)
	suite.NoError(result.Err())
	suite.NoError(result.Decode(foundTask))
	suite.Equal(taskData, foundTask)
}

func (suite *MongoTaskRepositorySuite) TestReplaceTask_NotFound() {
	taskData := mocks.GetNewTask()

	err := suite.repo.ReplaceTask(taskData.ID, taskData)
	suite.Error(err)
}

func (suite *MongoTaskRepositorySuite) TestReplaceTask_Found() {
	taskData := mocks.GetNewTask()

	_, err := suite.collection.InsertOne(context.Background(), taskData)
	suite.NoError(err)

	newTaskData := mocks.GetNewTask2()
	newTaskData.ID = taskData.ID
	err = suite.repo.ReplaceTask(taskData.ID, newTaskData)
	suite.NoError(err)

	taskFromDB := &domain.Task{}
	result := suite.collection.FindOne(context.Background(), bson.M{"_id": taskData.ID})
	suite.NoError(result.Decode(taskFromDB))
	suite.Equal(*newTaskData, *taskFromDB)
}

func (suite *MongoTaskRepositorySuite) TestUpdateTask_Found() {
	taskData := mocks.GetNewTask()

	_, err := suite.collection.InsertOne(context.Background(), taskData)
	suite.NoError(err)

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

	err = suite.repo.UpdateTask(taskData.ID, taskUpdate)
	suite.NoError(err)

	taskFromDB := &domain.Task{}
	err = suite.collection.FindOne(context.Background(), bson.M{"_id": taskData.ID}).Decode(taskFromDB)
	suite.NoError(err)
	suite.Equal(*expectedTask, *taskFromDB)
}

func (suite *MongoTaskRepositorySuite) TestUpdateTask_Second() {
	taskData := mocks.GetNewTask()

	_, err := suite.collection.InsertOne(context.Background(), taskData)
	suite.NoError(err)

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

	err = suite.repo.UpdateTask(taskData.ID, taskUpdate)
	suite.NoError(err)

	taskFromDB := &domain.Task{}
	err = suite.collection.FindOne(context.Background(), bson.M{"_id": taskData.ID}).Decode(taskFromDB)
	suite.NoError(err)
	suite.Equal(*expectedTask, *taskFromDB)
}

func (suite *MongoTaskRepositorySuite) TestUpdateTask_ID() {
	taskData := mocks.GetNewTask()

	_, err := suite.collection.InsertOne(context.Background(), taskData)
	suite.NoError(err)

	taskUpdate := bson.M{
		"_id": primitive.NewObjectID(),
	}

	err = suite.repo.UpdateTask(taskData.ID, taskUpdate)
	suite.Error(err)
}

func (suite *MongoTaskRepositorySuite) TestDeleteTask() {
	taskData := mocks.GetNewTask()

	_, err := suite.collection.InsertOne(context.Background(), taskData)
	suite.NoError(err)

	err = suite.repo.DeleteTask(taskData.ID)
	suite.NoError(err)

	_, err = suite.repo.GetTaskByID(taskData.ID)
	suite.Error(err)
}

func TestMongoTaskRepositorySuite(t *testing.T) {
	suite.Run(t, new(MongoTaskRepositorySuite))
}
