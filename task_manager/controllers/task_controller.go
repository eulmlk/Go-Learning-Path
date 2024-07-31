package controllers

import (
	"task_manager/data"
	"task_manager/models"
	"time"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	service *data.TaskService
}

func NewTaskController(service *data.TaskService) *TaskController {
	return &TaskController{service}
}

func (tc *TaskController) GetTasks(ctx *gin.Context) {
	tasks := tc.service.GetTasks()

	ctx.JSON(200, gin.H{
		"count": len(tasks),
		"tasks": tasks,
	})
}

func (tc *TaskController) GetTaskByID(ctx *gin.Context) {
	taskID := ctx.Param("id")
	task := tc.service.GetTaskByID(taskID)

	if task == nil {
		ctx.JSON(404, gin.H{"error": "Task not found"})
		return
	}

	ctx.JSON(200, task)
}

func (tc *TaskController) CreateTask(ctx *gin.Context) {
	newTask := models.Task{}

	err := ctx.BindJSON(&newTask)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request!"})
		return
	}

	if newTask.ID != "" {
		ctx.JSON(400, gin.H{"error": "id field is not allowed"})
		return
	}

	if newTask.Title == "" {
		ctx.JSON(400, gin.H{"error": "title field is required"})
		return
	}

	if newTask.DueDate == (time.Time{}) {
		ctx.JSON(400, gin.H{"error": "due_date field is required"})
		return
	}

	if newTask.Status == "" {
		newTask.Status = "Pending"
	}

	newTask.ID = tc.service.GenerateID()

	tc.service.CreateTask(&newTask)
	ctx.JSON(201, newTask)
}

func (tc *TaskController) UpdateTaskPut(ctx *gin.Context) {
	newTask := models.Task{}

	err := ctx.BindJSON(&newTask)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

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
	}

	taskID := ctx.Param("id")
	task := tc.service.UpdateTask(taskID, &newTask)
	if task == nil {
		ctx.JSON(404, gin.H{"error": "Task not found"})
		return
	}

	ctx.JSON(200, task)
}

func (tc *TaskController) UpdateTaskPatch(ctx *gin.Context) {
	newTask := models.Task{}

	err := ctx.BindJSON(&newTask)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	taskID := ctx.Param("id")
	task := tc.service.UpdateTask(taskID, &newTask)
	if task == nil {
		ctx.JSON(404, gin.H{"error": "Task not found"})
		return
	}

	ctx.JSON(200, task)
}

func (tc *TaskController) DeleteTask(ctx *gin.Context) {
	taskID := ctx.Param("id")
	err := tc.service.DeleteTask(taskID)
	if err != nil {
		ctx.JSON(404, gin.H{"error": "Task not found"})
		return
	}

	ctx.JSON(204, nil)
}
