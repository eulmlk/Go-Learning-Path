package repository

import (
	"context"
	"task_manager/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// This struct is a MongoDB implementation of the UserRepository interface.
type MongoUserRepository struct {
	collection *mongo.Collection
}

// A constructor that creates a new instance of MongoUserRepository.
func NewMongoUserRepository(collection *mongo.Collection) *MongoUserRepository {
	return &MongoUserRepository{
		collection: collection,
	}
}

// A method that adds a new user.
func (r *MongoUserRepository) AddUser(user *domain.User) error {
	// Insert the user into the database.
	_, err := r.collection.InsertOne(context.Background(), user)
	return err
}

// A method that returns all users.
func (r *MongoUserRepository) GetUsers() ([]domain.User, error) {
	users := []domain.User{}

	// Query the database for all users.
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	// Iterate over the cursor and decode each user into a User struct.
	err = cursor.All(context.Background(), &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// A method that returns a user with the given id.
func (r *MongoUserRepository) GetUserByID(id primitive.ObjectID) (*domain.User, error) {
	user := &domain.User{}

	// Query the database for a user with the given ID.
	result := r.collection.FindOne(context.Background(), id)
	if err := result.Decode(user); err != nil {
		return nil, err
	}

	return user, nil
}

// A method that returns a user with the given username.
func (r *MongoUserRepository) GetUserByUsername(username string) (*domain.User, error) {
	user := &domain.User{}

	// Query the database for a user with the given username.
	result := r.collection.FindOne(context.Background(), bson.M{"username": username})
	if err := result.Decode(user); err != nil {
		return nil, err
	}

	return user, nil
}

// A method that updates a user with the given ID.
func (r *MongoUserRepository) UpdateUser(id primitive.ObjectID, userData bson.M) (*domain.User, error) {
	// Update the user in the database.
	_, err := r.collection.UpdateOne(context.Background(), id, bson.M{"$set": userData})
	if err != nil {
		return nil, err
	}

	// Get the updated user from the database.
	updatedUser := &domain.User{}
	err = r.collection.FindOne(context.Background(), id).Decode(updatedUser)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

// A method that deletes a user with the given ID.
func (r *MongoUserRepository) DeleteUser(id primitive.ObjectID) error {
	// Delete the user from the database.
	_, err := r.collection.DeleteOne(context.Background(), id)
	return err
}
