package repository_test

import (
	"task_manager/domain"
	"task_manager/mocks"
	"task_manager/repository"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// A suite that contains tests for the MongoTaskRepository.
type MongoTaskRepositoryTestSuite struct {
	suite.Suite
	repo       *repository.MongoTaskRepository
	collection *mocks.Collection
}

// A method that initializes the test suite.
func (suite *MongoTaskRepositoryTestSuite) SetupSuite() {
	suite.collection = new(mocks.Collection)
	suite.repo = repository.NewMongoTaskRepository(suite.collection)
}

// A method that finalizes the test suite.
func (suite *MongoTaskRepositoryTestSuite) TearDownSuite() {
	suite.collection.AssertExpectations(suite.T())
}

// A test for the MongoUserRepository.GetAllTasks method.
func (suite *MongoTaskRepositoryTestSuite) TestGetAllTasks() {
	// A testcase for the successful retrieval of tasks.
	suite.Run("GetAllTasks_Success", func() {
		tasks := mocks.GetManyTasks()

		cursor := new(mocks.Cursor)
		cursor.On("All", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			taskPtr := args.Get(1).(*[]domain.Task)
			*taskPtr = append(*taskPtr, tasks...)
		})

		suite.collection.On("Find", mock.Anything, mock.Anything).Return(cursor, nil).Once()

		result, err := suite.repo.GetAllTasks()
		suite.NoError(err)
		suite.Equal(tasks, result)
	})

	// A testcase for the failure of retrieving tasks.
	suite.Run("GetAllTasks_Failure", func() {
		suite.collection.On("Find", mock.Anything, mock.Anything).Return(new(mocks.Cursor), mongo.ErrNoDocuments).Once()

		result, err := suite.repo.GetAllTasks()
		suite.Error(err)
		suite.Nil(result)
	})
}

// A test for the MongoUserRepository.GetTaskByID method.
func (suite *MongoTaskRepositoryTestSuite) TestGetTaskByID() {
	// A testcase for the successful retrieval of a task.
	suite.Run("GetTaskByID_Success", func() {
		task := mocks.GetNewTask()

		res := new(mocks.SingleResult)
		res.On("Decode", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			taskPtr := args.Get(0).(*domain.Task)
			*taskPtr = *task
		})

		suite.collection.On("FindOne", mock.Anything, mock.Anything).Return(res, nil).Once()

		result, err := suite.repo.GetTaskByID(task.ID)
		suite.NoError(err)
		suite.Equal(task, result)
	})

	// A testcase for the failure of retrieving a task.
	suite.Run("GetTaskByID_Failure", func() {
		id := primitive.NewObjectID()
		res := new(mocks.SingleResult)
		res.On("Decode", mock.Anything).Return(mongo.ErrNoDocuments).Once()
		suite.collection.On("FindOne", mock.Anything, mock.Anything).Return(res, nil).Once()

		result, err := suite.repo.GetTaskByID(id)
		suite.Error(err)
		suite.Equal(domain.Task{}, *result)
	})
}

// A test for the MongoUserRepository.AddTask method.
func (suite *MongoTaskRepositoryTestSuite) TestAddTask() {
	// A testcase for the successful addition of a task.
	suite.Run("AddTask_Success", func() {
		task := mocks.GetNewTask()
		suite.collection.On("InsertOne", mock.Anything, task).Return(&mongo.InsertOneResult{}, nil).Once()

		err := suite.repo.AddTask(task)
		suite.NoError(err)
	})

	// A testcase for the failure of adding a task.
	suite.Run("AddTask_Failure", func() {
		task := mocks.GetNewTask()
		suite.collection.On("InsertOne", mock.Anything, task).Return(nil, mongo.ErrNoDocuments).Once()

		err := suite.repo.AddTask(task)
		suite.Error(err)
	})
}

// A test for the MongoUserRepository.ReplaceTask method.
func (suite *MongoTaskRepositoryTestSuite) TestReplaceTask() {
	// A testcase for the successful replacement of a task.
	suite.Run("ReplaceTask_Success", func() {
		task := mocks.GetNewTask()
		id := task.ID

		res := new(mocks.SingleResult)
		res.On("Err").Return(nil)
		suite.collection.On("FindOneAndReplace", mock.Anything, mock.Anything, task).Return(res, nil).Once()

		err := suite.repo.ReplaceTask(id, task)
		suite.NoError(err)
	})

	// A testcase for the failure of replacing a task.
	suite.Run("ReplaceTask_Failure", func() {
		task := mocks.GetNewTask()
		id := task.ID

		res := new(mocks.SingleResult)
		res.On("Err").Return(mongo.ErrNoDocuments)
		suite.collection.On("FindOneAndReplace", mock.Anything, mock.Anything, task).Return(res, nil).Once()

		err := suite.repo.ReplaceTask(id, task)
		suite.Error(err)
	})
}

// A test for the MongoUserRepository.UpdateTask method.
func (suite *MongoTaskRepositoryTestSuite) TestUpdateTask() {
	// A testcase for the successful update of a task.
	suite.Run("UpdateTask_Success", func() {
		task := mocks.GetNewTask()
		id := task.ID
		taskData := bson.M{"title": "New Title"}

		suite.collection.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).Return(&mongo.UpdateResult{}, nil).Once()

		err := suite.repo.UpdateTask(id, taskData)
		suite.NoError(err)
	})

	// A testcase for the failure of updating a task.
	suite.Run("UpdateTask_Failure", func() {
		task := mocks.GetNewTask()
		id := task.ID
		taskData := bson.M{"title": "New Title"}

		suite.collection.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).Return(&mongo.UpdateResult{}, mongo.ErrClientDisconnected).Once()

		err := suite.repo.UpdateTask(id, taskData)
		suite.Error(err)
	})
}

// A test for the MongoUserRepository.DeleteTask method.
func (suite *MongoTaskRepositoryTestSuite) TestDeleteTask() {
	// A testcase for the successful deletion of a task.
	suite.Run("DeleteTask_Success", func() {
		task := mocks.GetNewTask()
		id := task.ID

		suite.collection.On("DeleteOne", mock.Anything, mock.Anything).Return(&mongo.DeleteResult{}, nil).Once()

		err := suite.repo.DeleteTask(id)
		suite.NoError(err)
	})

	// A testcase for the failure of deleting a task.
	suite.Run("DeleteTask_Failure", func() {
		task := mocks.GetNewTask()
		id := task.ID

		suite.collection.On("DeleteOne", mock.Anything, mock.Anything).Return(&mongo.DeleteResult{}, mongo.ErrClientDisconnected).Once()

		err := suite.repo.DeleteTask(id)
		suite.Error(err)
	})
}

// A function that runs the TestSuite.
func Test_MongoTaskRepository(t *testing.T) {
	suite.Run(t, new(MongoTaskRepositoryTestSuite))
}
