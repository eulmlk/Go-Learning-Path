package data

import (
	"errors"
	"fmt"
	"math/rand"
	"task_manager/models"
)

type TaskService struct {
	tasks map[string]*models.Task
}

func NewTaskService() *TaskService {
	return &TaskService{tasks: make(map[string]*models.Task)}
}

func (ts *TaskService) GetTasks() []models.Task {
	fmt.Println(ts.tasks)

	tasks := []models.Task{}

	for _, task := range ts.tasks {
		tasks = append(tasks, *task)
	}

	return tasks
}

func (ts *TaskService) GetTaskByID(id string) *models.Task {
	if task, ok := ts.tasks[id]; ok {
		return task
	}

	return nil
}

func (ts *TaskService) CreateTask(task *models.Task) {
	ts.tasks[task.ID] = &models.Task{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		DueDate:     task.DueDate,
		Status:      task.Status,
	}
}

func (ts *TaskService) randomString(length int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, length)

	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}

func (ts *TaskService) GenerateID() string {
	var id string

	for {
		id = ts.randomString(30)
		if ts.GetTaskByID(id) == nil {
			break
		}
	}

	return id
}

func (ts *TaskService) UpdateTask(id string, taskData *models.Task) *models.Task {
	_, ok := ts.tasks[id]
	if !ok {
		return nil
	}

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

func (ts *TaskService) DeleteTask(id string) error {
	if _, ok := ts.tasks[id]; ok {
		delete(ts.tasks, id)
		return nil
	}

	return errors.New("task not found")
}
