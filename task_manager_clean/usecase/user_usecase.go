package usecase

import (
	"errors"
	"net/http"
	"task_manager/domain"
	"task_manager/infrastructure"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// A struct that defines the services for users.
type UserUsecase struct {
	userRepo domain.UserRepository
}

// A constructor that creates a new instance of UserUsecase.
func NewUserUsecase(userRepo domain.UserRepository) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}

// A method that adds a new user.
func (u *UserUsecase) AddUser(userData *domain.CreateUserData, claims *domain.Claims) (*domain.User, *domain.Error) {
	// Create a new user object.
	user := &domain.User{
		ID:       primitive.NewObjectID(),
		Username: userData.Username,
		Password: userData.Password,
		Role:     userData.Role,
	}

	// Check if the user has the correct role.
	_err := canManipulateUser(claims, user, "add")
	if _err != nil {
		return nil, _err
	}

	// Check if the username is already taken and hash the password.
	_err = u.validate(user)
	if _err != nil {
		return nil, _err
	}

	// Add the user to the database.
	err := u.userRepo.AddUser(user)
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

	// Check if the user name is already taken and hash the password.
	_err := u.validate(user)
	if _err != nil {
		return nil, _err
	}

	// Add the user to the database.
	err := u.userRepo.AddUser(user)
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
	err = infrastructure.ComparePasswords(user.Password, userData.Password)
	if err != nil {
		return "", &domain.Error{
			Err:        err,
			StatusCode: http.StatusUnauthorized,
			Message:    "Invalid username or password",
		}
	}

	// Generate a JWT token for the user.
	token, err := infrastructure.GenerateToken(user)
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
		if err == mongo.ErrNoDocuments {
			return nil, &domain.Error{
				Err:        err,
				StatusCode: http.StatusNotFound,
				Message:    "User not found",
			}
		}

		return nil, &domain.Error{
			Err:        err,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}

	return user, nil
}

// A method that updates a user by ID.
func (u *UserUsecase) UpdateUser(objectID primitive.ObjectID, userData *domain.UpdateUserData, claims *domain.Claims) (*domain.User, *domain.Error) {
	// Get the user from the database.
	user, err := u.userRepo.GetUserByID(objectID)
	if err != nil {
		return nil, &domain.Error{
			Err:        err,
			StatusCode: http.StatusNotFound,
			Message:    "User not found",
		}
	}

	// Check if the user has the correct role.
	_err := canManipulateUser(claims, user, "update")
	if _err != nil {
		return nil, _err
	}

	// Check if the username is already taken and hash the password.
	_err = u.validate(&domain.User{Username: userData.Username, Password: userData.Password})
	if _err != nil {
		return nil, _err
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
		if userData.Role != "user" && claims.Role != "root" {
			return nil, &domain.Error{
				Err:        errors.New("forbidden"),
				StatusCode: http.StatusForbidden,
				Message:    "Only root user can update role",
			}
		}

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
	// Get the user from the database.
	user, err := u.userRepo.GetUserByID(objectID)
	if err != nil {
		return &domain.Error{
			Err:        err,
			StatusCode: http.StatusNotFound,
			Message:    "User not found",
		}
	}

	// Check if the user has the correct role.
	_err := canManipulateUser(claims, user, "delete")
	if _err != nil {
		return _err
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

// A function that checks if a the logged in user can manipulate the target user.
func canManipulateUser(claims *domain.Claims, user *domain.User, manip string) *domain.Error {
	// If the user is a regular user, they can only manipulate their own account.
	if claims.Role == "user" {
		if user.ID != claims.ID {
			var message string
			if manip == "add" {
				message = "A User cannot add a new user"
			} else {
				message = "A User cannot " + manip + " another user"
			}

			return &domain.Error{
				Err:        errors.New("unauthorized"),
				StatusCode: http.StatusForbidden,
				Message:    message,
			}
		}

		return nil
	}

	// If the user is an admin, they can manipulate all users except root user and other admin users.
	if claims.Role == "admin" {
		if user.Role == "root" {
			return &domain.Error{
				Err:        errors.New("forbidden"),
				StatusCode: http.StatusForbidden,
				Message:    "Cannot " + manip + " root user",
			}
		}

		if user.Role == "admin" && claims.ID != user.ID {
			return &domain.Error{
				Err:        errors.New("unauthorized"),
				StatusCode: http.StatusForbidden,
				Message:    "Admin cannot " + manip + " another admin user",
			}
		}
	}

	// If the user is a root user, they can manipulate all users.
	return nil
}

// A helper method that checks if a username is already taken and hashes the password.
func (u *UserUsecase) validate(user *domain.User) *domain.Error {
	// Check if the username is already taken.
	if user.Username != "" {
		_, err := u.userRepo.GetUserByUsername(user.Username)
		if err == nil {
			return &domain.Error{
				Err:        errors.New("conflict"),
				StatusCode: http.StatusConflict,
				Message:    "Username already exists",
			}
		}
	}

	// Hash the user's password.
	if user.Password != "" {
		var err error
		user.Password, err = infrastructure.HashPassword(user.Password)
		if err != nil {
			return &domain.Error{
				Err:        err,
				StatusCode: http.StatusInternalServerError,
				Message:    "Internal server error",
			}
		}
	}

	return nil
}
