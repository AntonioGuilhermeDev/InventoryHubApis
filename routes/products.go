package routes

import (
	"net/http"

	"github.com/AntonioGuilhermeDev/InventoryHubApis/models"
	"github.com/gin-gonic/gin"
)

func createProduct(ctx *gin.Context) {
	var product models.Product

	err := ctx.ShouldBindJSON(&product)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Requisição incompleta. Todos os campos obrigatórios devem ser preenchidos.",
		})
		return
	}

	exists, err := models.SKUExists(product.SKU)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar SKU"})
		return
	}

	if exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "SKU já cadastrado"})
		return
	}

	err = product.Save()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Não foi possível cadastrar o produto. Tente novamente mais tarde.",
		})
		return
	}

	ctx.JSON(http.StatusCreated, product)
}
