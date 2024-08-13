package usecase_test

import (
	"errors"
	"net/http"
	"os"
	"task_manager/domain"
	"task_manager/infrastructure"
	"task_manager/mocks"
	"task_manager/usecase"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	mockUser   = mock.AnythingOfType("*domain.User")
	mockString = mock.AnythingOfType("string")
)

// A suite that tests the user usecase.
type UserUsecaseSuite struct {
	suite.Suite
	userRepo    *mocks.UserRepository
	userUsecase *usecase.UserUsecase
}

// A method that sets up the test suite.
func (suite *UserUsecaseSuite) SetupTest() {
	suite.userRepo = new(mocks.UserRepository)
	suite.userUsecase = usecase.NewUserUsecase(suite.userRepo)
	os.Setenv("JWT_KEY", "test_key")
}

// A method that tears down the test suite.
func (suite *UserUsecaseSuite) TearDownTest() {
	suite.userRepo.AssertExpectations(suite.T())
	os.Unsetenv("JWT_KEY")
}

// A test for the UserUsecase.AddUser method.
func (suite *UserUsecaseSuite) Test_AddUser() {
	// A testcase that tests the successful addition of a user.
	suite.Run("AddUser_Success", func() {
		userData := mocks.GetCreateUserData()
		claims := mocks.GetClaims2() // An admin user.
		user := mocks.GetUser(userData)
		passWord := user.Password

		suite.userRepo.On("AddUser", mockUser).Return(nil).Once()
		suite.userRepo.On("GetUserByUsername", mockString).Return(nil, errors.New("some error")).Once()

		foundUser, err := suite.userUsecase.AddUser(userData, claims)
		suite.Nil(err)
		suite.Nil(infrastructure.ComparePasswords(foundUser.Password, passWord))
		user.Password = foundUser.Password
		suite.Equal(user, foundUser)
	})

	// A testcase that tests the failure of adding a user.
	suite.Run("AddUser_Failure", func() {
		userData := mocks.GetCreateUserData()
		claims := mocks.GetClaims2() // An admin user.

		suite.userRepo.On("AddUser", mockUser).Return(errors.New("some error")).Once()
		suite.userRepo.On("GetUserByUsername", mockString).Return(nil, errors.New("some error")).Once()

		expectedError := &domain.Error{
			Err:        errors.New("some error"),
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}

		foundUser, err := suite.userUsecase.AddUser(userData, claims)
		suite.Nil(foundUser)
		suite.Equal(expectedError, err)
	})

	// A testcase where the user already exists.
	suite.Run("AddUser_UserExists", func() {
		userData := mocks.GetCreateUserData()
		claims := mocks.GetClaims2() // An admin user.

		suite.userRepo.On("GetUserByUsername", mockString).Return(mocks.GetUser(userData), nil).Once()

		expectedError := &domain.Error{
			Err:        errors.New("conflict"),
			StatusCode: http.StatusConflict,
			Message:    "Username already exists",
		}

		foundUser, err := suite.userUsecase.AddUser(userData, claims)
		suite.Nil(foundUser)
		suite.Equal(expectedError, err)
	})

	// A testcase where the user is not authorized.
	suite.Run("AddUser_Unauthorized", func() {
		userData := mocks.GetCreateUserData()
		claims := mocks.GetClaims() // A regular user.

		expectedError := &domain.Error{
			Err:        errors.New("unauthorized"),
			StatusCode: http.StatusForbidden,
			Message:    "A User cannot add a new user",
		}

		foundUser, err := suite.userUsecase.AddUser(userData, claims)
		suite.Nil(foundUser)
		suite.Equal(expectedError, err)
	})
}

// A test for the UserUsecase.RegisterUser method.
func (suite *UserUsecaseSuite) Test_RegisterUser() {
	// A testcase that tests the successful registration of a user.
	suite.Run("RegisterUser_Success", func() {
		userData := mocks.GetAuthUserData()
		user := mocks.GetUser3(userData)
		passWord := user.Password

		suite.userRepo.On("AddUser", mockUser).Return(nil).Once()
		suite.userRepo.On("GetUserByUsername", mockString).Return(nil, errors.New("some error")).Once()

		foundUser, err := suite.userUsecase.RegisterUser(userData)
		suite.Nil(err)
		suite.Nil(infrastructure.ComparePasswords(foundUser.Password, passWord))
		user.Password = foundUser.Password
		suite.Equal(user, foundUser)
	})

	// A testcase that tests the failure of registering a user.
	suite.Run("RegisterUser_Failure", func() {
		userData := mocks.GetAuthUserData()

		suite.userRepo.On("AddUser", mockUser).Return(errors.New("some error")).Once()
		suite.userRepo.On("GetUserByUsername", mockString).Return(nil, errors.New("some error")).Once()

		expectedError := &domain.Error{
			Err:        errors.New("some error"),
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}

		foundUser, err := suite.userUsecase.RegisterUser(userData)
		suite.Nil(foundUser)
		suite.Equal(expectedError, err)
	})

	// A testcase where the user already exists.
	suite.Run("RegisterUser_UserExists", func() {
		userData := mocks.GetAuthUserData()

		suite.userRepo.On("GetUserByUsername", mockString).Return(mocks.GetUser3(userData), nil).Once()

		expectedError := &domain.Error{
			Err:        errors.New("conflict"),
			StatusCode: http.StatusConflict,
			Message:    "Username already exists",
		}

		foundUser, err := suite.userUsecase.RegisterUser(userData)
		suite.Nil(foundUser)
		suite.Equal(expectedError, err)
	})
}

// A test for the UserUsecase.LoginUser method.
func (suite *UserUsecaseSuite) Test_LoginUser() {
	// A testcase that tests the successful login of a user.
	suite.Run("LoginUser_Success", func() {
		user := mocks.GetNewUser()
		userData := mocks.GetAuthData(user)
		user.Password, _ = infrastructure.HashPassword(user.Password)

		suite.userRepo.On("GetUserByUsername", mockString).Return(user, nil).Once()

		token, err := suite.userUsecase.LoginUser(userData)
		suite.Nil(err)
		suite.NotEmpty(token)

	})

	// A testcase that tests the failure of logging in a user.
	suite.Run("LoginUser_Failure", func() {
		userData := mocks.GetAuthUserData()

		suite.userRepo.On("GetUserByUsername", mockString).Return(nil, errors.New("some error")).Once()

		expectedError := &domain.Error{
			Err:        errors.New("some error"),
			StatusCode: http.StatusNotFound,
			Message:    "Invalid username or password",
		}

		token, err := suite.userUsecase.LoginUser(userData)
		suite.Empty(token)
		suite.Equal(expectedError, err)
	})

	// A testcase where the password is incorrect.
	suite.Run("LoginUser_IncorrectPassword", func() {
		userData := mocks.GetAuthUserData()
		user := mocks.GetUser3(userData)
		userData.Password = "wrong_password"
		user.Password, _ = infrastructure.HashPassword(user.Password)

		suite.userRepo.On("GetUserByUsername", mockString).Return(user, nil).Once()

		expectedError := &domain.Error{
			Err:        errors.New("crypto/bcrypt: hashedPassword is not the hash of the given password"),
			StatusCode: http.StatusUnauthorized,
			Message:    "Invalid username or password",
		}

		token, err := suite.userUsecase.LoginUser(userData)
		suite.Empty(token)
		suite.Equal(expectedError, err)
	})
}

// A test for the UserUsecase.GetUsers method.
func (suite *UserUsecaseSuite) Test_GetUsers() {
	// A testcase that tests the successful retrieval of all users.
	suite.Run("GetUsers_Success", func() {
		users := mocks.GetManyUsers()

		suite.userRepo.On("GetUsers").Return(users, nil).Once()

		foundUsers, err := suite.userUsecase.GetUsers()
		suite.Nil(err)
		suite.Equal(users, foundUsers)
	})

	// A testcase that tests the failure of getting all users.
	suite.Run("GetUsers_Failure", func() {
		suite.userRepo.On("GetUsers").Return(nil, errors.New("some error")).Once()

		expectedError := &domain.Error{
			Err:        errors.New("some error"),
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}

		foundUsers, err := suite.userUsecase.GetUsers()
		suite.Nil(foundUsers)
		suite.Equal(expectedError, err)
	})
}

// A test for the UserUsecase.GetUserByID method.
func (suite *UserUsecaseSuite) Test_GetUserByID() {
	// A testcase that tests the successful retrieval of a user by ID.
	suite.Run("GetUserByID_Success", func() {
		user := mocks.GetNewUser()

		suite.userRepo.On("GetUserByID", mockObjectID).Return(user, nil).Once()

		foundUser, err := suite.userUsecase.GetUserByID(user.ID)
		suite.Nil(err)
		suite.Equal(user, foundUser)
	})

	// A testcase that tests the failure of getting a user by ID.
	suite.Run("GetUserByID_Failure", func() {
		suite.userRepo.On("GetUserByID", mockObjectID).Return(nil, errors.New("some error")).Once()

		expectedError := &domain.Error{
			Err:        errors.New("some error"),
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}

		foundUser, err := suite.userUsecase.GetUserByID(mocks.GetPrimitiveID1())
		suite.Nil(foundUser)
		suite.Equal(expectedError, err)
	})

	// A testcase where the user is not found.
	suite.Run("GetUserByID_NotFound", func() {
		suite.userRepo.On("GetUserByID", mockObjectID).Return(nil, mongo.ErrNoDocuments).Once()

		expectedError := &domain.Error{
			Err:        mongo.ErrNoDocuments,
			StatusCode: http.StatusNotFound,
			Message:    "User not found",
		}

		foundUser, err := suite.userUsecase.GetUserByID(mocks.GetPrimitiveID1())
		suite.Nil(foundUser)
		suite.Equal(expectedError, err)
	})
}

// A test for the UserUsecase.UpdateUser method.
func (suite *UserUsecaseSuite) Test_UpdateUser() {
	// A testcase that tests the successful update of a user.
	suite.Run("UpdateUser_Success", func() {
		userData := mocks.GetUpdateUserData()
		user := mocks.GetUser4(userData)
		claims := mocks.GetClaims2() // An admin user.

		suite.userRepo.On("GetUserByID", mockObjectID).Return(user, nil).Twice()
		suite.userRepo.On("GetUserByUsername", mockString).Return(nil, errors.New("some error")).Once()
		suite.userRepo.On("UpdateUser", mockObjectID, mockBSON).Return(nil).Once()

		foundUser, err := suite.userUsecase.UpdateUser(user.ID, userData, claims)
		suite.Nil(err)
		suite.Equal(user, foundUser)
	})

	// A testcase where the user is not found.
	suite.Run("UpdateUser_NotFound", func() {
		userData := mocks.GetUpdateUserData()
		claims := mocks.GetClaims2() // An admin user.

		suite.userRepo.On("GetUserByID", mockObjectID).Return(nil, errors.New("some error")).Once()

		expectedError := &domain.Error{
			Err:        errors.New("some error"),
			StatusCode: http.StatusNotFound,
			Message:    "User not found",
		}

		foundUser, err := suite.userUsecase.UpdateUser(mocks.GetPrimitiveID1(), userData, claims)
		suite.Nil(foundUser)
		suite.Equal(expectedError, err)
	})

	// A testcase that tests the failure of updating a user.
	suite.Run("UpdateUser_Failure", func() {
		userData := mocks.GetUpdateUserData()
		user := mocks.GetUser4(userData)
		claims := mocks.GetClaims2() // An admin user.

		suite.userRepo.On("GetUserByID", mockObjectID).Return(user, nil).Once()
		suite.userRepo.On("GetUserByUsername", mockString).Return(nil, errors.New("some error")).Once()
		suite.userRepo.On("UpdateUser", mockObjectID, mockBSON).Return(errors.New("some error")).Once()

		expectedError := &domain.Error{
			Err:        errors.New("some error"),
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}

		foundUser, err := suite.userUsecase.UpdateUser(user.ID, userData, claims)
		suite.Nil(foundUser)
		suite.Equal(expectedError, err)
	})

	// A testcase where the user is not authorized.
	suite.Run("UpdateUser_Unauthorized", func() {
		userData := mocks.GetUpdateUserData()
		user := mocks.GetUser4(userData)
		user.Role = "root"
		claims := mocks.GetClaims2() // An admin user.
		objectID := primitive.NewObjectID()

		suite.userRepo.On("GetUserByID", mockObjectID).Return(user, nil).Once()

		expectedError := &domain.Error{
			Err:        errors.New("forbidden"),
			StatusCode: http.StatusForbidden,
			Message:    "Cannot update root user",
		}

		foundUser, err := suite.userUsecase.UpdateUser(objectID, userData, claims)
		suite.Nil(foundUser)
		suite.Equal(expectedError, err)
	})
}

// A test for the UserUsecase.DeleteUser method.
func (suite *UserUsecaseSuite) Test_DeleteUser() {
	// A testcase that tests the successful deletion of a user.
	suite.Run("DeleteUser_Success", func() {
		user := mocks.GetNewUser()
		claims := mocks.GetClaims2() // An admin user.

		suite.userRepo.On("GetUserByID", mockObjectID).Return(user, nil).Once()
		suite.userRepo.On("DeleteUser", mockObjectID).Return(nil).Once()

		err := suite.userUsecase.DeleteUser(user.ID, claims)
		suite.Nil(err)
	})

	// A testcase where the user is not found.
	suite.Run("DeleteUser_NotFound", func() {
		claims := mocks.GetClaims2() // An admin user.

		suite.userRepo.On("GetUserByID", mockObjectID).Return(nil, errors.New("some error")).Once()

		expectedError := &domain.Error{
			Err:        errors.New("some error"),
			StatusCode: http.StatusNotFound,
			Message:    "User not found",
		}

		err := suite.userUsecase.DeleteUser(mocks.GetPrimitiveID1(), claims)
		suite.Equal(expectedError, err)
	})

	// A testcase that tests the failure of deleting a user.
	suite.Run("DeleteUser_Failure", func() {
		user := mocks.GetNewUser()
		claims := mocks.GetClaims2() // An admin user.

		suite.userRepo.On("GetUserByID", mockObjectID).Return(user, nil).Once()
		suite.userRepo.On("DeleteUser", mockObjectID).Return(errors.New("some error")).Once()

		expectedError := &domain.Error{
			Err:        errors.New("some error"),
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}

		err := suite.userUsecase.DeleteUser(user.ID, claims)
		suite.Equal(expectedError, err)
	})

	// A testcase where the user is not authorized.
	suite.Run("DeleteUser_Unauthorized", func() {
		user := mocks.GetNewUser()
		claims := mocks.GetClaims() // A regular user.

		suite.userRepo.On("GetUserByID", mockObjectID).Return(user, nil).Once()

		expectedError := &domain.Error{
			Err:        errors.New("unauthorized"),
			StatusCode: http.StatusForbidden,
			Message:    "A User cannot delete another user",
		}

		err := suite.userUsecase.DeleteUser(user.ID, claims)
		suite.Equal(expectedError, err)
	})
}

// A function that runs the TestSuite.
func Test_UsecaseSuite(t *testing.T) {
	suite.Run(t, new(UserUsecaseSuite))
}
