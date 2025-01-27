package server

import "github.com/gin-gonic/gin"

func respondWithErorr(err error) gin.H {
	return gin.H{"errror": err.Error()}
}
