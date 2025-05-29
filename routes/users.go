package routes

import (
	"net/http"
	"strconv"
	"time"

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

func getUsers(ctx *gin.Context) {
	users, err := models.GetAllUsers()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Não foi possivel listar os usuarios"})
		return
	}

	ctx.JSON(http.StatusOK, users)

}

func getUser(ctx *gin.Context) {
	userId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Não foi possível converter o id."})
		return
	}

	user, err := models.GetUserById(userId)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Usuário não encontrado."})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func updateUser(ctx *gin.Context) {
	userId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Não foi possível converter o id."})
		return
	}

	user, err := models.GetUserById(userId)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Usuário não encontrado."})
		return
	}

	var updatedUser models.PublicUser

	err = ctx.ShouldBindJSON(&updatedUser)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Erro na requisição. Verifique os parametros obrigatórios e tente novamente."})
		return
	}

	updatedUser.ID = user.ID
	updatedUser.UpdatedAt = time.Now()

	exists, err := models.EmailExistsExcludingUser(updatedUser.Email, updatedUser.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar Email"})
		return
	}

	if exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "O email já está sendo utilizado por outro usuário"})
		return
	}

	err = updatedUser.Update()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Não foi possivel atualizar o usuário"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Usuário atualizado com sucesso",
		"usuário": updatedUser,
	})
}

func deleteUser(ctx *gin.Context) {
	userId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Não foi possível converter o id."})
		return
	}

	user, err := models.GetUserById(userId)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Usuário não encontrado."})
		return
	}

	var deletedUser models.User

	deletedUser.ID = user.ID

	err = deletedUser.Delete()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Não foi possível deletar o usuário."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Usúario deletado com sucesso"})
}
