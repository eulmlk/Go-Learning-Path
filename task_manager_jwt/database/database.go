package database

import (
	"context"
	"errors"
	"log"
	"os"
	"task_manager/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func Init() (*mongo.Client, error) {
	// Get mongo connection uri
	uri, ok := os.LookupEnv("MONGODB_URI")
	if !ok {
		return nil, errors.New("MONGODB_URI is not set")
	}

	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB!")
	return client, nil
}

func CreateRootUser(client *mongo.Client) error {
	// Get the root username and password
	rootUsername, ok := os.LookupEnv("ROOT_USERNAME")
	if !ok {
		return errors.New("ROOT_USERNAME is not set")
	}

	rootPassword, ok := os.LookupEnv("ROOT_PASSWORD")
	if !ok {
		return errors.New("ROOT_PASSWORD is not set")
	}

	// Get the user collection
	collection := client.Database("task_manager").Collection("users")

	// Check if the root user exists
	var user map[string]interface{}
	err := collection.FindOne(context.Background(), map[string]string{"username": rootUsername}).Decode(&user)
	if err == nil {
		return nil
	}

	// Create the root user
	rootUser := models.User{
		ID:       primitive.NewObjectID(),
		Username: rootUsername,
		Password: rootPassword,
		Role:     "root",
	}

	// Hash the password
	bytes, err := bcrypt.GenerateFromPassword([]byte(rootUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	rootUser.Password = string(bytes)

	// Insert the root user
	_, err = collection.InsertOne(context.Background(), rootUser)
	if err != nil {
		return err
	}

	log.Println("Root user created!")
	return nil
}
