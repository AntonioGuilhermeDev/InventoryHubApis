package routes

import "github.com/gin-gonic/gin"

func RegisterRoutes(server *gin.Engine) {
	server.POST("/signup", signup)
	server.POST("/login", login)
	server.GET("/products", getProducts)
	server.GET("/products/:id", getProductById)
	server.POST("/products", createProduct)
	server.PUT("/products/:id", updateProduct)
	server.DELETE("/products/:id", deleteProduct)
}
