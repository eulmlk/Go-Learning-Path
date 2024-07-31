package controllers

import (
	"task_manager/data"
	"task_manager/models"
	"time"

	"github.com/gin-gonic/gin"
)

// TaskController is a struct that holds a pointer to the TaskService.
type TaskController struct {
	service *data.TaskService
}

// NewTaskController creates a new instance of TaskController.
func NewTaskController(service *data.TaskService) *TaskController {
	return &TaskController{service}
}

// GetTasks is a handler function that returns all tasks.
func (tc *TaskController) GetTasks(ctx *gin.Context) {
	tasks := tc.service.GetTasks()

	ctx.JSON(200, gin.H{
		"count": len(tasks),
		"tasks": tasks,
	})
}

// GetTaskByID is a handler function that returns a task by ID.
func (tc *TaskController) GetTaskByID(ctx *gin.Context) {
	taskID := ctx.Param("id")
	task := tc.service.GetTaskByID(taskID)

	// If the task is not found, return a 404 response.
	if task == nil {
		ctx.JSON(404, gin.H{"error": "Task not found"})
		return
	}

	// Otherwise, return a 200 response with the task.
	ctx.JSON(200, task)
}

func (tc *TaskController) CreateTask(ctx *gin.Context) {
	newTask := models.Task{}

	// Bind the request body to the newTask struct.
	err := ctx.BindJSON(&newTask)

	// If there is an error binding the request body, return a 400 response.
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request!"})
		return
	}

	// If the ID field is present in the request body, return a 400 response.
	if newTask.ID != "" {
		ctx.JSON(400, gin.H{"error": "id field is not allowed"})
		return
	}

	// If any of the required fields are missing, return a 400 response.
	if newTask.Title == "" {
		ctx.JSON(400, gin.H{"error": "title field is required"})
		return
	}

	if newTask.DueDate == (time.Time{}) {
		ctx.JSON(400, gin.H{"error": "due_date field is required"})
		return
	}

	// If the status field is missing, set it to "Pending" by default.
	if newTask.Status == "" {
		newTask.Status = "Pending"

		// If the status field is present, check if it is one of the allowed values.
	} else if newTask.Status != "Pending" && newTask.Status != "Completed" && newTask.Status != "In Progress" {
		ctx.JSON(400, gin.H{"error": "status field must be one of: Pending, Completed, In Progress"})
		return
	}

	// Generate a new ID for the task.
	newTask.ID = tc.service.GenerateID()

	// Create the task using the TaskService.
	tc.service.CreateTask(&newTask)

	// Return a 201 response with the new task.
	ctx.JSON(201, newTask)
}

func (tc *TaskController) UpdateTaskPut(ctx *gin.Context) {
	newTask := models.Task{}

	// Bind the request body to the newTask struct.
	err := ctx.BindJSON(&newTask)

	// If there is an error binding the request body, return a 400 response.
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// If any of the required fields are missing, return a 400 response.
	if newTask.Title == "" {
		ctx.JSON(400, gin.H{"error": "title field is required"})
		return
	}

	if newTask.Description == "" {
		ctx.JSON(400, gin.H{"error": "description field is required"})
		return
	}

	if newTask.DueDate == (time.Time{}) {
		ctx.JSON(400, gin.H{"error": "due_date field is required"})
		return
	}

	if newTask.Status == "" {
		ctx.JSON(400, gin.H{"error": "status field is required"})
		return
	} else if newTask.Status != "Pending" && newTask.Status != "Completed" && newTask.Status != "In Progress" {
		ctx.JSON(400, gin.H{"error": "status field must be one of: Pending, Completed, In Progress"})
		return
	}

	// Update the task using the TaskService.
	taskID := ctx.Param("id")
	task := tc.service.UpdateTask(taskID, &newTask)

	// If the task is not found, return a 404 response.
	if task == nil {
		ctx.JSON(404, gin.H{"error": "Task not found"})
		return
	}

	// Otherwise, return a 200 response with the updated task.
	ctx.JSON(200, task)
}

func (tc *TaskController) UpdateTaskPatch(ctx *gin.Context) {
	newTask := models.Task{}

	// Bind the request body to the newTask struct.
	err := ctx.BindJSON(&newTask)

	// If there is an error binding the request body, return a 400 response.
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// Update the task using the TaskService.
	taskID := ctx.Param("id")
	task := tc.service.UpdateTask(taskID, &newTask)

	// If the task is not found, return a 404 response.
	if task == nil {
		ctx.JSON(404, gin.H{"error": "Task not found"})
		return
	}

	// Otherwise, return a 200 response with the updated task.
	ctx.JSON(200, task)
}

func (tc *TaskController) DeleteTask(ctx *gin.Context) {
	taskID := ctx.Param("id")

	// Delete the task using the TaskService.
	err := tc.service.DeleteTask(taskID)

	// If the task is not found, return a 404 response.
	if err != nil {
		ctx.JSON(404, gin.H{"error": "Task not found"})
		return
	}

	// Otherwise, return a 204 response.
	ctx.JSON(204, nil)
}
