package routes

import (
	"net/http"
	"strconv"
	"time"

	"github.com/AntonioGuilhermeDev/InventoryHubApis/db"
	"github.com/AntonioGuilhermeDev/InventoryHubApis/models"
	"github.com/AntonioGuilhermeDev/InventoryHubApis/utils"
	"github.com/gin-gonic/gin"
)

func createEstablishment(ctx *gin.Context) {
	var establishment models.Establishment

	err := ctx.ShouldBindJSON(&establishment)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Erro na estrutura da requisição. Verifique os parametros obrigatórios"})
		return
	}

	formatedDoc, err := utils.FormatAndValidateCpfCnpj(establishment.CPFCNPJ)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "CPF ou CNPJ inválido."})
		return
	}

	establishment.CPFCNPJ = formatedDoc

	tx, err := db.DB.Begin()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao criar o estabelecimento. Falha interna."})
		return
	}

	err = establishment.Endereco.Save(tx)

	if err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Não foi possivel criar o estabelecimento. Erro ao salvar o endereço."})
		return
	}

	establishment.EnderecoID = establishment.Endereco.ID

	exists, err := utils.CpfCnpjExists(establishment.CPFCNPJ)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar CPF/CNPJ"})
		return
	}

	if exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Esse CPF ou CNPJ já foi cadastrado por outro estabelecimento"})
		return
	}

	err = establishment.Save(tx)

	if err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Não foi possível criar o estabelecimento"})
		return
	}

	err = tx.Commit()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao criar o estabelecimento. Falha interna."})
		return
	}

	ctx.JSON(http.StatusCreated, establishment)
}

func getEstablishments(ctx *gin.Context) {
	establishment, err := models.GetAllEstablishments()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Não foi possivel listar os estabelecimentos."})
		return
	}

	ctx.JSON(http.StatusOK, establishment)
}

func getEstablishment(ctx *gin.Context) {
	establishmentId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Não foi possível converter o id"})
		return
	}

	establishment, err := models.GetEstablishmentByID(establishmentId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Não foi encontrado nenhum estabelecimento com esse ID."})
		return
	}

	ctx.JSON(http.StatusOK, establishment)

}

func updateEstablishment(ctx *gin.Context) {
	establishmentId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Não foi possível converter o id"})
		return
	}

	establishment, err := models.GetEstablishmentByID(establishmentId)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Não foi possível encontrar nenhum estabelecimento com o id"})
		return
	}

	var updatedEstablishment models.Establishment

	err = ctx.ShouldBindJSON(&updatedEstablishment)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Erro na requisição. Verifique os parametros obrigatórios e tente novamente."})
		return
	}

	updatedEstablishment.ID = establishment.ID
	updatedEstablishment.UpdatedAt = time.Now()
	updatedEstablishment.EnderecoID = establishment.EnderecoID
	updatedEstablishment.Endereco.UpdatedAt = updatedEstablishment.UpdatedAt

	exists, err := utils.CpfCnpjExistsExcludingEc(updatedEstablishment.CPFCNPJ, updatedEstablishment.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar o CPF ou CNPJ."})
		return
	}

	if exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "CPF ou CNPJ já cadastrado por outro estabelecimento."})
		return
	}

	tx, err := db.DB.Begin()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao atualizar o estabelecimento. Falha interna."})
		return
	}

	err = updatedEstablishment.Endereco.Update(tx)

	if err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Não foi possível atualizar o endereço."})
		return
	}

	err = updatedEstablishment.Update(tx)

	if err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Não foi possivel atualizar o estabelecimento."})
		return
	}

	err = tx.Commit()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao atualizar o estabelecimento. Falha interna."})
		return
	}

	ctx.JSON(http.StatusOK, updatedEstablishment)

}

func deleteEstablishment(ctx *gin.Context) {
	establishmentId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Não foi possível converter o id"})
		return
	}

	establishment, err := models.GetEstablishmentByID(establishmentId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Não foi encontrado nenhum estabelecimento com esse ID."})
		return
	}

	var deletedEstablishment models.Establishment

	deletedEstablishment.ID = establishment.ID

	tx, err := db.DB.Begin()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao deletar o estabelecimento. Falha interna."})
		return
	}

	err = deletedEstablishment.Delete(tx)

	if err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Não foi possível deletar o estabelecimento."})
		return
	}

	var deletedAddress models.Address

	deletedAddress.ID = establishment.EnderecoID

	err = deletedAddress.Delete(tx)

	if err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Não foi possível deletar o estabelecimento. Erro ao deletar o endereço"})
		return
	}

	err = tx.Commit()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao deletar o estabelecimento. Falha interna."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Estabelecimento deletado com sucesso."})

}
