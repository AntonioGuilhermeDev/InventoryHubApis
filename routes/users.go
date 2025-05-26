package routes

import (
	"net/http"

	"github.com/AntonioGuilhermeDev/InventoryHubApis/models"
	"github.com/AntonioGuilhermeDev/InventoryHubApis/utils"
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

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Não foi possível cadastrar. Tente novamente mais tarde.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Cadastro realizado com sucesso",
	})
}

func login(ctx *gin.Context) {
	var input models.LoginInput

	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Email e senha são obrigatórios"})
		return
	}

	user := models.User{
		Email:    input.Email,
		Password: input.Password,
	}

	err = user.ValidateCredentials()
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Credenciais inválidas"})
		return
	}

	token, err := utils.GenerateToken(user.Email, user.Role, user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao gerar token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Login realizado com sucesso!", "token": token})
}
