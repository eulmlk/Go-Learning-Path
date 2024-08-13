package controllers

import (
	"log"
	"net/http"
	"task_manager/data"
	"task_manager/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	service *data.UserService
}

func NewUserController(service *data.UserService) *UserController {
	return &UserController{service: service}
}

func (uc *UserController) RegisterUser(ctx *gin.Context) {
	newUser := models.User{}

	err := ctx.BindJSON(&newUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	newUser.ID = primitive.NewObjectID()
	newUser.Role = "user"

	if newUser.Username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}

	if newUser.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Password is required"})
		return
	}

	_err := uc.service.AddUser(&newUser)
	if _err != nil {
		log.Println(_err.Err)
		ctx.JSON(_err.StatusCode, gin.H{"error": _err.Message})
		return
	}

	ctx.JSON(http.StatusCreated, newUser)
}

func (uc *UserController) Login(ctx *gin.Context) {
	user := models.User{}

	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	if user.Username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}

	if user.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Password is required"})
		return
	}

	if user.Role != "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Role cannot be set"})
		return
	}

	token, _err := uc.service.Login(&user)
	if _err != nil {
		log.Println(_err.Err)
		ctx.JSON(_err.StatusCode, gin.H{"error": _err.Message})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (uc *UserController) AddUser(ctx *gin.Context) {
	role := ctx.MustGet("user_role").(string)
	if role != "admin" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user := models.User{}
	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	user.ID = primitive.NewObjectID()
	if user.Username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}

	if user.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Password is required"})
		return
	}

	if user.Role == "" {
		if role == "root" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Role is required"})
			return
		} else {
			user.Role = "user"
		}
	}

	if role != "root" && user.Role == "admin" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Only the root user can create new admins"})
		return
	}

	if user.Role != "user" && user.Role != "admin" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Role can only be 'user' or 'admin'"})
		return
	}

	_err := uc.service.AddUser(&user)
	if _err != nil {
		log.Println(_err.Err)
		ctx.JSON(_err.StatusCode, gin.H{"error": _err.Message})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func (uc *UserController) GetUsers(ctx *gin.Context) {
	role := ctx.MustGet("user_role").(string)
	if role != "admin" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	users, _err := uc.service.GetUsers()
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

func (uc *UserController) GetUserByID(ctx *gin.Context) {
	role := ctx.MustGet("user_role").(string)
	if role != "admin" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := ctx.Param("id")

	user, _err := uc.service.GetUserByID(id)
	if _err != nil {
		log.Println(_err.Err)
		ctx.JSON(_err.StatusCode, gin.H{"error": _err.Message})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (uc *UserController) UpdateUserPatch(ctx *gin.Context) {
	role := ctx.MustGet("user_role").(string)
	if role != "admin" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := ctx.Param("id")
	userData := models.User{}

	err := ctx.BindJSON(&userData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	user, _err := uc.service.UpdateUser(id, &userData, role == "root")
	if _err != nil {
		log.Println(_err.Err)
		ctx.JSON(_err.StatusCode, gin.H{"error": _err.Message})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (uc *UserController) DeleteUser(ctx *gin.Context) {
	role := ctx.MustGet("user_role").(string)
	if role != "admin" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := ctx.Param("id")
	_err := uc.service.DeleteUser(id, role == "root")
	if _err != nil {
		log.Println(_err.Err)
		ctx.JSON(_err.StatusCode, gin.H{"error": _err.Message})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
