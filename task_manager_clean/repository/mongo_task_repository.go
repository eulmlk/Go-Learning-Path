package repository

import (
	"context"
	"task_manager/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// This struct is a MongoDB implementation of the TaskRepository interface.
type MongoTaskRepository struct {
	collection *mongo.Collection
}

// A constructor that creates a new instance of MongoTaskRepository.
func NewMongoTaskRepository(collection *mongo.Collection) *MongoTaskRepository {
	return &MongoTaskRepository{
		collection: collection,
	}
}

// A method that returns all tasks.
func (r *MongoTaskRepository) GetAllTasks() ([]domain.Task, error) {
	tasks := []domain.Task{}

	// Query the database for all tasks.
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	// Iterate over the cursor and decode each task into a Task struct.
	err = cursor.All(context.Background(), &tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// A method that returns a task with the given ID.
func (r *MongoTaskRepository) GetTaskByID(id primitive.ObjectID) (*domain.Task, error) {
	task := &domain.Task{}

	// Query the database for a task with the given ID.
	result := r.collection.FindOne(context.Background(), bson.M{"_id": id})
	if err := result.Decode(task); err != nil {
		return nil, err
	}

	return task, nil
}

// A method that adds a new task.
func (r *MongoTaskRepository) AddTask(task *domain.Task) error {
	// Insert the task into the database.
	_, err := r.collection.InsertOne(context.Background(), task)
	return err
}

// A method that replaces a task with the given ID, with the new task.
func (r *MongoTaskRepository) ReplaceTask(id primitive.ObjectID, newTask *domain.Task) (*domain.Task, error) {
	newTask.ID = id
	result := r.collection.FindOneAndReplace(context.Background(), bson.M{"_id": id}, newTask)
	if err := result.Decode(newTask); err != nil {
		return nil, err
	}

	return newTask, nil
}

// A method that updates a task with the given ID.
func (r *MongoTaskRepository) UpdateTask(id primitive.ObjectID, taskData bson.M) (*domain.Task, error) {
	// Update the task with the given ID.
	_, err := r.collection.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": taskData})
	if err != nil {
		return nil, err
	}

	// Return the updated task.
	updatedTask := &domain.Task{}
	err = r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(updatedTask)
	return updatedTask, err
}

// A method that deletes a task with the given ID.
func (r *MongoTaskRepository) DeleteTask(id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}
