package infrastructure

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// A middleware that checks if the task ID is valid
func IDMiddleware(idType string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the task ID from the request
		taskID := ctx.Param("id")

		// Convert the task ID to an ObjectID
		objectID, err := primitive.ObjectIDFromHex(taskID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid " + idType + " ID"})
			ctx.Abort()
			return
		}

		// Set the task ID in the context
		ctx.Set(idType+"_id", objectID)
		ctx.Next()
	}
}
