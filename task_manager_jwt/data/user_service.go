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
		Role:     "user",
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

func (us *UserService) RegisterUser(user *models.User) *models.Error {
	var err error
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
