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

var mockUser = mock.AnythingOfType("*domain.User")

// A suite that contains tests for the MongoTaskRepository.
type MongoUserRepositoryTestSuite struct {
	suite.Suite
	repo       *repository.MongoUserRepository
	collection *mocks.Collection
}

// A method that initializes the test suite.
func (suite *MongoUserRepositoryTestSuite) SetupSuite() {
	suite.collection = new(mocks.Collection)
	suite.repo = repository.NewMongoUserRepository(suite.collection)
}

// A method that finalizes the test suite.
func (suite *MongoUserRepositoryTestSuite) TearDownSuite() {
	suite.collection.AssertExpectations(suite.T())
}

// A test for the MongoUserRepository.AddUser method.
func (suite *MongoUserRepositoryTestSuite) TestAddUser() {
	// A testcase for the successful addition of a user.
	suite.Run("AddUser_Success", func() {
		user := mocks.GetNewUser()
		suite.collection.On("InsertOne", mock.Anything, user).Return(&mongo.InsertOneResult{}, nil).Once()

		err := suite.repo.AddUser(user)
		suite.NoError(err)
	})

	// A testcase for the failure of adding a user.
	suite.Run("AddUser_Failure", func() {
		user := mocks.GetNewUser()
		suite.collection.On("InsertOne", mock.Anything, user).Return(&mongo.InsertOneResult{}, mongo.ErrClientDisconnected).Once()

		err := suite.repo.AddUser(user)
		suite.Error(err)
	})
}

// A test for the MongoUserRepository.GetUsers method.
func (suite *MongoUserRepositoryTestSuite) TestGetUsers() {
	// A testcase for the successful retrieval of users.
	suite.Run("GetUsers_Success", func() {
		users := mocks.GetManyUsers()
		cursor := new(mocks.Cursor)

		suite.collection.On("Find", mock.Anything, mock.Anything).Return(cursor, nil).Once()
		cursor.On("All", mock.Anything, &[]domain.User{}).Return(nil).Once().Run(func(args mock.Arguments) {
			usersPtr := args.Get(1).(*[]domain.User)
			*usersPtr = append(*usersPtr, users...)
		})

		result, err := suite.repo.GetUsers()

		suite.NoError(err)
		suite.Equal(users, result)
	})

	// A testcase for the failure of retrieving users.
	suite.Run("GetUsers_Failure", func() {
		suite.collection.On("Find", mock.Anything, mock.Anything).Return(nil, mongo.ErrClientDisconnected).Once()

		result, err := suite.repo.GetUsers()
		suite.Error(err)
		suite.Nil(result)
	})
}

// A test for the MongoUserRepository.GetUserByID method.
func (suite *MongoUserRepositoryTestSuite) TestGetUserByID() {
	// A testcase for the successful retrieval of a user by ID.
	suite.Run("GetUserByID_Success", func() {
		user := mocks.GetNewUser()
		id := user.ID

		res := new(mocks.SingleResult)
		res.On("Decode", mock.Anything).Return(nil).Once().RunFn = func(args mock.Arguments) {
			userPtr := args.Get(0).(*domain.User)
			*userPtr = *user
		}

		suite.collection.On("FindOne", mock.Anything, mock.Anything).Return(res).Once()

		result, err := suite.repo.GetUserByID(id)

		suite.NoError(err)
		suite.Equal(user, result)
	})

	// A testcase for the failure of retrieving a user by ID.
	suite.Run("GetUserByID_Failure", func() {
		id := primitive.NewObjectID()

		res := new(mocks.SingleResult)
		res.On("Decode", mockUser).Return(mongo.ErrNoDocuments).Once()
		suite.collection.On("FindOne", mock.Anything, mock.Anything).Return(res).Once()

		result, err := suite.repo.GetUserByID(id)
		suite.Error(err)
		suite.Nil(result)
	})
}

// A test for the MongoUserRepository.GetUserByUsername method.
func (suite *MongoUserRepositoryTestSuite) TestGetUserByUsername() {
	// A testcase for the successful retrieval of a user by username.
	suite.Run("GetUserByUsername_Success", func() {
		user := mocks.GetNewUser()
		username := user.Username

		res := new(mocks.SingleResult)
		res.On("Decode", mock.Anything).Return(nil).Once().Run(func(args mock.Arguments) {
			userPtr := args.Get(0).(*domain.User)
			*userPtr = *user
		})

		suite.collection.On("FindOne", mock.Anything, mock.Anything).Return(res).Once()

		result, err := suite.repo.GetUserByUsername(username)

		suite.NoError(err)
		suite.Equal(user, result)
	})

	// A testcase for the failure of retrieving a user by username.
	suite.Run("GetUserByUsername_Failure", func() {
		username := "nonexistent"

		res := new(mocks.SingleResult)
		res.On("Decode", mockUser).Return(mongo.ErrNoDocuments).Once()
		suite.collection.On("FindOne", mock.Anything, mock.Anything).Return(res).Once()

		result, err := suite.repo.GetUserByUsername(username)
		suite.Error(err)
		suite.Nil(result)
	})
}

// A test for the MongoUserRepository.UpdateUser method.
func (suite *MongoUserRepositoryTestSuite) TestUpdateUser() {
	// A testcase for the successful update of a user.
	suite.Run("UpdateUser_Success", func() {
		user := mocks.GetNewUser()
		id := user.ID
		userData := bson.M{"username": "new_username"}

		suite.collection.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).Return(&mongo.UpdateResult{}, nil).Once()

		err := suite.repo.UpdateUser(id, userData)
		suite.NoError(err)
	})

	// A testcase for the failure of updating a user.
	suite.Run("UpdateUser_Failure", func() {
		user := mocks.GetNewUser()
		id := user.ID
		userData := bson.M{"username": "new_username"}

		suite.collection.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).Return(&mongo.UpdateResult{}, mongo.ErrClientDisconnected).Once()

		err := suite.repo.UpdateUser(id, userData)
		suite.Error(err)
	})
}

// A test for the MongoUserRepository.DeleteUser method.
func (suite *MongoUserRepositoryTestSuite) TestDeleteUser() {
	// A testcase for the successful deletion of a user.
	suite.Run("DeleteUser_Success", func() {
		user := mocks.GetNewUser()
		id := user.ID

		suite.collection.On("DeleteOne", mock.Anything, mock.Anything).Return(&mongo.DeleteResult{}, nil).Once()

		err := suite.repo.DeleteUser(id)
		suite.NoError(err)
	})

	// A testcase for the failure of deleting a user.
	suite.Run("DeleteUser_Failure", func() {
		user := mocks.GetNewUser()
		id := user.ID

		suite.collection.On("DeleteOne", mock.Anything, mock.Anything).Return(&mongo.DeleteResult{}, mongo.ErrClientDisconnected).Once()

		err := suite.repo.DeleteUser(id)
		suite.Error(err)
	})
}

// A function that runs the MongoUserRepositoryTestSuite.
func TestMongoUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(MongoUserRepositoryTestSuite))
}
