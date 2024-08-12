package router

import (
	"task_manager/database"
	"task_manager/delivery/controllers"
	"task_manager/domain"
	"task_manager/infrastructure"
	"task_manager/repository"
	"task_manager/usecase"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// Sets up the public routes
func PublicRoutes(router *gin.Engine, userController *controllers.UserController) {
	router.POST("/register", userController.RegisterUser)
	router.POST("/login", userController.Login)

	router.GET("/users", userController.GetUsers)
	router.GET("/users/:id", infrastructure.IDMiddleware("user"), userController.GetUserByID)
}

// Protected Routes related to tasks
func ProtectedTaskRoutes(router *gin.Engine, taskController *controllers.TaskController) {
	router.GET("/tasks", taskController.GetTasks)
	router.POST("/tasks", taskController.CreateTask)

	router.GET("/tasks/:id", infrastructure.IDMiddleware("task"), taskController.GetTaskByID)
	router.PUT("/tasks/:id", infrastructure.IDMiddleware("task"), taskController.UpdateTaskPut)
	router.PATCH("/tasks/:id", infrastructure.IDMiddleware("task"), taskController.UpdateTaskPatch)
	router.DELETE("/tasks/:id", infrastructure.IDMiddleware("task"), taskController.DeleteTask)
}

// Protected Routes related to users
func ProtectedUserRoutes(router *gin.Engine, userController *controllers.UserController) {
	router.POST("/users", userController.AddUser)
	router.PATCH("/users/:id", infrastructure.IDMiddleware("user"), userController.UpdateUserPatch)
	router.DELETE("/users/:id", infrastructure.IDMiddleware("user"), userController.DeleteUser)
}

func GetTaskController(db *mongo.Database) *controllers.TaskController {
	taskRepository := repository.NewMongoTaskRepository(db.Collection(domain.TaskCollection))
	taskUsecase := usecase.NewTaskUsecase(taskRepository)
	taskController := controllers.NewTaskController(taskUsecase)
	return taskController
}

func GetUserController(db *mongo.Database) *controllers.UserController {
	userRepository := repository.NewMongoUserRepository(db.Collection(domain.UserCollection))
	userUsecase := usecase.NewUserUsecase(userRepository)
	userController := controllers.NewUserController(userUsecase)
	return userController
}

// InitializeRouter initializes the Gin router and sets up the routes
func InitializeRouter(client *mongo.Client) *gin.Engine {
	// Create a new Gin router
	router := gin.Default()

	// Get the task and user controllers
	db := client.Database(database.DatabaseName)
	taskController := GetTaskController(db)
	userController := GetUserController(db)

	// Public routes
	PublicRoutes(router, userController)

	// Protected routes
	router.Use(infrastructure.AuthMiddleware)
	{
		ProtectedTaskRoutes(router, taskController)
		ProtectedUserRoutes(router, userController)
	}

	return router
}
