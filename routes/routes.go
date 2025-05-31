package routes

import (
	"github.com/AntonioGuilhermeDev/InventoryHubApis/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.POST("/signup", signup)
	server.POST("/login", login)

	api := server.Group("/")
	api.Use(middlewares.AuthMiddleware())

	// Usuários
	api.GET("/users", getUsers)
	api.GET("/users/:id", getUser)
	api.PUT("/users/:id", updateUser)
	api.DELETE("/users/:id", deleteUser)

	// Produtos
	api.GET("/products", getProducts)
	api.GET("/products/:id", getProductById)
	api.POST("/products", createProduct)
	api.PUT("/products/:id", updateProduct)
	api.DELETE("/products/:id", deleteProduct)

	// Estabelecimentos
	api.POST("/establishments", createEstablishment)
	api.GET("/establishments", getEstablishments)
	api.GET("/establishments/:id", getEstablishment)
	api.PUT("/establishments/:id", updateEstablishment)
	api.DELETE("/establishments/:id", deleteEstablishment)
}
