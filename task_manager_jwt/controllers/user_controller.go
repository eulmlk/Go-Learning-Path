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

func (u *UserController) RegisterUser(ctx *gin.Context) {
	newUser := models.User{}

	err := ctx.BindJSON(&newUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	newUser.ID = primitive.NewObjectID()
	_err := u.service.RegisterUser(&newUser)
	if _err != nil {
		log.Println(_err.Err)
		ctx.JSON(_err.StatusCode, gin.H{"error": _err.Message})
		return
	}

	ctx.JSON(http.StatusCreated, newUser)
}

func (u *UserController) Login(ctx *gin.Context) {
	user := models.User{}

	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	token, _err := u.service.Login(&user)
	if _err != nil {
		log.Println(_err.Err)
		ctx.JSON(_err.StatusCode, gin.H{"error": _err.Message})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
