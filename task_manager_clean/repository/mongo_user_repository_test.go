package repository_test

import (
	"context"
	"task_manager/database"
	"task_manager/domain"
	"task_manager/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	userCollection = "users_test"
)

// A test for the MongoUserRepository.AddUser method.
func TestMongoUserRepository_AddUser(t *testing.T) {
	client, collection := setupUserCollection(t)
	defer client.Disconnect(context.Background())

	repo := repository.NewMongoUserRepository(collection)
	user := getNewUser()

	// A test case for adding a new user
	t.Run("AddUser", func(t *testing.T) {
		// 1. Check the number of users in the collection
		count, err := collection.CountDocuments(context.Background(), bson.M{})
		require.NoError(t, err)
		require.Equal(t, int64(0), count)

		// 2. Add a new user
		err = repo.AddUser(user)
		require.NoError(t, err)

		// 3. Check the number of users in the collection again
		count, err = collection.CountDocuments(context.Background(), bson.M{})
		require.NoError(t, err)
		require.Equal(t, int64(1), count)

		// 4. Check if the user was added correctly
		addedUser, err := repo.GetUserByID(user.ID)
		require.NoError(t, err)
		require.Equal(t, *user, *addedUser)
	})
}

// A test for the MongoUserRepository.GetUsers method.
func TestMongoUserRepository_GetUsers(t *testing.T) {
	client, collection := setupUserCollection(t)
	defer client.Disconnect(context.Background())

	repo := repository.NewMongoUserRepository(collection)
	user1 := getNewUser()
	user2 := getNewUser2()

	// A test case for an empty collection
	t.Run("GetUsers_Empty", func(t *testing.T) {
		users, err := repo.GetUsers()
		require.NoError(t, err)
		require.Len(t, users, 0)
	})

	_, err := collection.InsertMany(context.Background(), []interface{}{user1, user2})
	require.NoError(t, err)

	// A test case for a collection with two users
	t.Run("GetUsers_NotEmpty", func(t *testing.T) {
		users, err := repo.GetUsers()
		require.NoError(t, err)
		require.Len(t, users, 2)
	})
}

// A test for the MongoUserRepository.GetUserByID method.
func TestMongoUserRepository_GetUserByID(t *testing.T) {
	client, collection := setupUserCollection(t)
	defer client.Disconnect(context.Background())

	repo := repository.NewMongoUserRepository(collection)
	user := getNewUser()

	// A test case for a user not found
	t.Run("GetUserByID_NotFound", func(t *testing.T) {
		foundUser, err := repo.GetUserByID(user.ID)
		require.Error(t, err)
		assert.Nil(t, foundUser)
	})

	_, err := collection.InsertOne(context.Background(), user)
	require.NoError(t, err)

	// A test case for a user found
	t.Run("GetUserByID_Found", func(t *testing.T) {
		added, err := repo.GetUserByID(user.ID)
		require.NoError(t, err)
		require.Equal(t, user, added)
	})
}

// A test for the MongoUserRepository.DeleteUser method.
func TestMongoUserRepository_GetUserByUsername(t *testing.T) {
	client, collection := setupUserCollection(t)
	defer client.Disconnect(context.Background())

	repo := repository.NewMongoUserRepository(collection)
	user := getNewUser()

	// A test case for a user not found
	t.Run("GetUserByUsername_NotFound", func(t *testing.T) {
		foundUser, err := repo.GetUserByUsername(user.Username)
		require.Error(t, err)
		assert.Nil(t, foundUser)
	})

	_, err := collection.InsertOne(context.Background(), user)
	require.NoError(t, err)

	// A test case for a user found
	t.Run("GetUserByUsername_Found", func(t *testing.T) {
		userFromDB, err := repo.GetUserByUsername(user.Username)
		require.NoError(t, err)
		require.Equal(t, user, userFromDB)
	})
}

func TestMongoUserRepository_UpdateUser(t *testing.T) {
	client, collection := setupUserCollection(t)
	defer client.Disconnect(context.Background())

	repo := repository.NewMongoUserRepository(collection)
	user := getNewUser()

	_, err := collection.InsertOne(context.Background(), user)
	require.NoError(t, err)

	// A test case for a user found
	t.Run("UpdateUser_Found", func(t *testing.T) {
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

		err := repo.UpdateUser(user.ID, update)
		require.NoError(t, err)

		userFromDB := &domain.User{}
		err = collection.FindOne(context.Background(), bson.M{"_id": user.ID}).Decode(userFromDB)
		require.NoError(t, err)
		assert.Equal(t, expectedUser, userFromDB)
	})
}

// A test for the MongoUserRepository.DeleteUser method.
func TestMongoUserRepository_DeleteUser(t *testing.T) {
	client, collection := setupUserCollection(t)
	defer client.Disconnect(context.Background())

	repo := repository.NewMongoUserRepository(collection)
	user := getNewUser()

	_, err := collection.InsertOne(context.Background(), user)
	require.NoError(t, err)

	// A test case for a user found
	t.Run("DeleteUser_Found", func(t *testing.T) {
		err := repo.DeleteUser(user.ID)
		require.NoError(t, err)

		count, err := collection.CountDocuments(context.Background(), bson.M{})
		require.NoError(t, err)
		require.Equal(t, int64(0), count)
	})
}

func setupUserCollection(t *testing.T) (*mongo.Client, *mongo.Collection) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	require.NoError(t, err)
	require.NoError(t, client.Ping(context.Background(), nil))

	db := client.Database(database.DatabaseName)
	collection := db.Collection(userCollection)
	require.NoError(t, collection.Drop(context.Background()))

	return client, collection
}

func getNewUser() *domain.User {
	return &domain.User{
		ID:       primitive.NewObjectID(),
		Username: "user1",
		Password: "password1",
		Role:     "user",
	}
}

func getNewUser2() *domain.User {
	return &domain.User{
		ID:       primitive.NewObjectID(),
		Username: "user2",
		Password: "password2",
		Role:     "user",
	}
}
