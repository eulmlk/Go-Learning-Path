package data

import (
	"context"
	"errors"
	"net/http"
	"os"
	"task_manager/models"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	collection *mongo.Collection
}

func NewUserService(collection *mongo.Collection) *UserService {
	return &UserService{
		collection: collection,
	}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func createToken(user models.User) (string, error) {
	jwtKeyHex, ok := os.LookupEnv("JWT_KEY")
	if !ok {
		return "", errors.New("JWT_KEY not found")
	}

	jwtKey := []byte(jwtKeyHex)
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &models.Claims{
		ID:       user.ID.Hex(),
		Username: user.Username,
		Password: user.Password,
		Role:     user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (us *UserService) AddUser(user *models.User) *models.Error {
	var err error
	err = us.collection.FindOne(context.Background(), bson.M{"username": user.Username}).Err()
	if err == nil {
		return &models.Error{
			Err:        err,
			StatusCode: http.StatusConflict,
			Message:    "Username already exists",
		}
	}

	user.Password, err = hashPassword(user.Password)
	if err != nil {
		return &models.Error{
			Err:        err,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}

	_, _err := us.collection.InsertOne(context.Background(), user)
	if _err != nil {
		return &models.Error{
			Err:        _err,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}

	return nil
}

func (us *UserService) Login(user *models.User) (string, *models.Error) {
	var result models.User

	// find user by username
	err := us.collection.FindOne(context.Background(), bson.M{"username": user.Username}).Decode(&result)
	if err != nil {
		return "", &models.Error{
			Err:        err,
			StatusCode: http.StatusUnauthorized,
			Message:    "Incorrect Username or Password",
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))
	if err != nil {
		return "", &models.Error{
			Err:        err,
			StatusCode: http.StatusUnauthorized,
			Message:    "Incorrect Username or Password",
		}
	}

	token, _err := createToken(result)
	if _err != nil {
		return "", &models.Error{
			Err:        _err,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}

	return token, nil
}

func (us *UserService) GetUsers() ([]models.User, *models.Error) {
	var users []models.User

	cursor, err := us.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, &models.Error{
			Err:        err,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}

	err = cursor.All(context.Background(), &users)
	if err != nil {
		return nil, &models.Error{
			Err:        err,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}

	return users, nil
}

func (us *UserService) GetUserByID(id string) (*models.User, *models.Error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, &models.Error{
			Err:        err,
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid ID",
		}
	}

	user := &models.User{}
	err = us.collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(user)
	if err != nil {
		return nil, &models.Error{
			Err:        err,
			StatusCode: http.StatusNotFound,
			Message:    "User not found",
		}
	}

	return user, nil
}

func (us *UserService) UpdateUser(id string, user *models.User, isRoot bool) (*models.User, *models.Error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, &models.Error{
			Err:        err,
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid ID",
		}
	}

	rootUsername, ok := os.LookupEnv("ROOT_USERNAME")
	if !ok {
		return nil, &models.Error{
			Err:        errors.New("ROOT_USERNAME is not set"),
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}

	existingUser := &models.User{}
	err = us.collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(existingUser)
	if err != nil {
		return nil, &models.Error{
			Err:        err,
			StatusCode: http.StatusNotFound,
			Message:    "User not found",
		}
	}

	if existingUser.Username == rootUsername {
		return nil, &models.Error{
			Err:        errors.New("cannot update root user"),
			StatusCode: http.StatusForbidden,
			Message:    "Cannot update root user",
		}
	}

	if !isRoot && existingUser.Role == "admin" {
		return nil, &models.Error{
			Err:        errors.New("unauthorized"),
			StatusCode: http.StatusForbidden,
			Message:    "Only root user can update admin user",
		}
	}

	if user.Username != "" {
		err = us.collection.FindOne(context.Background(), bson.M{"username": user.Username}).Err()
		if err != nil {
			return nil, &models.Error{
				Err:        err,
				StatusCode: http.StatusConflict,
				Message:    "Username already exists",
			}
		}
	}

	if user.Password != "" {
		user.Password, err = hashPassword(user.Password)
		if err != nil {
			return nil, &models.Error{
				Err:        err,
				StatusCode: http.StatusInternalServerError,
				Message:    "Internal server error",
			}
		}
	}

	update := bson.M{
		"id": objectID,
	}

	if user.Username != "" {
		update["username"] = user.Username
	}

	if user.Password != "" {
		update["password"] = user.Password
	}

	if user.Role != "" {
		update["role"] = user.Role
	}

	result, _err := us.collection.UpdateOne(context.Background(), bson.M{"_id": objectID}, bson.M{"$set": update})
	if _err != nil {
		return nil, &models.Error{
			Err:        _err,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}

	if result.MatchedCount == 0 {
		return nil, &models.Error{
			Err:        err,
			StatusCode: http.StatusNotFound,
			Message:    "User not found",
		}
	}

	updatedUser := &models.User{}
	err = us.collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(updatedUser)
	if err != nil {
		return nil, &models.Error{
			Err:        err,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}

	return updatedUser, nil
}

func (us *UserService) DeleteUser(id string, isRoot bool) *models.Error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &models.Error{
			Err:        err,
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid ID",
		}
	}

	rootUsername, ok := os.LookupEnv("ROOT_USERNAME")
	if !ok {
		return &models.Error{
			Err:        errors.New("ROOT_USERNAME is not set"),
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}

	existingUser := &models.User{}
	err = us.collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(existingUser)
	if err != nil {
		return &models.Error{
			Err:        err,
			StatusCode: http.StatusNotFound,
			Message:    "User not found",
		}
	}

	if existingUser.Username == rootUsername {
		return &models.Error{
			Err:        errors.New("cannot delete root user"),
			StatusCode: http.StatusForbidden,
			Message:    "Cannot delete root user",
		}
	}

	if !isRoot && existingUser.Role == "admin" {
		return &models.Error{
			Err:        errors.New("unauthorized"),
			StatusCode: http.StatusForbidden,
			Message:    "Only root user can delete admin user",
		}
	}

	result, _err := us.collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if _err != nil {
		return &models.Error{
			Err:        _err,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}

	if result.DeletedCount == 0 {
		return &models.Error{
			Err:        err,
			StatusCode: http.StatusNotFound,
			Message:    "User not found",
		}
	}

	return nil
}
