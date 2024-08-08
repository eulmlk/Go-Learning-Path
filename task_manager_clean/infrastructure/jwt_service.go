package infrastructure

import (
	"errors"
	"os"
	"task_manager/domain"
	"time"

	"github.com/golang-jwt/jwt"
)

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
