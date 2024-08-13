package router

import (
	"task_manager/controllers"
	"task_manager/data"
	"task_manager/middleware"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// InitializeRouter initializes the Gin router and sets up the routes
func InitializeRouter(client *mongo.Client) *gin.Engine {
	// Create a new Gin router
	router := gin.Default()

	taskService := data.NewTaskService(client.Database("task_manager").Collection("tasks"))
	taskController := controllers.NewTaskController(taskService)

	userService := data.NewUserService(client.Database("task_manager").Collection("users"))
	userController := controllers.NewUserController(userService)

	// Public routes
	// A route to register a new user
	router.POST("/register", userController.RegisterUser)

	// A route to login
	router.POST("/login", userController.Login)

	// Protected routes - User
	router.Use(middleware.AuthMiddleware()) // Apply user authentication middleware
	{
		// A route to get all the tasks (user)
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

		// A route to add a new user
		router.POST("/users", userController.AddUser)

		// A route to get all the users
		router.GET("/users", userController.GetUsers)

		// A route to get a user by ID
		router.GET("/users/:id", userController.GetUserByID)

		// A route to update a user by ID
		router.PATCH("/users/:id", userController.UpdateUserPatch)

		// A route to delete a user by ID
		router.DELETE("/users/:id", userController.DeleteUser)
	}

	return router
}
