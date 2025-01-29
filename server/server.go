package server

import (
	database "github.com/andreiz53/cookinator/database/handlers"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	store  database.Store
}

func NewServer(store database.Store) *Server {
	server := &Server{store: store}

	router := gin.Default()

	router.POST("/users", server.createUser)
	router.GET("/users", server.getUsers)
	router.GET("/users/:id", server.getUserByID)
	router.PUT("/users/email", server.updateUserEmail)
	router.PUT("/users/password", server.updateUserPassword)
	router.PUT("/users/info", server.updateUserInfo)
	router.DELETE("/users/:id", server.deleteUser)

	router.POST("/ingredients", server.createIngredient)
	router.GET("/ingredients", server.getIngredients)
	router.GET("/ingredients/:id", server.getIngredientByID)
	router.PUT("/ingredients", server.updateIngredient)
	router.DELETE("/ingredients/:id", server.deleteIngredient)

	server.router = router
	return server
}

func (s *Server) Run(address string) error {
	return s.router.Run(address)
}
