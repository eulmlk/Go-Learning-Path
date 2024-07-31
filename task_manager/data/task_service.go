package data

import (
	"errors"
	"math/rand"
	"task_manager/models"
)

// TaskService is a struct that holds a map of tasks.
type TaskService struct {
	tasks map[string]*models.Task
}

// NewTaskService creates a new instance of TaskService.
// It initializes the tasks map and returns a pointer to the TaskService.
func NewTaskService() *TaskService {
	return &TaskService{tasks: make(map[string]*models.Task)}
}

// GetTasks returns a slice of all tasks.
func (ts *TaskService) GetTasks() []models.Task {
	tasks := []models.Task{}

	// Loop through the tasks map and append each task to the tasks slice.
	for _, task := range ts.tasks {
		tasks = append(tasks, *task)
	}

	return tasks
}

// GetTaskByID returns the task with the given ID.
func (ts *TaskService) GetTaskByID(id string) *models.Task {
	// Check if the task with the given ID exists in the tasks map.
	// If it exists, return the task, otherwise return nil.
	if task, ok := ts.tasks[id]; ok {
		return task
	}

	return nil
}

// CreateTask creates a new task and adds it to the tasks map.
func (ts *TaskService) CreateTask(task *models.Task) {
	ts.tasks[task.ID] = task
}

// randomString generates a random string of the given length.
func (ts *TaskService) randomString(length int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, length)

	// Generate a random string of the given length using the letterRunes slice.
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}

// GenerateID generates a random ID for a new task.
func (ts *TaskService) GenerateID() string {
	var id string

	// Generate a random ID and check if it already exists in the tasks map.
	// If it exists, generate a new ID until a unique one is found.
	for {
		id = ts.randomString(30)
		if ts.GetTaskByID(id) == nil {
			break
		}
	}

	return id
}

// UpdateTask updates the task with the given ID using the provided task data.
func (ts *TaskService) UpdateTask(id string, taskData *models.Task) *models.Task {
	_, ok := ts.tasks[id]
	if !ok {
		return nil
	}

	// Update the task fields with the provided task data.
	// If a field is not provided, the existing value is retained.
	if taskData.Title != "" {
		ts.tasks[id].Title = taskData.Title
	}

	if taskData.Description != "" {
		ts.tasks[id].Description = taskData.Description
	}

	if !taskData.DueDate.IsZero() {
		ts.tasks[id].DueDate = taskData.DueDate
	}

	if taskData.Status != "" {
		ts.tasks[id].Status = taskData.Status
	}

	return ts.tasks[id]
}

// DeleteTask deletes the task with the given ID.
func (ts *TaskService) DeleteTask(id string) error {
	// Check if the task with the given ID exists in the tasks map.
	// If it exists, delete the task, otherwise return an error.

	if _, ok := ts.tasks[id]; ok {
		delete(ts.tasks, id)
		return nil
	}

	return errors.New("task not found")
}
