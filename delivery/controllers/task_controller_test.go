package controllers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"task_manager/delivery/controllers"
	"task_manager/domain"
	"task_manager/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

// A suite to test the TaskController.
type TaskControllerTestSuite struct {
	suite.Suite
	controller *controllers.TaskController
	usecase    *mocks.TaskUsecase
}

// A method that initializes the TaskControllerTestSuite.
func (suite *TaskControllerTestSuite) SetupSuite() {
	suite.usecase = new(mocks.TaskUsecase)
	suite.controller = controllers.NewTaskController(suite.usecase)
}

// A method that closes the suite.
func (suite *TaskControllerTestSuite) TearDownSuite() {
	suite.usecase.AssertExpectations(suite.T())
}

// A test for the TaskController.GetTasks method.
func (suite *TaskControllerTestSuite) TestGetTasks() {
	// A testcase when the usecase returns a slice of tasks.
	suite.Run("Tasks", func() {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		tasks := mocks.GetManyTasks()
		suite.usecase.On("GetTasks").Return(tasks, nil).Once()

		ctx.Request = httptest.NewRequest("GET", "/tasks", nil)

		suite.controller.GetTasks(ctx)

		expected, err := json.Marshal(gin.H{
			"count": len(tasks),
			"tasks": tasks,
		})
		suite.Nil(err)

		suite.Equal(200, w.Code)
		suite.Equal(string(expected), w.Body.String())
	})

	// A testcase when the usecase returns an error.
	suite.Run("Error", func() {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		suite.usecase.On("GetTasks").Return(nil, &domain.Error{
			Err:        errors.New("some error"),
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
		}).Once()

		ctx.Request = httptest.NewRequest("GET", "/tasks", nil)

		suite.controller.GetTasks(ctx)

		expected, err := json.Marshal(gin.H{"error": "Internal Server Error"})
		suite.Nil(err)

		suite.Equal(500, w.Code)
		suite.Equal(string(expected), w.Body.String())
	})
}

// A test for the TaskController.GetTaskByID method.
func (suite *TaskControllerTestSuite) TestGetTaskByID() {
	// A testcase when the task is found.
	suite.Run("TaskFound", func() {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		taskID := mocks.GetPrimitiveID1()
		task := mocks.GetNewTask()
		suite.usecase.On("GetTaskByID", taskID).Return(task, nil).Once()
		ctx.Set("task_id", taskID)

		suite.controller.GetTaskByID(ctx)

		expected, err := json.Marshal(task)
		suite.Nil(err)

		suite.Equal(200, w.Code)
		suite.Equal(string(expected), w.Body.String())
	})

	// A testcase when the task is not found.
	suite.Run("TaskNotFound", func() {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		taskID := mocks.GetPrimitiveID1()
		suite.usecase.On("GetTaskByID", taskID).Return(nil, &domain.Error{
			Err:        errors.New("task not found"),
			StatusCode: http.StatusNotFound,
			Message:    "Task Not Found",
		}).Once()
		ctx.Set("task_id", taskID)

		suite.controller.GetTaskByID(ctx)
		expected, err := json.Marshal(gin.H{"error": "Task Not Found"})
		suite.Nil(err)

		suite.Equal(404, w.Code)
		suite.Equal(string(expected), w.Body.String())
	})
}

// A test for the TaskController.CreateTask method.
func (suite *TaskControllerTestSuite) TestCreateTask() {
	// A testcase when the task is created successfully.
	suite.Run("TaskCreated", func() {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		taskData := mocks.GetCreateTaskData()
		claims := mocks.GetClaims()
		ctx.Set("claims", claims)
		taskView := mocks.GetView(taskData, claims)

		body, err := json.Marshal(taskData)
		suite.Nil(err)

		ctx.Request = httptest.NewRequest("POST", "/tasks", strings.NewReader(string(body)))
		suite.usecase.On("CreateTask", taskData, claims).Return(taskView, nil).Once()

		suite.controller.CreateTask(ctx)
		expected, err := json.Marshal(taskView)
		suite.Nil(err)

		suite.Equal(201, w.Code)
		suite.Equal(string(expected), w.Body.String())
	})

	// A testcase when the request body is invalid.
	suite.Run("InvalidRequestBody", func() {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		claims := mocks.GetClaims()
		ctx.Set("claims", claims)
		ctx.Request = httptest.NewRequest("POST", "/tasks", nil)

		suite.controller.CreateTask(ctx)

		expected, err := json.Marshal(gin.H{"error": "Invalid request"})
		suite.Nil(err)

		suite.Equal(400, w.Code)
		suite.Equal(string(expected), w.Body.String())
	})

	// A testcase when the status field is invalid.
	suite.Run("InvalidStatus", func() {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		taskData := mocks.GetCreateTaskData()
		taskData.Status = "Invalid"
		claims := mocks.GetClaims()
		ctx.Set("claims", claims)

		body, err := json.Marshal(taskData)
		suite.Nil(err)

		ctx.Request = httptest.NewRequest("POST", "/tasks", strings.NewReader(string(body)))

		suite.controller.CreateTask(ctx)

		expected, err := json.Marshal(gin.H{"error": "status field must be one of: Pending, Completed, In Progress"})
		suite.Nil(err)

		suite.Equal(400, w.Code)
		suite.Equal(string(expected), w.Body.String())
	})

	// A testcase when the usecase returns an error.
	suite.Run("Error", func() {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		taskData := mocks.GetCreateTaskData()
		claims := mocks.GetClaims()
		ctx.Set("claims", claims)
		body, err := json.Marshal(taskData)
		suite.Nil(err)
		ctx.Request = httptest.NewRequest("POST", "/tasks", strings.NewReader(string(body)))

		suite.usecase.On("CreateTask", taskData, claims).Return(nil, &domain.Error{
			Err:        errors.New("some error"),
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
		}).Once()

		suite.controller.CreateTask(ctx)

		expected, err := json.Marshal(gin.H{"error": "Internal Server Error"})
		suite.Nil(err)

		suite.Equal(500, w.Code)
		suite.Equal(string(expected), w.Body.String())
	})

	// A testcase when the status field is missing.
	suite.Run("StatusMissing", func() {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		taskData := mocks.GetCreateTaskData()
		claims := mocks.GetClaims()
		ctx.Set("claims", claims)

		taskData.Status = ""
		body, err := json.Marshal(taskData)
		suite.Nil(err)
		ctx.Request = httptest.NewRequest("POST", "/tasks", strings.NewReader(string(body)))

		taskData.Status = "Pending"
		taskView := mocks.GetView(taskData, claims)
		suite.usecase.On("CreateTask", taskData, claims).Return(taskView, nil).Once()

		suite.controller.CreateTask(ctx)

		expected, err := json.Marshal(taskView)
		suite.Nil(err)

		suite.Equal(201, w.Code)
		suite.Equal(string(expected), w.Body.String())
	})
}

// A test for the TaskController.UpdateTaskPut method.
func (suite *TaskControllerTestSuite) TestUpdateTaskPut() {
	// A testcase when the task is updated successfully.
	suite.Run("TaskUpdated", func() {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		taskID := mocks.GetPrimitiveID1()
		claims := mocks.GetClaims()
		ctx.Set("claims", claims)
		taskData := mocks.GetReplaceTaskData()
		ctx.Set("task_id", taskID)
		taskView := mocks.GetView2(taskData, claims)

		body, err := json.Marshal(taskData)
		suite.Nil(err)
		ctx.Request = httptest.NewRequest("PUT", "/tasks/"+taskID.Hex(), strings.NewReader(string(body)))
		suite.usecase.On("ReplaceTask", taskID, taskData, claims).Return(taskView, nil).Once()

		suite.controller.UpdateTaskPut(ctx)
		expected, err := json.Marshal(taskView)
		suite.Nil(err)

		suite.Equal(200, w.Code)
		suite.Equal(string(expected), w.Body.String())
	})

	// A testcase when the request body is invalid.
	suite.Run("InvalidRequestBody", func() {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		taskID := mocks.GetPrimitiveID1()
		claims := mocks.GetClaims()
		ctx.Set("task_id", taskID)
		ctx.Set("claims", claims)

		ctx.Request = httptest.NewRequest("PUT", "/tasks/"+taskID.Hex(), nil)

		suite.controller.UpdateTaskPut(ctx)

		expected, err := json.Marshal(gin.H{"error": "Invalid request"})
		suite.Nil(err)

		suite.Equal(400, w.Code)
		suite.Equal(string(expected), w.Body.String())
	})

	// A testcase when the status field is invalid.
	suite.Run("InvalidStatus", func() {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		taskID := mocks.GetPrimitiveID1()
		claims := mocks.GetClaims()
		ctx.Set("task_id", taskID)
		ctx.Set("claims", claims)
		taskData := mocks.GetReplaceTaskData()
		taskData.Status = "Invalid"

		body, err := json.Marshal(taskData)
		suite.Nil(err)
		ctx.Request = httptest.NewRequest("PUT", "/tasks/"+taskID.Hex(), strings.NewReader(string(body)))

		suite.controller.UpdateTaskPut(ctx)

		expected, err := json.Marshal(gin.H{"error": "status field must be one of: Pending, Completed, In Progress"})
		suite.Nil(err)

		suite.Equal(400, w.Code)
		suite.Equal(string(expected), w.Body.String())
	})

	// A testcase when the usecase returns an error.
	suite.Run("Error", func() {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		taskID := mocks.GetPrimitiveID1()
		claims := mocks.GetClaims()
		ctx.Set("task_id", taskID)
		ctx.Set("claims", claims)
		taskData := mocks.GetReplaceTaskData()

		body, err := json.Marshal(taskData)
		suite.Nil(err)
		ctx.Request = httptest.NewRequest("PUT", "/tasks/"+taskID.Hex(), strings.NewReader(string(body)))

		suite.usecase.On("ReplaceTask", taskID, taskData, claims).Return(nil, &domain.Error{
			Err:        errors.New("some error"),
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
		}).Once()

		suite.controller.UpdateTaskPut(ctx)

		expected, err := json.Marshal(gin.H{"error": "Internal Server Error"})
		suite.Nil(err)

		suite.Equal(500, w.Code)
		suite.Equal(string(expected), w.Body.String())
	})
}

// A test for the TaskController.UpdateTaskPatch method.
func (suite *TaskControllerTestSuite) TestUpdateTaskPatch() {
	// A testcase when the task is updated successfully.
	suite.Run("TaskUpdated", func() {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		taskID := mocks.GetPrimitiveID1()
		claims := mocks.GetClaims()
		ctx.Set("claims", claims)
		taskData := mocks.GetUpdateTaskData()
		ctx.Set("task_id", taskID)
		taskView := mocks.GetView3(taskData, claims)

		body, err := json.Marshal(taskData)
		suite.Nil(err)
		ctx.Request = httptest.NewRequest("PATCH", "/tasks/"+taskID.Hex(), strings.NewReader(string(body)))
		suite.usecase.On("UpdateTask", taskID, taskData, claims).Return(taskView, nil).Once()

		suite.controller.UpdateTaskPatch(ctx)
		expected, err := json.Marshal(taskView)
		suite.Nil(err)

		suite.Equal(200, w.Code)
		suite.Equal(string(expected), w.Body.String())
	})

	// A testcase when the request body is invalid.
	suite.Run("InvalidRequestBody", func() {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		taskID := mocks.GetPrimitiveID1()
		claims := mocks.GetClaims()
		ctx.Set("task_id", taskID)
		ctx.Set("claims", claims)

		ctx.Request = httptest.NewRequest("PATCH", "/tasks/"+taskID.Hex(), nil)

		suite.controller.UpdateTaskPatch(ctx)

		expected, err := json.Marshal(gin.H{"error": "Invalid request"})
		suite.Nil(err)

		suite.Equal(400, w.Code)
		suite.Equal(string(expected), w.Body.String())
	})

	// A testcase when the status field is invalid.
	suite.Run("InvalidStatus", func() {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		taskID := mocks.GetPrimitiveID1()
		claims := mocks.GetClaims()
		ctx.Set("task_id", taskID)
		ctx.Set("claims", claims)
		taskData := mocks.GetUpdateTaskData()
		taskData.Status = "Invalid"

		body, err := json.Marshal(taskData)
		suite.Nil(err)
		ctx.Request = httptest.NewRequest("PATCH", "/tasks/"+taskID.Hex(), strings.NewReader(string(body)))

		suite.controller.UpdateTaskPatch(ctx)

		expected, err := json.Marshal(gin.H{"error": "status field must be one of: Pending, Completed, In Progress"})
		suite.Nil(err)

		suite.Equal(400, w.Code)
		suite.Equal(string(expected), w.Body.String())
	})

	// A testcase when the usecase returns an error.
	suite.Run("Error", func() {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		taskID := mocks.GetPrimitiveID1()
		claims := mocks.GetClaims()
		ctx.Set("task_id", taskID)
		ctx.Set("claims", claims)
		taskData := mocks.GetUpdateTaskData()

		body, err := json.Marshal(taskData)
		suite.Nil(err)
		ctx.Request = httptest.NewRequest("PATCH", "/tasks/"+taskID.Hex(), strings.NewReader(string(body)))

		suite.usecase.On("UpdateTask", taskID, taskData, claims).Return(nil, &domain.Error{
			Err:        errors.New("some error"),
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
		}).Once()

		suite.controller.UpdateTaskPatch(ctx)

		expected, err := json.Marshal(gin.H{"error": "Internal Server Error"})
		suite.Nil(err)

		suite.Equal(500, w.Code)
		suite.Equal(string(expected), w.Body.String())
	})
}

// A test for the TaskController.DeleteTask method.
func (suite *TaskControllerTestSuite) TestDeleteTask() {
	// A testcase when the task is deleted successfully.
	suite.Run("TaskDeleted", func() {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		taskID := mocks.GetPrimitiveID1()
		claims := mocks.GetClaims()
		ctx.Set("claims", claims)
		ctx.Set("task_id", taskID)

		suite.usecase.On("DeleteTask", taskID, claims).Return(nil).Once()

		suite.controller.DeleteTask(ctx)

		suite.Equal(204, w.Code)
		suite.Empty(w.Body.String())
	})

	// A testcase when the usecase returns an error.
	suite.Run("Error", func() {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		taskID := mocks.GetPrimitiveID1()
		claims := mocks.GetClaims()
		ctx.Set("claims", claims)
		ctx.Set("task_id", taskID)

		suite.usecase.On("DeleteTask", taskID, claims).Return(&domain.Error{
			Err:        errors.New("some error"),
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
		}).Once()

		suite.controller.DeleteTask(ctx)

		expected, err := json.Marshal(gin.H{"error": "Internal Server Error"})
		suite.Nil(err)

		suite.Equal(500, w.Code)
		suite.Equal(string(expected), w.Body.String())
	})
}

// A function that runs the TaskControllerTestSuite.
func Test_TaskControllerTest(t *testing.T) {
	suite.Run(t, new(TaskControllerTestSuite))
}
