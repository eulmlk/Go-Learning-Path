package usecase

import (
	"errors"
	"net/http"
	"task_manager/domain"
	"task_manager/internal"
	"task_manager/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// A struct that defines the services for users.
type UserUsecase struct {
	userRepo repository.UserRepository
}

// A constructor that creates a new instance of UserUsecase.
func NewUserUsecase(userRepo repository.UserRepository) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}

// A method that adds a new user.
func (u *UserUsecase) AddUser(userData *domain.CreateUserData, claims *domain.Claims) (*domain.User, *domain.Error) {
	// Check if the user has the correct role.
	if claims.Role != "admin" {
		return nil, &domain.Error{
			Err:        errors.New("unauthorized"),
			StatusCode: http.StatusUnauthorized,
			Message:    "Only admin can add a new user",
		}
	}

	// Check if the user is trying to add a root user.
	if userData.Role == "root" {
		return nil, &domain.Error{
			Err:        errors.New("forbidden"),
			StatusCode: http.StatusForbidden,
			Message:    "Cannot add root user",
		}
	}

	// Check if the user is trying to add an admin user.
	if claims.Role != "root" && userData.Role == "admin" {
		return nil, &domain.Error{
			Err:        errors.New("unauthorized"),
			StatusCode: http.StatusForbidden,
			Message:    "Only root user can add admin user",
		}
	}

	// Create a new user.
	user := &domain.User{
		ID:       primitive.NewObjectID(),
		Username: userData.Username,
		Password: userData.Password,
		Role:     userData.Role,
	}

	// Hash the user's password.
	var err error
	user.Password, err = hashPassword(user.Password)
	if err != nil {
		return nil, &domain.Error{
			Err:        err,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}

	// Add the user to the database.
	err = u.userRepo.AddUser(user)
	if err != nil {
		return nil, &domain.Error{
			Err:        err,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}

	return user, nil
}

// A method that registers a new user.
func (u *UserUsecase) RegisterUser(userData *domain.AuthUserData) (*domain.User, *domain.Error) {
	// Create a new user.
	user := &domain.User{
		ID:       primitive.NewObjectID(),
		Username: userData.Username,
		Password: userData.Password,
		Role:     "user",
	}

	// Check if the username already exists.
	_, err := u.userRepo.GetUserByUsername(user.Username)
	if err == nil {
		return nil, &domain.Error{
			Err:        errors.New("conflict"),
			StatusCode: http.StatusConflict,
			Message:    "Username already exists",
		}
	}

	// Hash the user's password.
	user.Password, err = hashPassword(user.Password)
	if err != nil {
		return nil, &domain.Error{
			Err:        err,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}

	// Add the user to the database.
	err = u.userRepo.AddUser(user)
	if err != nil {
		return nil, &domain.Error{
			Err:        err,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}

	return user, nil
}

// A method that logs in a user.
func (u *UserUsecase) LoginUser(userData *domain.AuthUserData) (string, *domain.Error) {
	// Get the user from the database.
	user, err := u.userRepo.GetUserByUsername(userData.Username)
	if err != nil {
		return "", &domain.Error{
			Err:        err,
			StatusCode: http.StatusNotFound,
			Message:    "Invalid username or password",
		}
	}

	// Compare the user's password with the given password.
	err = internal.ComparePasswords(user.Password, userData.Password)
	if err != nil {
		return "", &domain.Error{
			Err:        err,
			StatusCode: http.StatusUnauthorized,
			Message:    "Invalid username or password",
		}
	}

	// Generate a JWT token for the user.
	token, err := internal.GenerateToken(user)
	if err != nil {
		return "", &domain.Error{
			Err:        err,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}

	return token, nil
}

// A method that gets all users.
func (u *UserUsecase) GetUsers() ([]domain.User, *domain.Error) {
	// Get all users from the database.
	users, err := u.userRepo.GetUsers()
	if err != nil {
		return nil, &domain.Error{
			Err:        err,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}

	return users, nil
}

// A method that gets a user by ID.
func (u *UserUsecase) GetUserByID(objectID primitive.ObjectID) (*domain.User, *domain.Error) {
	// Get the user from the database.
	user, err := u.userRepo.GetUserByID(objectID)
	if err != nil {
		return nil, &domain.Error{
			Err:        err,
			StatusCode: http.StatusNotFound,
			Message:    "User not found",
		}
	}

	return user, nil
}

// A method that updates a user by ID.
func (u *UserUsecase) UpdateUser(objectID primitive.ObjectID, userData *domain.UpdateUserData, claims *domain.Claims) (*domain.User, *domain.Error) {
	// Check if the user has the correct role.
	if claims.Role != "admin" {
		return nil, &domain.Error{
			Err:        errors.New("unauthorized"),
			StatusCode: http.StatusUnauthorized,
			Message:    "Only admin can update a user",
		}
	}

	// Get the user from the database.
	user, err := u.userRepo.GetUserByID(objectID)
	if err != nil {
		return nil, &domain.Error{
			Err:        err,
			StatusCode: http.StatusNotFound,
			Message:    "User not found",
		}
	}

	// Check if the user is trying to update a root user.
	if user.Role == "root" {
		return nil, &domain.Error{
			Err:        errors.New("forbidden"),
			StatusCode: http.StatusForbidden,
			Message:    "Cannot update root user",
		}
	}

	// Check if the user is trying to update an admin user.
	if claims.Role != "root" && user.Role == "admin" {
		return nil, &domain.Error{
			Err:        errors.New("unauthorized"),
			StatusCode: http.StatusForbidden,
			Message:    "Only root user can update admin user",
		}
	}

	// Check if the username already exists.
	if userData.Username != "" {
		_, err := u.userRepo.GetUserByUsername(userData.Username)
		if err == nil {
			return nil, &domain.Error{
				Err:        errors.New("conflict"),
				StatusCode: http.StatusConflict,
				Message:    "Username already exists",
			}
		}
	}

	// Hash the user's password.
	if userData.Password != "" {
		userData.Password, err = hashPassword(userData.Password)
		if err != nil {
			return nil, &domain.Error{
				Err:        err,
				StatusCode: http.StatusInternalServerError,
				Message:    "Internal server error",
			}
		}
	}

	// Create a map to store the updated data.
	updateData := bson.M{}
	if userData.Username != "" {
		updateData["username"] = userData.Username
	}
	if userData.Password != "" {
		updateData["password"] = userData.Password
	}
	if userData.Role != "" {
		updateData["role"] = userData.Role
	}

	// Update the user in the database.
	user, err = u.userRepo.UpdateUser(objectID, bson.M{"$set": updateData})
	if err != nil {
		return nil, &domain.Error{
			Err:        err,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}

	return user, nil
}

// A method that deletes a user by ID.
func (u *UserUsecase) DeleteUser(objectID primitive.ObjectID, claims *domain.Claims) *domain.Error {
	// Check if the user has the correct role.
	if claims.Role != "admin" {
		return &domain.Error{
			Err:        errors.New("unauthorized"),
			StatusCode: http.StatusUnauthorized,
			Message:    "Only admin can delete a user",
		}
	}

	// Get the user from the database.
	user, err := u.userRepo.GetUserByID(objectID)
	if err != nil {
		return &domain.Error{
			Err:        err,
			StatusCode: http.StatusNotFound,
			Message:    "User not found",
		}
	}

	// Check if the user is trying to delete a root user.
	if user.Role == "root" {
		return &domain.Error{
			Err:        errors.New("forbidden"),
			StatusCode: http.StatusForbidden,
			Message:    "Cannot delete root user",
		}
	}

	// Check if the user is trying to delete an admin user.
	if claims.Role != "root" && user.Role == "admin" {
		return &domain.Error{
			Err:        errors.New("unauthorized"),
			StatusCode: http.StatusForbidden,
			Message:    "Only root user can delete admin user",
		}
	}

	// Delete the user from the database.
	err = u.userRepo.DeleteUser(objectID)
	if err != nil {
		return &domain.Error{
			Err:        err,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}

	return nil
}
