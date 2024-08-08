package usecase

import (
	"errors"
	"net/http"
	"task_manager/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// A struct that defines the services for tasks.
type TaskUsecase struct {
	taskRepo domain.TaskRepository
}

// A constructor that creates a new instance of TaskUsecase.
func NewTaskUsecase(taskRepo domain.TaskRepository) *TaskUsecase {
	return &TaskUsecase{
		taskRepo: taskRepo,
	}
}

// A method that returns all tasks.
func (tu *TaskUsecase) GetTasks() ([]domain.Task, *domain.Error) {
	// Query the database for all tasks.
	tasks, err := tu.taskRepo.GetAllTasks()
	if err != nil {
		return nil, &domain.Error{
			Err:        err,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}

	return tasks, nil
}

// A method that returns a task with the given ID.
func (tu *TaskUsecase) GetTaskByID(objectID primitive.ObjectID) (*domain.Task, *domain.Error) {
	task, err := tu.taskRepo.GetTaskByID(objectID)
	if err != nil {
		// Check if the task is not found.
		if err == mongo.ErrNoDocuments {
			return nil, &domain.Error{
				Err:        err,
				StatusCode: http.StatusNotFound,
				Message:    "Task not found",
			}
		}

		return nil, &domain.Error{
			Err:        err,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}

	return task, nil
}

// A method that creates a new task.
func (tu *TaskUsecase) CreateTask(taskData *domain.CreateTaskData, claims *domain.Claims) (*domain.TaskView, *domain.Error) {
	// Create the task object.
	task := &domain.Task{
		ID:          primitive.NewObjectID(),
		Title:       taskData.Title,
		Description: taskData.Description,
		DueDate:     taskData.DueDate,
		Status:      taskData.Status,
		UserID:      claims.ID,
	}

	// Insert the task into the database.
	err := tu.taskRepo.AddTask(task)
	if err != nil {
		return nil, &domain.Error{
			Err:        err,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}

	// Create the task view object.
	taskView := &domain.TaskView{
		ID:          task.ID.Hex(),
		Title:       task.Title,
		Description: task.Description,
		DueDate:     task.DueDate,
		Status:      task.Status,
	}

	return taskView, nil
}

// A method that fully replaces a task with the given ID with the new task data.
func (tu *TaskUsecase) ReplaceTask(objectID primitive.ObjectID, taskData *domain.ReplaceTaskData, claims *domain.Claims) (*domain.TaskView, *domain.Error) {
	// Check if the task exists.
	_, _err := tu.GetTaskByID(objectID)
	if _err != nil {
		return nil, _err
	}

	// Check if the user is an admin or the owner of the task.
	if claims.Role == "user" && claims.ID != objectID {
		return nil, &domain.Error{
			Err:        errors.New("trying to replace another user's task"),
			StatusCode: http.StatusForbidden,
			Message:    "A User can only update their own task",
		}
	}

	// Create the new task object.
	task := &domain.Task{
		ID:          objectID,
		Title:       taskData.Title,
		Description: taskData.Description,
		DueDate:     taskData.DueDate,
		Status:      taskData.Status,
		UserID:      claims.ID,
	}

	// Replace the task in the database.
	err := tu.taskRepo.ReplaceTask(objectID, task)
	if err != nil {
		return nil, &domain.Error{
			Err:        err,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}

	newTask, err := tu.taskRepo.GetTaskByID(objectID)
	if err != nil {
		return nil, &domain.Error{
			Err:        err,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}

	// Create the task view object.
	taskView := &domain.TaskView{
		ID:          newTask.ID.Hex(),
		Title:       newTask.Title,
		Description: newTask.Description,
		DueDate:     newTask.DueDate,
		Status:      newTask.Status,
	}

	return taskView, nil
}

// A method that partially updates a task with the given ID with the only the provided task data.
func (tu *TaskUsecase) UpdateTask(objectID primitive.ObjectID, taskData *domain.UpdateTaskData, claims *domain.Claims) (*domain.Task, *domain.Error) {
	// Check if the task exists.
	_, _err := tu.GetTaskByID(objectID)
	if _err != nil {
		return nil, _err
	}

	// Check if the user is an admin or the owner of the task.
	if claims.Role == "user" && claims.ID != objectID {
		return nil, &domain.Error{
			Err:        errors.New("trying to update another user's task"),
			StatusCode: http.StatusForbidden,
			Message:    "A User can only update their own task",
		}
	}

	// Get the data to update.
	updateData := bson.M{}
	if taskData.Title != "" {
		updateData["title"] = taskData.Title
	}
	if taskData.Description != "" {
		updateData["description"] = taskData.Description
	}
	if taskData.Status != "" {
		updateData["status"] = taskData.Status
	}
	if !taskData.DueDate.IsZero() {
		updateData["due_date"] = taskData.DueDate
	}

	// Update the task in the database.
	err := tu.taskRepo.UpdateTask(objectID, updateData)
	if err != nil {
		return nil, &domain.Error{
			Err:        err,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}

	task, err := tu.taskRepo.GetTaskByID(objectID)
	if err != nil {
		return nil, &domain.Error{
			Err:        err,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}

	return task, nil
}

func (tu *TaskUsecase) DeleteTask(objectID primitive.ObjectID, claims *domain.Claims) *domain.Error {
	_, _err := tu.GetTaskByID(objectID)
	if _err != nil {
		return _err
	}

	// Check if the user is an admin or the owner of the task.
	if claims.Role == "user" && claims.ID != objectID {
		return &domain.Error{
			Err:        errors.New("trying to delete another user's task"),
			StatusCode: http.StatusForbidden,
			Message:    "A User can only delete their own task",
		}
	}

	// Delete the task from the database.
	err := tu.taskRepo.DeleteTask(objectID)
	if err != nil {
		return &domain.Error{
			Err:        err,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}

	return nil
}
