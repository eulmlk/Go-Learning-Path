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
	router.POST("/register", userController.RegisterUser)
	router.POST("/login", userController.Login)

	// Protected routes - User
	userRoutes := router.Group("/user")
	userRoutes.Use(middleware.AuthMiddleware("user")) // Apply user authentication middleware
	{
		// A route to get all the tasks (user)
		userRoutes.GET("/tasks", taskController.GetTasks)

		// A route to get a task by ID
		userRoutes.GET("/tasks/:id", taskController.GetTaskByID)

		// A route to create a new task
		userRoutes.POST("/tasks", taskController.CreateTask)

		// A route to update a task by ID
		userRoutes.PUT("/tasks/:id", taskController.UpdateTaskPut)

		// A route to update a task by ID
		userRoutes.PATCH("/tasks/:id", taskController.UpdateTaskPatch)

		// A route to delete a task by ID
		userRoutes.DELETE("/tasks/:id", taskController.DeleteTask)
	}

	// Protected routes - Admin
	adminRoutes := router.Group("/admin")
	adminRoutes.Use(middleware.AuthMiddleware("admin")) // Apply admin authentication middleware
	{
		// A route to get all the tasks (user)
		adminRoutes.GET("/tasks", taskController.GetTasks)

		// A route to get a task by ID
		adminRoutes.GET("/tasks/:id", taskController.GetTaskByID)

		// A route to create a new task
		adminRoutes.POST("/tasks", taskController.CreateTask)

		// A route to update a task by ID
		adminRoutes.PUT("/tasks/:id", taskController.UpdateTaskPut)

		// A route to update a task by ID
		adminRoutes.PATCH("/tasks/:id", taskController.UpdateTaskPatch)

		// A route to delete a task by ID
		adminRoutes.DELETE("/tasks/:id", taskController.DeleteTask)
	}

	return router
}
