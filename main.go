package main

import (
	"github.com/AntonioGuilhermeDev/InventoryHubApis/db"
	"github.com/AntonioGuilhermeDev/InventoryHubApis/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		panic("Erro ao carregar o .env")
	}

	db.InitDB()
	server := gin.Default()

	routes.RegisterRoutes(server)

	server.Run(":8080")
}
