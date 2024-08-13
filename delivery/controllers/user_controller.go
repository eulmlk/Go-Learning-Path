package controllers

import (
	"log"
	"net/http"
	"task_manager/domain"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// A struct that handles user operations by calling the usecase methods.
type UserController struct {
	usecase domain.UserUsecase
}

// A constructor that creates a new instance of UserController.
func NewUserController(usecase domain.UserUsecase) *UserController {
	return &UserController{usecase: usecase}
}

// A handler function that registers a new user.
func (uc *UserController) RegisterUser(ctx *gin.Context) {
	newUser := &domain.AuthUserData{}

	// Bind the request body to the newUser struct.
	err := ctx.BindJSON(&newUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Add the user to the database using the user usecase.
	addedUser, _err := uc.usecase.RegisterUser(newUser)
	if _err != nil {
		log.Println(_err.Err)
		ctx.JSON(_err.StatusCode, gin.H{"error": _err.Message})
		return
	}

	ctx.JSON(http.StatusCreated, addedUser)
}

// A handler function that logs a user in.
func (uc *UserController) Login(ctx *gin.Context) {
	user := &domain.AuthUserData{}

	// Bind the request body to the user struct.
	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	// Log the user in using the user usecase.
	token, _err := uc.usecase.LoginUser(user)
	if _err != nil {
		log.Println(_err.Err)
		ctx.JSON(_err.StatusCode, gin.H{"error": _err.Message})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

// A handler function that adds a new user.
func (uc *UserController) AddUser(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*domain.Claims)

	// Bind the request body to the user struct.
	user := &domain.CreateUserData{}
	err := ctx.BindJSON(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	// Add the user to the database using the user usecase.
	addedUser, _err := uc.usecase.AddUser(user, claims)
	if _err != nil {
		log.Println(_err.Err)
		ctx.JSON(_err.StatusCode, gin.H{"error": _err.Message})
		return
	}

	ctx.JSON(http.StatusCreated, addedUser)
}

// A handler function that returns all users.
func (uc *UserController) GetUsers(ctx *gin.Context) {
	// Get all users using the user usecase.
	log.Println("GetUsers")
	users, _err := uc.usecase.GetUsers()
	if _err != nil {
		log.Println(_err.Err)
		ctx.JSON(_err.StatusCode, gin.H{"error": _err.Message})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"count": len(users),
		"users": users,
	})
}

// A handler function that returns a user with the given ID.
func (uc *UserController) GetUserByID(ctx *gin.Context) {
	// Get the user ID from the context.
	userID := ctx.MustGet("user_id").(primitive.ObjectID)

	// Get the user using the user usecase.
	user, _err := uc.usecase.GetUserByID(userID)
	if _err != nil {
		log.Println(_err.Err)
		ctx.JSON(_err.StatusCode, gin.H{"error": _err.Message})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// A handler function that updates a user with the given ID.
func (uc *UserController) UpdateUserPatch(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*domain.Claims)
	userID := ctx.MustGet("user_id").(primitive.ObjectID)

	// Bind the request body to the user struct.
	userData := &domain.UpdateUserData{}
	err := ctx.BindJSON(userData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	// Update the user using the user usecase.
	user, _err := uc.usecase.UpdateUser(userID, userData, claims)
	if _err != nil {
		log.Println(_err.Err)
		ctx.JSON(_err.StatusCode, gin.H{"error": _err.Message})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// A handler function that deletes a user with the given ID.
func (uc *UserController) DeleteUser(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*domain.Claims)
	userID := ctx.MustGet("user_id").(primitive.ObjectID)

	// Delete the user using the user usecase.
	_err := uc.usecase.DeleteUser(userID, claims)
	if _err != nil {
		log.Println(_err.Err)
		ctx.JSON(_err.StatusCode, gin.H{"error": _err.Message})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
