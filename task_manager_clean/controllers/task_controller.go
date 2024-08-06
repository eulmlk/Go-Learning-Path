package controllers

import (
	"log"
	"net/http"
	"task_manager/data"
	"task_manager/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	tasks, _err := tc.service.GetTasks()
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

// GetTaskByID is a handler function that returns a task by ID.
func (tc *TaskController) GetTaskByID(ctx *gin.Context) {
	role := ctx.GetString("user_role")
	userID, ok := ctx.Get("user_id")
	if !ok {
		log.Println("user_id not found")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		log.Println("Invalid user_id:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	taskID := ctx.Param("id")
	task, _err := tc.service.GetTaskByID(taskID, role, objectID)
	if _err != nil {
		log.Println(_err.Err)
		ctx.JSON(_err.StatusCode, gin.H{"error": _err.Message})
		return
	}

	// Otherwise, return a 200 response with the task.
	ctx.JSON(http.StatusOK, task)
}

func (tc *TaskController) CreateTask(ctx *gin.Context) {
	newTask := models.Task{}

	// Bind the request body to the newTask struct.
	err := ctx.BindJSON(&newTask)

	// If there is an error binding the request body, return a 400 response.
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request!"})
		return
	}

	// If any of the required fields are missing, return a 400 response.
	if newTask.Title == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "title field is required"})
		return
	}

	if newTask.DueDate == (time.Time{}) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "due_date field is required"})
		return
	}

	// If the status field is missing, set it to "Pending" by default.
	if newTask.Status == "" {
		newTask.Status = "Pending"

		// If the status field is present, check if it is one of the allowed values.
	} else if newTask.Status != "Pending" && newTask.Status != "Completed" && newTask.Status != "In Progress" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "status field must be one of: Pending, Completed, In Progress"})
		return
	}

	// Generate a new ObjectID for the task.
	newTask.ID = primitive.NewObjectID()
	userID, ok := ctx.Get("user_id")
	if !ok {
		log.Println("user_id not found")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		log.Println("Invalid user_id:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	newTask.UserID = objectID

	// Create the task using the TaskService.
	_err := tc.service.CreateTask(&newTask)
	if _err != nil {
		log.Println(_err.Err)
		ctx.JSON(_err.StatusCode, gin.H{"error": _err.Message})
		return
	}

	// Return a 201 response with the new task.
	ctx.JSON(http.StatusCreated, newTask)
}

func (tc *TaskController) UpdateTaskPut(ctx *gin.Context) {
	role := ctx.GetString("user_role")

	userID, ok := ctx.Get("user_id")
	if !ok {
		log.Println("user_id not found")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		log.Println("Invalid user_id:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	newTask := models.Task{}

	// Bind the request body to the newTask struct.
	err = ctx.BindJSON(&newTask)

	// If there is an error binding the request body, return a 400 response.
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// If the ID field is present in the request body, return a 400 response.
	if newTask.ID != primitive.NilObjectID {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id field is not allowed"})
		return
	}

	// If any of the required fields are missing, return a 400 response.
	if newTask.Title == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "title field is required"})
		return
	}

	if newTask.Description == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "description field is required"})
		return
	}

	if newTask.DueDate == (time.Time{}) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "due_date field is required"})
		return
	}

	if newTask.Status == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "status field is required"})
		return
	} else if newTask.Status != "Pending" && newTask.Status != "Completed" && newTask.Status != "In Progress" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "status field must be one of: Pending, Completed, In Progress"})
		return
	}

	// Update the task using the TaskService.
	taskID := ctx.Param("id")
	task, _err := tc.service.ReplaceTask(taskID, &newTask, objectID, role)
	if _err != nil {
		log.Println(_err.Err)
		ctx.JSON(_err.StatusCode, gin.H{"error": _err.Message})
		return
	}

	if task == nil {
		panic("task not found")
	}

	// Otherwise, return a 200 response with the updated task.
	ctx.JSON(http.StatusOK, task)
}

func (tc *TaskController) UpdateTaskPatch(ctx *gin.Context) {
	role := ctx.GetString("user_role")
	userID, ok := ctx.Get("user_id")
	if !ok {
		log.Println("user_id not found")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		log.Println("Invalid user_id:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	newTask := models.Task{}

	// Bind the request body to the newTask struct.
	err = ctx.BindJSON(&newTask)

	// If there is an error binding the request body, return a 400 response.
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if newTask.Status != "" && newTask.Status != "Pending" && newTask.Status != "Completed" && newTask.Status != "In Progress" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "status field must be one of: Pending, Completed, In Progress"})
		return
	}

	// Update the task using the TaskService.
	taskID := ctx.Param("id")
	task, _err := tc.service.UpdateTask(taskID, &newTask, objectID, role)
	if _err != nil {
		log.Println(_err.Err)
		ctx.JSON(_err.StatusCode, gin.H{"error": _err.Message})
		return
	}

	// If the task is not found, return a 404 response.
	if task == nil {
		panic("task not found")
	}

	// Otherwise, return a 200 response with the updated task.
	ctx.JSON(http.StatusOK, task)
}

func (tc *TaskController) DeleteTask(ctx *gin.Context) {
	// Get the user role and ID from the context.
	role := ctx.GetString("user_role")
	userID, ok := ctx.Get("user_id")
	if !ok {
		log.Println("user_id not found")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		log.Println("Invalid user_id:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	taskID := ctx.Param("id")

	// Delete the task using the TaskService.
	_err := tc.service.DeleteTask(taskID, role, objectID)

	// If the task is not found, return a 404 response.
	if _err != nil {
		ctx.JSON(_err.StatusCode, gin.H{"error": _err.Message})
		return
	}

	// Otherwise, return a 204 response.
	ctx.JSON(http.StatusNoContent, nil)
}
