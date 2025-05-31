package utils

import (
	"errors"
	"regexp"

	"github.com/AntonioGuilhermeDev/InventoryHubApis/db"
	"github.com/klassmann/cpfcnpj"
)

func EmailExists(email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	err := db.DB.QueryRow(query, email).Scan(&exists)
	return exists, err
}

func EmailExistsExcludingUser(email string, id int64) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND id != $2)`
	err := db.DB.QueryRow(query, email, id).Scan(&exists)
	return exists, err
}

func CpfCnpjExists(cpf_cnpj string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM estabelecimentos WHERE cpf_cnpj = $1)`
	err := db.DB.QueryRow(query, cpf_cnpj).Scan(&exists)
	return exists, err
}

func CpfCnpjExistsExcludingEc(cpf_cnpj string, id int64) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM estabelecimentos WHERE cpf_cnpj = $1 AND id != $2)`
	err := db.DB.QueryRow(query, cpf_cnpj, id).Scan(&exists)
	return exists, err
}

func FormatAndValidateCpfCnpj(doc string) (string, error) {
	formated := regexp.MustCompile(`\D`).ReplaceAllString(doc, "")

	switch len(formated) {
	case 11:
		if !cpfcnpj.ValidateCPF(formated) {
			return "", errors.New("cpf inválido")
		}
		return formated, nil
	case 14:
		if !cpfcnpj.ValidateCNPJ(formated) {
			return "", errors.New("cnpj inválido")
		}
		return formated, nil
	default:
		return "", errors.New("documento inválido")
	}
}

func SKUExists(sku string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM products WHERE sku = $1)`
	err := db.DB.QueryRow(query, sku).Scan(&exists)
	return exists, err
}

func SKUExistsForOtherProduct(sku string, id int64) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM products WHERE sku = $1 AND id != $2)`
	err := db.DB.QueryRow(query, sku, id).Scan(&exists)
	return exists, err
}
