package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"task_manager/delivery/controllers"
	"task_manager/domain"
	"task_manager/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

// A suite to test the UserController.
type UserControllerTestSuite struct {
	suite.Suite
	controller  *controllers.UserController
	mockUsecase *mocks.UserUsecase
}

// A method that initializes the UserControllerTestSuite.
func (suite *UserControllerTestSuite) SetupSuite() {
	suite.mockUsecase = new(mocks.UserUsecase)
	suite.controller = controllers.NewUserController(suite.mockUsecase)
}

// A method that cleans up the UserControllerTestSuite.
func (suite *UserControllerTestSuite) TearDownSuite() {
	suite.mockUsecase.AssertExpectations(suite.T())
}

// A test for the UserController.RegisterUser method.
func (suite *UserControllerTestSuite) TestRegisterUser() {
	// A testcase for a successful user registration.
	suite.Run("RegisterUser_Success", func() {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		userData := mocks.GetAuthUserData()
		user := mocks.GetUser3(userData)
		suite.mockUsecase.On("RegisterUser", userData).Return(user, nil).Once()

		body, err := json.Marshal(userData)
		suite.Nil(err)
		ctx.Request = httptest.NewRequest("POST", "/users", bytes.NewReader(body))

		suite.controller.RegisterUser(ctx)
		expected, err := json.Marshal(user)
		suite.Nil(err)

		suite.Equal(201, w.Code)
		suite.Equal(string(expected), w.Body.String())
	})

	// A testcase for an invalid request.
	suite.Run("RegisterUser_InvalidRequest", func() {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Request = httptest.NewRequest("POST", "/users", nil)

		suite.controller.RegisterUser(ctx)
		expected, err := json.Marshal(gin.H{"error": "Invalid request"})
		suite.Nil(err)

		suite.Equal(400, w.Code)
		suite.Equal(string(expected), w.Body.String())
	})

	// A testcase for an error during user registration.
	suite.Run("RegisterUser_Error", func() {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		userData := mocks.GetAuthUserData()
		suite.mockUsecase.On("RegisterUser", userData).Return(nil, &domain.Error{
			Err:        errors.New("some error"),
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
		}).Once()

		body, err := json.Marshal(userData)
		suite.Nil(err)
		ctx.Request = httptest.NewRequest("POST", "/users", bytes.NewReader(body))

		suite.controller.RegisterUser(ctx)
		expected, err := json.Marshal(gin.H{"error": "Internal Server Error"})
		suite.Nil(err)

		suite.Equal(500, w.Code)
		suite.Equal(string(expected), w.Body.String())
	})
}

// A test for the UserController.Login method.
func (suite *UserControllerTestSuite) TestLogin() {
	// A testcase for a successful user login.
	suite.Run("Login_Success", func() {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		userData := mocks.GetAuthUserData()
		token := "some.random.token.after.login"
		suite.mockUsecase.On("LoginUser", userData).Return(token, nil).Once()

		body, err := json.Marshal(userData)
		suite.Nil(err)
		ctx.Request = httptest.NewRequest("POST", "/users/login", bytes.NewReader(body))

		suite.controller.Login(ctx)
		expected, err := json.Marshal(gin.H{"token": token})
		suite.Nil(err)

		suite.Equal(200, w.Code)
		suite.Equal(string(expected), w.Body.String())
	})

	// A testcase for an invalid request.
	suite.Run("Login_InvalidRequest", func() {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Request = httptest.NewRequest("POST", "/users/login", nil)

		suite.controller.Login(ctx)
		expected, err := json.Marshal(gin.H{"error": "Invalid Request"})
		suite.Nil(err)

		suite.Equal(400, w.Code)
		suite.Equal(string(expected), w.Body.String())
	})

	// A testcase for an error during user login.
	suite.Run("Login_Error", func() {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		userData := mocks.GetAuthUserData()
		suite.mockUsecase.On("LoginUser", userData).Return("", &domain.Error{
			Err:        errors.New("some error"),
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
		}).Once()

		body, err := json.Marshal(userData)
		suite.Nil(err)
		ctx.Request = httptest.NewRequest("POST", "/users/login", bytes.NewReader(body))

		suite.controller.Login(ctx)
		expected, err := json.Marshal(gin.H{"error": "Internal Server Error"})
		suite.Nil(err)

		suite.Equal(500, w.Code)
		suite.Equal(string(expected), w.Body.String())
	})
}

// A test for the UserController.AddUser method.
func (suite *UserControllerTestSuite) TestAddUser() {
	// A testcase for a successful user addition.
	suite.Run("AddUser_Success", func() {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		userData := mocks.GetCreateUserData()
		user := mocks.GetUser(userData)
		claims := mocks.GetClaims()
		suite.mockUsecase.On("AddUser", userData, claims).Return(user, nil).Once()

		body, err := json.Marshal(userData)
		suite.Nil(err)
		ctx.Request = httptest.NewRequest("POST", "/users", bytes.NewReader(body))
		ctx.Set("claims", claims)

		suite.controller.AddUser(ctx)
		expected, err := json.Marshal(user)
		suite.Nil(err)

		suite.Equal(201, w.Code)
		suite.Equal(string(expected), w.Body.String())
	})

	// A testcase for an invalid request.
	suite.Run("AddUser_InvalidRequest", func() {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		claims := mocks.GetClaims()
		ctx.Set("claims", claims)
		ctx.Request = httptest.NewRequest("POST", "/users", nil)

		suite.controller.AddUser(ctx)
		expected, err := json.Marshal(gin.H{"error": "Invalid Request"})
		suite.Nil(err)

		suite.Equal(400, w.Code)
		suite.Equal(string(expected), w.Body.String())
	})

	// A testcase for an error during user addition.
	suite.Run("AddUser_Error", func() {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		userData := mocks.GetCreateUserData()
		claims := mocks.GetClaims()
		suite.mockUsecase.On("AddUser", userData, claims).Return(nil, &domain.Error{
			Err:        errors.New("some error"),
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
		}).Once()

		body, err := json.Marshal(userData)
		suite.Nil(err)
		ctx.Request = httptest.NewRequest("POST", "/users", bytes.NewReader(body))
		ctx.Set("claims", claims)

		suite.controller.AddUser(ctx)
		expected, err := json.Marshal(gin.H{"error": "Internal Server Error"})
		suite.Nil(err)

		suite.Equal(500, w.Code)
		suite.Equal(string(expected), w.Body.String())
	})
}

// A test for the UserController.GetUsers method.
func (suite *UserControllerTestSuite) TestGetUsers() {
	// A testcase for a successful user retrieval.
	suite.Run("GetUsers_Success", func() {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		users := mocks.GetManyUsers()
		suite.mockUsecase.On("GetUsers").Return(users, nil).Once()

		ctx.Request = httptest.NewRequest("GET", "/users", nil)

		suite.controller.GetUsers(ctx)
		expected, err := json.Marshal(gin.H{
			"count": len(users),
			"users": users,
		})
		suite.Nil(err)

		suite.Equal(200, w.Code)
		suite.Equal(string(expected), w.Body.String())
	})

	// A testcase for an error during user retrieval.
	suite.Run("GetUsers_Error", func() {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		suite.mockUsecase.On("GetUsers").Return(nil, &domain.Error{
			Err:        errors.New("some error"),
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
		}).Once()

		ctx.Request = httptest.NewRequest("GET", "/users", nil)

		suite.controller.GetUsers(ctx)
		expected, err := json.Marshal(gin.H{"error": "Internal Server Error"})
		suite.Nil(err)

		suite.Equal(500, w.Code)
		suite.Equal(string(expected), w.Body.String())
	})
}

// A test for the UserController.GetUserByID method.
func (suite *UserControllerTestSuite) TestGetUser() {
	// A testcase for a successful user retrieval.
	suite.Run("GetUser_Success", func() {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		user := mocks.GetNewUser()
		suite.mockUsecase.On("GetUserByID", user.ID).Return(user, nil).Once()

		ctx.Set("user_id", user.ID)
		ctx.Request = httptest.NewRequest("GET", "/users/"+user.ID.Hex(), nil)

		suite.controller.GetUserByID(ctx)
		expected, err := json.Marshal(user)
		suite.Nil(err)

		suite.Equal(200, w.Code)
		suite.Equal(string(expected), w.Body.String())
	})

	// A testcase for an error during user retrieval.
	suite.Run("GetUser_Error", func() {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		user := mocks.GetNewUser()
		suite.mockUsecase.On("GetUserByID", user.ID).Return(nil, &domain.Error{
			Err:        errors.New("some error"),
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
		}).Once()

		ctx.Set("user_id", user.ID)
		ctx.Request = httptest.NewRequest("GET", "/users/"+user.ID.Hex(), nil)

		suite.controller.GetUserByID(ctx)
		expected, err := json.Marshal(gin.H{"error": "Internal Server Error"})
		suite.Nil(err)

		suite.Equal(500, w.Code)
		suite.Equal(string(expected), w.Body.String())
	})
}

// A test for the UserController.UpdateUserPatch method.
func (suite *UserControllerTestSuite) TestUpdateUserPatch() {
	// A testcase for a successful user update.
	suite.Run("UpdateUserPatch_Success", func() {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		user := mocks.GetNewUser()
		userData := mocks.GetUpdateUserData()
		claims := mocks.GetClaims()
		suite.mockUsecase.On("UpdateUser", user.ID, userData, claims).Return(user, nil).Once()

		body, err := json.Marshal(userData)
		suite.Nil(err)
		ctx.Request = httptest.NewRequest("PATCH", "/users/"+user.ID.Hex(), bytes.NewReader(body))
		ctx.Set("user_id", user.ID)
		ctx.Set("claims", claims)

		suite.controller.UpdateUserPatch(ctx)
		expected, err := json.Marshal(user)
		suite.Nil(err)

		suite.Equal(200, w.Code)
		suite.Equal(string(expected), w.Body.String())
	})

	// A testcase for an invalid request.
	suite.Run("UpdateUserPatch_InvalidRequest", func() {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		user := mocks.GetNewUser()
		claims := mocks.GetClaims()
		ctx.Set("user_id", user.ID)
		ctx.Set("claims", claims)
		ctx.Request = httptest.NewRequest("PATCH", "/users/"+user.ID.Hex(), nil)

		suite.controller.UpdateUserPatch(ctx)
		expected, err := json.Marshal(gin.H{"error": "Invalid Request"})
		suite.Nil(err)

		suite.Equal(400, w.Code)
		suite.Equal(string(expected), w.Body.String())
	})

	// A testcase for an error during user update.
	suite.Run("UpdateUserPatch_Error", func() {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		user := mocks.GetNewUser()
		userData := mocks.GetUpdateUserData()
		claims := mocks.GetClaims()
		suite.mockUsecase.On("UpdateUser", user.ID, userData, claims).Return(nil, &domain.Error{
			Err:        errors.New("some error"),
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
		}).Once()

		body, err := json.Marshal(userData)
		suite.Nil(err)
		ctx.Request = httptest.NewRequest("PATCH", "/users/"+user.ID.Hex(), bytes.NewReader(body))
		ctx.Set("user_id", user.ID)
		ctx.Set("claims", claims)

		suite.controller.UpdateUserPatch(ctx)
		expected, err := json.Marshal(gin.H{"error": "Internal Server Error"})
		suite.Nil(err)

		suite.Equal(500, w.Code)
		suite.Equal(string(expected), w.Body.String())
	})
}

// A test for the UserController.DeleteUser method.
func (suite *UserControllerTestSuite) TestDeleteUser() {
	// A testcase for a successful user deletion.
	suite.Run("DeleteUser_Success", func() {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		user := mocks.GetNewUser()
		claims := mocks.GetClaims()
		suite.mockUsecase.On("DeleteUser", user.ID, claims).Return(nil).Once()

		ctx.Set("user_id", user.ID)
		ctx.Set("claims", claims)
		ctx.Request = httptest.NewRequest("DELETE", "/users/"+user.ID.Hex(), nil)

		suite.controller.DeleteUser(ctx)

		suite.Equal(204, w.Code)
		suite.Empty(w.Body.String())
	})

	// A testcase for an error during user deletion.
	suite.Run("DeleteUser_Error", func() {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		user := mocks.GetNewUser()
		claims := mocks.GetClaims()
		suite.mockUsecase.On("DeleteUser", user.ID, claims).Return(&domain.Error{
			Err:        errors.New("some error"),
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
		}).Once()

		ctx.Set("user_id", user.ID)
		ctx.Set("claims", claims)
		ctx.Request = httptest.NewRequest("DELETE", "/users/"+user.ID.Hex(), nil)

		suite.controller.DeleteUser(ctx)
		expected, err := json.Marshal(gin.H{"error": "Internal Server Error"})
		suite.Nil(err)

		suite.Equal(500, w.Code)
		suite.Equal(string(expected), w.Body.String())
	})
}

// A function that runs the UserControllerTestSuite.
func Test_UserControllerTestSuite(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
}
