package routes

import (
	"log"
	"net/http"

	"github.com/AntonioGuilhermeDev/InventoryHubApis/models"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func signup(ctx *gin.Context) {
	var user models.User

	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Requisição incompleta. Todos os campos obrigatórios devem ser preenchidos.",
		})
		return
	}

	err = user.Save()
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			ctx.JSON(http.StatusConflict, gin.H{
				"message": "Esse email já está sendo utilizado.",
			})
			return
		}

		log.Println("Erro ao salvar usuário:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Não foi possível cadastrar. Tente novamente mais tarde.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Cadastro realizado com sucesso",
	})
}
