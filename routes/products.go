package routes

import (
	"net/http"
	"strconv"
	"time"

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

func getProducts(ctx *gin.Context) {
	products, err := models.GetAllProducts()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Não foi possivel listar os produtos."})
		return
	}

	ctx.JSON(http.StatusOK, products)
}

func getProductById(ctx *gin.Context) {
	productId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Não foi possível converter o id"})
		return
	}

	product, err := models.GetProduct(productId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Produto não encontrado."})
		return
	}

	ctx.JSON(http.StatusOK, product)
}

func updateProduct(ctx *gin.Context) {
	productId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Não foi possível converter o id"})
		return
	}

	product, err := models.GetProduct(productId)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Não foi possível encontrar nenhum produto com o id"})
		return
	}

	var updatedProduct models.Product

	err = ctx.ShouldBindJSON(&updatedProduct)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Erro na requisição. Verifique os parametros obrigatórios e tente novamente."})
		return
	}

	updatedProduct.ID = product.ID
	updatedProduct.UpdatedAt = time.Now()

	exists, err := models.SKUExistsForOtherProduct(updatedProduct.SKU, updatedProduct.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar SKU"})
		return
	}

	if exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "SKU já cadastrado"})
		return
	}

	err = updatedProduct.Update()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Não foi possivel atualizar o produto"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Produto atualizado com sucesso",
		"produto": updatedProduct,
	})
}

func deleteProduct(ctx *gin.Context) {
	productId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Não foi possível converter o id"})
		return
	}

	product, err := models.GetProduct(productId)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Não foi possível encontrar nenhum produto com o id"})
		return
	}

	var deletedProduct models.Product

	deletedProduct.ID = product.ID

	err = deletedProduct.Delete()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Não foi possivel deletar o produto"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Produto deletado com sucesso"})
}
