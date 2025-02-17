package server

import (
	"github.com/gin-gonic/gin"

	database "github.com/andreiz53/cookinator/database/handlers"
	"github.com/andreiz53/cookinator/token"
	"github.com/andreiz53/cookinator/util"
)

type Server struct {
	config     util.Config
	router     *gin.Engine
	store      database.Store
	tokenMaker token.Maker
}

func NewServer(config util.Config, store database.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.setupRoutes()
	return server, nil
}

func (server *Server) setupRoutes() {

	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	// no reason to expose this
	router.GET("/users", server.getUsers)
	// only for authenticated users
	router.GET("/users/:id", server.getUserByID)
	router.PUT("/users/email", server.updateUserEmail)
	router.PUT("/users/password", server.updateUserPassword)
	router.PUT("/users/info", server.updateUserInfo)
	router.DELETE("/users/:id", server.deleteUser)

	// no reason to expose this at the moment
	router.POST("/ingredients", server.createIngredient)
	router.GET("/ingredients", server.getIngredients)
	router.GET("/ingredients/:id", server.getIngredientByID)
	router.PUT("/ingredients", server.updateIngredient)
	router.DELETE("/ingredients/:id", server.deleteIngredient)

	// with authenticated user middleware
	authRouter := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRouter.POST("/families", server.createFamily)
	// no reason to expose this at the moment
	router.GET("/families", server.getFamilies)
	// only if the user is within that family
	router.GET("/families/:id", server.getFamilyByID)
	router.GET("/families/users/:user_id", server.getFamilyByUserID)
	router.PUT("/families", server.updateFamily)
	router.DELETE("/families/:id", server.deleteFamily)

	server.router = router
}

func (s *Server) Run(address string) error {
	return s.router.Run(address)
}
