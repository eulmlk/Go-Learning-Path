package repository_test

import (
	"context"
	"task_manager/database"
	"task_manager/domain"
	"task_manager/repository"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	userCollection = "users_test"
)

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

func TestMongoUserRepository_AddUser(t *testing.T) {
	client, collection := setupUserCollection(t)
	defer client.Disconnect(context.Background())

	repo := repository.NewMongoUserRepository(collection)

	user := &domain.User{
		ID:       primitive.NewObjectID(),
		Username: "user1",
		Password: "password1",
		Role:     "user",
	}

	err := repo.AddUser(user)
	require.NoError(t, err)

	userFromDB, err := repo.GetUserByID(user.ID)
	require.NoError(t, err)
	require.Equal(t, user, userFromDB)
}

func TestMongoUserRepository_GetUsers(t *testing.T) {
	client, collection := setupUserCollection(t)
	defer client.Disconnect(context.Background())

	repo := repository.NewMongoUserRepository(collection)

	user1 := &domain.User{
		ID:       primitive.NewObjectID(),
		Username: "user1",
		Password: "password1",
		Role:     "user",
	}

	user2 := &domain.User{
		ID:       primitive.NewObjectID(),
		Username: "user2",
		Password: "password2",
		Role:     "user",
	}

	err := repo.AddUser(user1)
	require.NoError(t, err)

	err = repo.AddUser(user2)
	require.NoError(t, err)

	users, err := repo.GetUsers()
	require.NoError(t, err)
	require.Len(t, users, 2)
}

func TestMongoUserRepository_GetUserByID(t *testing.T) {
	client, collection := setupUserCollection(t)
	defer client.Disconnect(context.Background())

	repo := repository.NewMongoUserRepository(collection)

	user := &domain.User{
		ID:       primitive.NewObjectID(),
		Username: "user1",
		Password: "password1",
		Role:     "user",
	}

	err := repo.AddUser(user)
	require.NoError(t, err)

	userFromDB, err := repo.GetUserByID(user.ID)
	require.NoError(t, err)
	require.Equal(t, user, userFromDB)
}

func TestMongoUserRepository_GetUserByUsername(t *testing.T) {
	client, collection := setupUserCollection(t)
	defer client.Disconnect(context.Background())

	repo := repository.NewMongoUserRepository(collection)

	user := &domain.User{
		ID:       primitive.NewObjectID(),
		Username: "user1",
		Password: "password1",
		Role:     "user",
	}

	err := repo.AddUser(user)
	require.NoError(t, err)

	userFromDB, err := repo.GetUserByUsername(user.Username)
	require.NoError(t, err)
	require.Equal(t, user, userFromDB)
}

func TestMongoUserRepository_UpdateUser(t *testing.T) {
	client, collection := setupUserCollection(t)
	defer client.Disconnect(context.Background())

	repo := repository.NewMongoUserRepository(collection)

	user := &domain.User{
		ID:       primitive.NewObjectID(),
		Username: "user1",
		Password: "password1",
		Role:     "user",
	}

	err := repo.AddUser(user)
	require.NoError(t, err)

	update := bson.M{
		"$set": bson.M{
			"username": "user2",
			"password": "password2",
			"role":     "admin",
		},
	}

	updatedUser, err := repo.UpdateUser(user.ID, update)
	require.NoError(t, err)

	userFromDB, err := repo.GetUserByID(user.ID)
	require.NoError(t, err)
	require.Equal(t, updatedUser, userFromDB)
}
