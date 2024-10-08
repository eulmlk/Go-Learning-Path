package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// AuthMiddleware is a middleware that checks if the request is authorized
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the token from the request header
		authHeader := ctx.GetHeader("Authorization")

		// Verify that it is a Bearer Token
		authWords := strings.Fields(authHeader)
		if len(authWords) != 2 || authWords[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			ctx.Abort()
			return
		}

		// Get the token string
		tokenString := authWords[1]

		// If the token is empty, return an error
		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			ctx.Abort()
			return
		}

		// Parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure the token method conforms to "SigningMethodHMAC"
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.NewValidationError("Unexpected signing method", jwt.ValidationErrorSignatureInvalid)
			}

			// Fetch the secret key from the environment
			secretKeyHex, ok := os.LookupEnv("JWT_KEY")
			if !ok {
				return nil, jwt.NewValidationError("JWT_KEY not found", jwt.ValidationErrorSignatureInvalid)
			}

			return []byte(secretKeyHex), nil
		})

		// If there's an error while parsing the token, return an error
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		// Extract claims and validate the token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		ctx.Set("user_role", claims["role"])
		ctx.Set("user_id", claims["id"])

		ctx.Next()
	}
}
