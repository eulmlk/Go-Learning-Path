package controllers

import (
	"log"
	"net/http"
	"task_manager/domain"
	"task_manager/usecase"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// A struct that handles task operations by calling the usecase methods.
type TaskController struct {
	usecase *usecase.TaskUsecase
}

// A constructor that creates a new instance of TaskController.
func NewTaskController(usecase *usecase.TaskUsecase) *TaskController {
	return &TaskController{usecase: usecase}
}

// A handler function that returns all tasks.
func (tc *TaskController) GetTasks(ctx *gin.Context) {
	tasks, _err := tc.usecase.GetTasks()
	if _err != nil {
		log.Println(_err.Err)
		ctx.JSON(_err.StatusCode, gin.H{"error": _err.Message})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"count": len(tasks),
		"tasks": tasks,
	})
}

// A handler function that returns a task with the given ID.
func (tc *TaskController) GetTaskByID(ctx *gin.Context) {
	// Get the task ID from the context.
	taskID := ctx.MustGet("task_id").(primitive.ObjectID)

	// Get the task using the TaskUsecase.
	task, _err := tc.usecase.GetTaskByID(taskID)
	if _err != nil {
		log.Println(_err.Err)
		ctx.JSON(_err.StatusCode, gin.H{"error": _err.Message})
		return
	}

	ctx.JSON(http.StatusOK, task)
}

// A handler function that creates a new task.
func (tc *TaskController) CreateTask(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*domain.Claims)
	taskData := &domain.CreateTaskData{}

	// Bind the request body to the struct.
	err := ctx.BindJSON(taskData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// If the status field is missing, set it to "Pending" by default.
	if taskData.Status == "" {
		taskData.Status = "Pending"
	}

	// If the status field is present, check if it is one of the allowed values.
	if taskData.Status != "Pending" && taskData.Status != "Completed" && taskData.Status != "In Progress" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "status field must be one of: Pending, Completed, In Progress"})
		return
	}

	// Create the task using the TaskUsecase.
	taskView, _err := tc.usecase.CreateTask(taskData, claims)
	if _err != nil {
		log.Println(_err.Err)
		ctx.JSON(_err.StatusCode, gin.H{"error": _err.Message})
		return
	}

	// Return a 201 response with the new task.
	ctx.JSON(http.StatusCreated, taskView)
}

// A handler function that replaces a task with the given ID.
func (tc *TaskController) UpdateTaskPut(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*domain.Claims)
	taskID := ctx.MustGet("task_id").(primitive.ObjectID)

	// Bind the request body to the struct.
	taskData := &domain.ReplaceTaskData{}
	err := ctx.BindJSON(taskData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Check if the status field is one of the allowed values.
	if taskData.Status != "Pending" && taskData.Status != "Completed" && taskData.Status != "In Progress" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "status field must be one of: Pending, Completed, In Progress"})
		return
	}

	// Replace the task using the TaskUsecase.
	task, _err := tc.usecase.ReplaceTask(taskID, taskData, claims)
	if _err != nil {
		log.Println(_err.Err)
		ctx.JSON(_err.StatusCode, gin.H{"error": _err.Message})
		return
	}

	// Otherwise, return a 200 response with the updated task.
	ctx.JSON(http.StatusOK, task)
}

// A handler function that updates a task with the given ID.
func (tc *TaskController) UpdateTaskPatch(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*domain.Claims)
	taskID := ctx.MustGet("task_id").(primitive.ObjectID)

	// Bind the request body to the struct.
	taskData := &domain.UpdateTaskData{}
	err := ctx.BindJSON(taskData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Check if the status field is one of the allowed values.
	if taskData.Status != "" && taskData.Status != "Pending" && taskData.Status != "Completed" && taskData.Status != "In Progress" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "status field must be one of: Pending, Completed, In Progress"})
		return
	}

	// Update the task using the TaskUsecase.
	task, _err := tc.usecase.UpdateTask(taskID, taskData, claims)
	if _err != nil {
		log.Println(_err.Err)
		ctx.JSON(_err.StatusCode, gin.H{"error": _err.Message})
		return
	}

	// Otherwise, return a 200 response with the updated task.
	ctx.JSON(http.StatusOK, task)
}

// A handler function that deletes a task with the given ID.
func (tc *TaskController) DeleteTask(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*domain.Claims)
	taskID := ctx.MustGet("task_id").(primitive.ObjectID)

	// Delete the task using the TaskUsecase.
	_err := tc.usecase.DeleteTask(taskID, claims)
	if _err != nil {
		log.Println(_err.Err)
		ctx.JSON(_err.StatusCode, gin.H{"error": _err.Message})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
