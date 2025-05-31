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
	api.GET("/users", middlewares.RoleMiddleware("OWNER", "MANAGER"), getUsers)
	api.GET("/users/:id", middlewares.RoleMiddleware("OWNER", "MANAGER"), getUser)
	api.PUT("/users/:id", middlewares.RoleMiddleware("OWNER", "MANAGER"), updateUser)
	api.DELETE("/users/:id", middlewares.RoleMiddleware("OWNER", "MANAGER"), deleteUser)

	// Produtos
	api.GET("/products", getProducts)
	api.GET("/products/:id", getProductById)
	api.POST("/products", middlewares.RoleMiddleware("OWNER", "MANAGER"), createProduct)
	api.PUT("/products/:id", middlewares.RoleMiddleware("OWNER", "MANAGER"), updateProduct)
	api.DELETE("/products/:id", middlewares.RoleMiddleware("OWNER", "MANAGER"), deleteProduct)

	// Estabelecimentos
	api.POST("/establishments", middlewares.RoleMiddleware("OWNER"), createEstablishment)
	api.GET("/establishments", middlewares.RoleMiddleware("OWNER"), getEstablishments)
	api.GET("/establishments/:id", middlewares.RoleMiddleware("OWNER"), getEstablishment)
	api.PUT("/establishments/:id", middlewares.RoleMiddleware("OWNER"), updateEstablishment)
	api.DELETE("/establishments/:id", middlewares.RoleMiddleware("OWNER"), deleteEstablishment)
}
