package internal

import (
	"errors"
	"os"
	"task_manager/domain"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// A function that hashes a password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// A function that compares a hashed password with a password.
func ComparePasswords(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// A function that generates a jwt token.
func GenerateToken(user *domain.User) (string, error) {
	// Get the jwt key from the environment variables.
	jwtKeyHex, ok := os.LookupEnv("JWT_KEY")
	if !ok {
		return "", errors.New("JWT_KEY not found")
	}
	jwtKey := []byte(jwtKeyHex)

	// Setup the claims.
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &domain.Claims{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
		Role:     user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Create the token.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
