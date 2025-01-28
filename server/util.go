package server

import "github.com/gin-gonic/gin"

func respondWithErorr(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func respondWithMessage(message string) gin.H {
	return gin.H{"message": message}
}
