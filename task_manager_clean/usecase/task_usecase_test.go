package usecase_test

import (
	"errors"
	"net/http"
	"task_manager/domain"
	"task_manager/mocks"
	"task_manager/usecase"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	mockTask     = mock.AnythingOfType("*domain.Task")
	mockObjectID = mock.AnythingOfType("primitive.ObjectID")
	mockBSON     = mock.AnythingOfType("primitive.M")
)

// A suite for the TaskUsecase.
type TaskUsecaseSuite struct {
	suite.Suite
	taskRepo *mocks.TaskRepository
	usecase  *usecase.TaskUsecase
}

// A method that sets up the TestSuite.
func (suite *TaskUsecaseSuite) SetupSuite() {
	suite.taskRepo = new(mocks.TaskRepository)
	suite.usecase = usecase.NewTaskUsecase(suite.taskRepo)
}

// A method that tears down the TestSuite.
func (suite *TaskUsecaseSuite) TearDownSuite() {
	suite.taskRepo.AssertExpectations(suite.T())
}

// A test for the TaskUsecase.GetTasks method.
func (suite *TaskUsecaseSuite) Test_GetTasks() {
	// A testcase where the task repository returns an empty list of tasks.
	suite.Run("GetTasks_Empty", func() {
		suite.taskRepo.On("GetAllTasks").Return([]domain.Task{}, nil).Once()
		tasks, err := suite.usecase.GetTasks()

		suite.Equal(0, len(tasks))
		suite.Nil(err)
	})

	// A testcase where the task repository returns a non-empty list of tasks.
	suite.Run("GetTasks_Success", func() {
		tasks := mocks.GetManyTasks()
		suite.taskRepo.On("GetAllTasks").Return(tasks, nil).Once()

		result, err := suite.usecase.GetTasks()
		suite.Nil(err)
		suite.Equal(tasks, result)
	})

	// A testcase where the task repository returns an error.
	suite.Run("GetTasks_Error", func() {
		suite.taskRepo.On("GetAllTasks").Return(nil, errors.New("some error")).Once()

		result, err := suite.usecase.GetTasks()
		suite.Nil(result)

		expectedErr := &domain.Error{
			Err:        errors.New("some error"),
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}

		suite.Equal(expectedErr, err)
	})
}

// A test for the TaskUsecase.GetTaskByID method.
func (suite *TaskUsecaseSuite) Test_GetTaskByID() {
	// A testcase where the task repository returns a task.
	suite.Run("GetTaskByID_Success", func() {
		task := mocks.GetNewTask()
		suite.taskRepo.On("GetTaskByID", mockObjectID).Return(task, nil).Once()

		result, err := suite.usecase.GetTaskByID(task.ID)
		suite.Equal(task, result)
		suite.Nil(err)
	})

	// A testcase where the task repository fails to find a task.
	suite.Run("GetTaskByID_NotFound", func() {
		id := primitive.NewObjectID()
		suite.taskRepo.On("GetTaskByID", mockObjectID).Return(nil, mongo.ErrNoDocuments).Once()

		result, err := suite.usecase.GetTaskByID(id)
		suite.Nil(result)

		expectedErr := &domain.Error{
			Err:        mongo.ErrNoDocuments,
			StatusCode: http.StatusNotFound,
			Message:    "Task not found",
		}

		suite.Equal(expectedErr, err)
	})

	// A testcase where the task repository returns an error.
	suite.Run("GetTaskByID_Error", func() {
		id := primitive.NewObjectID()
		suite.taskRepo.On("GetTaskByID", mockObjectID).Return(nil, errors.New("some error")).Once()

		result, err := suite.usecase.GetTaskByID(id)
		suite.Nil(result)

		expectedErr := &domain.Error{
			Err:        errors.New("some error"),
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}

		suite.Equal(expectedErr, err)
	})
}

// A test for the TaskUsecase.CreateTask method.
func (suite *TaskUsecaseSuite) Test_CreateTask() {
	// A testcase where the task repository successfully creates a task.
	suite.Run("CreateTask_Success", func() {
		taskData := mocks.GetCreateTaskData()
		claims := mocks.GetClaims()
		task := mocks.GetTask(taskData, claims)
		task.ID = mocks.GetNextID(primitive.NewObjectID())
		taskView := mocks.GetTaskView(task)

		suite.taskRepo.On("AddTask", mockTask).Return(nil).Once()

		result, err := suite.usecase.CreateTask(taskData, claims)
		suite.Equal(taskView, result)
		suite.Nil(err)
	})

	// A testcase where the task repository returns an error.
	suite.Run("CreateTask_Error", func() {
		taskData := mocks.GetCreateTaskData()
		claims := mocks.GetClaims()
		task := mocks.GetTask(taskData, claims)
		task.ID = mocks.GetNextID(primitive.NewObjectID())

		suite.taskRepo.On("AddTask", mockTask).Return(errors.New("some error")).Once()

		result, err := suite.usecase.CreateTask(taskData, claims)
		suite.Nil(result)

		expectedErr := &domain.Error{
			Err:        errors.New("some error"),
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}

		suite.Equal(expectedErr, err)
	})
}

// A test for the TaskUsecase.ReplaceTask method.
func (suite *TaskUsecaseSuite) Test_ReplaceTask() {
	// A testcase where the task repository successfully replaces a task.
	suite.Run("ReplaceTask_Success", func() {
		taskData := mocks.GetReplaceTaskData()
		claims := mocks.GetClaims()
		task := mocks.GetTask2(taskData, claims)
		task.ID = mocks.GetNextID(primitive.NewObjectID())
		objectID := task.ID
		taskView := mocks.GetTaskView(task)

		suite.taskRepo.On("GetTaskByID", mockObjectID).Return(task, nil).Once()
		suite.taskRepo.On("ReplaceTask", mockObjectID, mockTask).Return(nil).Once()

		result, err := suite.usecase.ReplaceTask(objectID, taskData, claims)
		suite.Equal(taskView, result)
		suite.Nil(err)
	})

	// A tescase where the get task function returns an error
	suite.Run("ReplaceTask_Error", func() {
		taskData := mocks.GetReplaceTaskData()
		claims := mocks.GetClaims()
		task := mocks.GetTask2(taskData, claims)
		task.ID = mocks.GetNextID(primitive.NewObjectID())
		objectID := task.ID

		suite.taskRepo.On("GetTaskByID", mockObjectID).Return(nil, errors.New("some error")).Once()

		result, err := suite.usecase.ReplaceTask(objectID, taskData, claims)
		suite.Nil(result)

		expectedErr := &domain.Error{
			Err:        errors.New("some error"),
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}

		suite.Equal(expectedErr, err)
	})

	// A testcase for replacing other user's task.
	suite.Run("ReplaceTask_OtherUser", func() {
		taskData := mocks.GetReplaceTaskData()
		claims := mocks.GetClaims()
		task := mocks.GetTask2(taskData, claims)
		task.ID = mocks.GetNextID(primitive.NewObjectID())
		objectID := task.ID
		task.UserID = mocks.GetNextID(primitive.NewObjectID())

		suite.taskRepo.On("GetTaskByID", mockObjectID).Return(task, nil).Once()

		result, err := suite.usecase.ReplaceTask(objectID, taskData, claims)
		suite.Nil(result)

		expectedErr := &domain.Error{
			Err:        errors.New("trying to replace another user's task"),
			StatusCode: http.StatusForbidden,
			Message:    "A User can only update their own task",
		}

		suite.Equal(expectedErr, err)
	})

	// A testcase where the replace task function returns an error.
	suite.Run("ReplaceTask_Error2", func() {
		taskData := mocks.GetReplaceTaskData()
		claims := mocks.GetClaims()
		task := mocks.GetTask2(taskData, claims)
		task.ID = mocks.GetNextID(primitive.NewObjectID())
		objectID := task.ID

		suite.taskRepo.On("GetTaskByID", mockObjectID).Return(task, nil).Once()
		suite.taskRepo.On("ReplaceTask", mockObjectID, mockTask).Return(errors.New("some error")).Once()

		result, err := suite.usecase.ReplaceTask(objectID, taskData, claims)
		suite.Nil(result)

		expectedErr := &domain.Error{
			Err:        errors.New("some error"),
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}

		suite.Equal(expectedErr, err)
	})
}

// A test for the TaskUsecase.UpdateTask method.
func (suite *TaskUsecaseSuite) Test_UpdateTask() {
	// A testcase where the task does not exist.
	suite.Run("UpdateTask_NotFound", func() {
		taskData := mocks.GetUpdateTaskData()
		claims := mocks.GetClaims()
		id := primitive.NewObjectID()

		suite.taskRepo.On("GetTaskByID", mockObjectID).Return(nil, mongo.ErrNoDocuments).Once()

		expectedErr := &domain.Error{
			Err:        mongo.ErrNoDocuments,
			StatusCode: http.StatusNotFound,
			Message:    "Task not found",
		}

		result, err := suite.usecase.UpdateTask(id, taskData, claims)
		suite.Nil(result)
		suite.Equal(expectedErr, err)
	})

	// A testcase where the task repository successfully updates a task.
	suite.Run("UpdateTask_Success", func() {
		taskData := mocks.GetUpdateTaskData()
		claims := mocks.GetClaims()
		task := mocks.GetTask3(taskData, claims)
		task.UserID = claims.ID
		objectID := task.ID
		taskView := mocks.GetTaskView(task)

		suite.taskRepo.On("GetTaskByID", mockObjectID).Return(task, nil).Once()
		suite.taskRepo.On("UpdateTask", mockObjectID, mockBSON).Return(nil).Once()
		taskData.DueDate = time.Time{}
		taskData.Status = ""

		result, err := suite.usecase.UpdateTask(objectID, taskData, claims)
		suite.Equal(taskView, result)
		suite.Nil(err)
	})

	// A second testcase where the task repository successfully updates a task.
	suite.Run("UpdateTask_Success2", func() {
		taskData := mocks.GetUpdateTaskData()
		claims := mocks.GetClaims()
		task := mocks.GetTask3(taskData, claims)
		task.UserID = claims.ID
		objectID := task.ID
		taskView := mocks.GetTaskView(task)

		suite.taskRepo.On("GetTaskByID", mockObjectID).Return(task, nil).Once()
		suite.taskRepo.On("UpdateTask", mockObjectID, mockBSON).Return(nil).Once()
		taskData.Title = ""
		taskData.Description = ""

		result, err := suite.usecase.UpdateTask(objectID, taskData, claims)
		suite.Equal(taskView, result)
		suite.Nil(err)
	})

	// A testcase where the task repository returns an error.
	suite.Run("UpdateTask_Error", func() {
		taskData := mocks.GetUpdateTaskData()
		claims := mocks.GetClaims()
		task := mocks.GetTask3(taskData, claims)
		task.UserID = claims.ID
		objectID := task.ID

		suite.taskRepo.On("GetTaskByID", mockObjectID).Return(task, nil).Once()
		suite.taskRepo.On("UpdateTask", mockObjectID, mockBSON).Return(errors.New("some error")).Once()

		result, err := suite.usecase.UpdateTask(objectID, taskData, claims)
		suite.Nil(result)

		expectedErr := &domain.Error{
			Err:        errors.New("some error"),
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}

		suite.Equal(expectedErr, err)
	})

	// A testcase for updating another user's task.
	suite.Run("UpdateTask_OtherUser", func() {
		taskData := mocks.GetUpdateTaskData()
		claims := mocks.GetClaims()
		task := mocks.GetTask3(taskData, claims)
		task.UserID = mocks.GetNextID(primitive.NewObjectID())
		objectID := task.ID

		suite.taskRepo.On("GetTaskByID", mockObjectID).Return(task, nil).Once()

		result, err := suite.usecase.UpdateTask(objectID, taskData, claims)
		suite.Nil(result)

		expectedErr := &domain.Error{
			Err:        errors.New("trying to update another user's task"),
			StatusCode: http.StatusForbidden,
			Message:    "A User can only update their own task",
		}

		suite.Equal(expectedErr, err)
	})
}

// A test for the TaskUsecase.DeleteTask method.
func (suite *TaskUsecaseSuite) Test_DeleteTask() {
	// A testcase where the task repository successfully deletes a task.
	suite.Run("DeleteTask_Success", func() {
		claims := mocks.GetClaims()
		task := mocks.GetNewTask()
		task.UserID = claims.ID
		objectID := task.ID

		suite.taskRepo.On("GetTaskByID", mockObjectID).Return(task, nil).Once()
		suite.taskRepo.On("DeleteTask", mockObjectID).Return(nil).Once()

		err := suite.usecase.DeleteTask(objectID, claims)
		suite.Nil(err)
	})

	// A testcase where the task repository returns an error.
	suite.Run("DeleteTask_Error", func() {
		claims := mocks.GetClaims()
		task := mocks.GetNewTask()
		task.UserID = claims.ID
		objectID := task.ID

		suite.taskRepo.On("GetTaskByID", mockObjectID).Return(task, nil).Once()
		suite.taskRepo.On("DeleteTask", mockObjectID).Return(errors.New("some error")).Once()

		err := suite.usecase.DeleteTask(objectID, claims)

		expectedErr := &domain.Error{
			Err:        errors.New("some error"),
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}

		suite.Equal(expectedErr, err)
	})

	// A testcase for deleting another user's task.
	suite.Run("DeleteTask_OtherUser", func() {
		claims := mocks.GetClaims()
		task := mocks.GetNewTask()
		objectID := task.ID
		task.UserID = mocks.GetNextID(primitive.NewObjectID())

		suite.taskRepo.On("GetTaskByID", mockObjectID).Return(task, nil).Once()

		err := suite.usecase.DeleteTask(objectID, claims)

		expectedErr := &domain.Error{
			Err:        errors.New("trying to delete another user's task"),
			StatusCode: http.StatusForbidden,
			Message:    "A User can only delete their own task",
		}

		suite.Equal(expectedErr, err)
	})

	// A testcase where the task does not exist.
	suite.Run("DeleteTask_TaskDoesNotExist", func() {
		claims := mocks.GetClaims()
		id := primitive.NewObjectID()

		suite.taskRepo.On("GetTaskByID", mockObjectID).Return(nil, mongo.ErrNoDocuments).Once()

		err := suite.usecase.DeleteTask(id, claims)

		expectedErr := &domain.Error{
			Err:        mongo.ErrNoDocuments,
			StatusCode: http.StatusNotFound,
			Message:    "Task not found",
		}

		suite.Equal(expectedErr, err)
	})
}

// A method that runs the TestSuite.
func TestTaskUsecaseSuite(t *testing.T) {
	suite.Run(t, new(TaskUsecaseSuite))
}
