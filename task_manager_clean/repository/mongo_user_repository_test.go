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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	userCollection = "users_test"
)

type MongoUserRepositorySuite struct {
	suite.Suite
	repo       *repository.MongoUserRepository
	collection *mongo.Collection
	client     *mongo.Client
}

func (suite *MongoUserRepositorySuite) SetupSuite() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	suite.NoError(err)
	suite.NoError(client.Ping(context.Background(), nil))
	suite.client = client

	db := client.Database(database.DatabaseName)
	suite.collection = db.Collection(userCollection)
	suite.NoError(suite.collection.Drop(context.Background()))

	suite.repo = repository.NewMongoUserRepository(suite.collection)
}

func (suite *MongoUserRepositorySuite) SetupTest() {
	err := suite.collection.Drop(context.Background())
	suite.NoError(err)
}

func (suite *MongoUserRepositorySuite) TearDownSuite() {
	err := suite.collection.Drop(context.Background())
	suite.NoError(err)
	err = suite.client.Disconnect(context.Background())
	suite.NoError(err)
}

func (suite *MongoUserRepositorySuite) TestAddUser() {
	user := mocks.GetNewUser()

	err := suite.repo.AddUser(user)
	suite.NoError(err)

	count, err := suite.collection.CountDocuments(context.Background(), bson.M{})
	suite.NoError(err)
	suite.Equal(int64(1), count)

	addedUser, err := suite.repo.GetUserByID(user.ID)
	suite.NoError(err)
	suite.Equal(*user, *addedUser)
}

func (suite *MongoUserRepositorySuite) TestGetUsers() {
	user1 := mocks.GetNewUser()
	user2 := mocks.GetNewUser2()

	users, err := suite.repo.GetUsers()
	suite.NoError(err)
	suite.Len(users, 0)

	_, err = suite.collection.InsertMany(context.Background(), []interface{}{user1, user2})
	suite.NoError(err)

	users, err = suite.repo.GetUsers()
	suite.NoError(err)
	suite.Len(users, 2)
}

func (suite *MongoUserRepositorySuite) TestGetUserByID() {
	user := mocks.GetNewUser()

	foundUser, err := suite.repo.GetUserByID(user.ID)
	suite.Error(err)
	suite.Nil(foundUser)

	_, err = suite.collection.InsertOne(context.Background(), user)
	suite.NoError(err)

	added, err := suite.repo.GetUserByID(user.ID)
	suite.NoError(err)
	suite.Equal(user, added)
}

func (suite *MongoUserRepositorySuite) TestGetUserByUsername() {
	user := mocks.GetNewUser()

	foundUser, err := suite.repo.GetUserByUsername(user.Username)
	suite.Error(err)
	suite.Nil(foundUser)

	_, err = suite.collection.InsertOne(context.Background(), user)
	suite.NoError(err)

	userFromDB, err := suite.repo.GetUserByUsername(user.Username)
	suite.NoError(err)
	suite.Equal(user, userFromDB)
}

func (suite *MongoUserRepositorySuite) TestUpdateUser() {
	user := mocks.GetNewUser()

	_, err := suite.collection.InsertOne(context.Background(), user)
	suite.NoError(err)

	update := bson.M{
		"username": "new_username",
		"password": "new_password",
	}

	expectedUser := &domain.User{
		ID:       user.ID,
		Username: "new_username",
		Password: "new_password",
		Role:     user.Role,
	}

	err = suite.repo.UpdateUser(user.ID, update)
	suite.NoError(err)

	userFromDB := &domain.User{}
	err = suite.collection.FindOne(context.Background(), bson.M{"_id": user.ID}).Decode(userFromDB)
	suite.NoError(err)
	suite.Equal(expectedUser, userFromDB)
}

func (suite *MongoUserRepositorySuite) TestDeleteUser() {
	user := mocks.GetNewUser()

	_, err := suite.collection.InsertOne(context.Background(), user)
	suite.NoError(err)

	err = suite.repo.DeleteUser(user.ID)
	suite.NoError(err)

	count, err := suite.collection.CountDocuments(context.Background(), bson.M{})
	suite.NoError(err)
	suite.Equal(int64(0), count)
}

func TestMongoUserRepositorySuite(t *testing.T) {
	suite.Run(t, new(MongoUserRepositorySuite))
}
