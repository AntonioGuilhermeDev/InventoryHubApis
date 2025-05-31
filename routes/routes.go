package routes

import "github.com/gin-gonic/gin"

func RegisterRoutes(server *gin.Engine) {
	server.POST("/signup", signup)
	server.POST("/login", login)

	// Usuarios
	server.GET("/users", getUsers)
	server.GET("/users/:id", getUser)
	server.PUT("/users/:id", updateUser)
	server.DELETE("/users/:id", deleteUser)

	//Produtos
	server.GET("/products", getProducts)
	server.GET("/products/:id", getProductById)
	server.POST("/products", createProduct)
	server.PUT("/products/:id", updateProduct)
	server.DELETE("/products/:id", deleteProduct)

	//Estabelecimentos
	server.POST("/establishments", createEstablishment)
	server.GET("/establishments", getEstablishments)
	server.GET("/establishments/:id", getEstablishment)
	server.PUT("/establishments/:id", updateEstablishment)
	server.DELETE("/establishments/:id", deleteEstablishment)
}
