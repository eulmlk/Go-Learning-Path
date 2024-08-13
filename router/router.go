package router

import (
	"task_manager/controllers"
	"task_manager/data"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// InitializeRouter initializes the Gin router and sets up the routes
func InitializeRouter(client *mongo.Client) *gin.Engine {
	// Create a new Gin router
	router := gin.Default()

	taskService := data.NewTaskService(client.Database("task_manager"))
	taskController := controllers.NewTaskController(taskService)

	// A route to get all the tasks
	router.GET("/tasks", taskController.GetTasks)

	// A route to get a task by ID
	router.GET("/tasks/:id", taskController.GetTaskByID)

	// A route to create a new task
	router.POST("/tasks", taskController.CreateTask)

	// A route to update a task by ID
	router.PUT("/tasks/:id", taskController.UpdateTaskPut)

	// A route to update a task by ID
	router.PATCH("/tasks/:id", taskController.UpdateTaskPatch)

	// A route to delete a task by ID
	router.DELETE("/tasks/:id", taskController.DeleteTask)

	return router
}
