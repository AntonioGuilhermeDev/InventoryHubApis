package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/AntonioGuilhermeDev/InventoryHubApis/db"
)

type Product struct {
	ID                int64     `json:"id"`
	Nome              string    `json:"nome" binding:"required"`
	Descricao         string    `json:"descricao" binding:"required"`
	Valor             float64   `json:"valor" binding:"required"`
	Estoque           float64   `json:"estoque" binding:"required"`
	EstabelecimentoID int64     `json:"estabelecimento_id" binding:"required"`
	SKU               string    `json:"sku" binding:"required"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type ProductFilter struct {
	SKU         string
	Description string
	Valor       string
	StartDate   string
	EndDate     string
}

func (p *Product) Save() error {
	query := `
		INSERT INTO products (nome, sku, descricao, valor, estoque, estabelecimento_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	err := db.DB.QueryRow(
		query,
		p.Nome, p.SKU, p.Descricao, p.Valor, p.Estoque, p.EstabelecimentoID,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)

	return err
}

func GetAllProducts(role, userId string, filter ProductFilter) ([]Product, error) {
	baseQuery := "SELECT * FROM products WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if role != "OWNER" {
		var estabelecimentoId int
		err := db.DB.QueryRow("SELECT estabelecimento_id FROM users WHERE id = $1", userId).Scan(&estabelecimentoId)

		if err != nil {
			return nil, err
		}

		baseQuery += fmt.Sprintf(" AND estabelecimento_id = $%d", argIndex)
		args = append(args, estabelecimentoId)
		argIndex++
	}

	if filter.SKU != "" {
		baseQuery += fmt.Sprintf(" AND sku = $%d", argIndex)
		args = append(args, filter.SKU)
		argIndex++
	}
	if filter.Description != "" {
		baseQuery += fmt.Sprintf(" AND descricao ILIKE $%d", argIndex)
		args = append(args, "%"+filter.Description+"%")
		argIndex++
	}
	if filter.Valor != "" {
		valorStr := strings.ReplaceAll(filter.Valor, ",", ".")
		valorFloat, err := strconv.ParseFloat(valorStr, 64)
		if err == nil {
			baseQuery += fmt.Sprintf(" AND valor = $%d", argIndex)
			args = append(args, valorFloat)
			argIndex++
		}
	}
	if filter.StartDate != "" && filter.EndDate != "" {
		layoutBR := "02/01/2006"
		startDate, err1 := time.Parse(layoutBR, filter.StartDate)
		endDate, err2 := time.Parse(layoutBR, filter.EndDate)
		if err1 == nil && err2 == nil {
			endDate = endDate.Add(time.Hour*23 + time.Minute*59 + time.Second*59)

			baseQuery += fmt.Sprintf(" AND created_at BETWEEN $%d AND $%d", argIndex, argIndex+1)
			args = append(args, startDate, endDate)
			argIndex += 2
		}
	}

	rows, err := db.DB.Query(baseQuery, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []Product

	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Nome, &product.SKU, &product.Descricao, &product.Valor, &product.Estoque, &product.CreatedAt, &product.UpdatedAt, &product.EstabelecimentoID)

		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

func GetProduct(id int64, role, userId string) (*Product, error) {
	query := "SELECT * FROM products WHERE id = $1"
	row := db.DB.QueryRow(query, id)

	var product Product

	err := row.Scan(&product.ID, &product.Nome, &product.SKU, &product.Descricao, &product.Valor, &product.Estoque, &product.CreatedAt, &product.UpdatedAt, &product.EstabelecimentoID)

	if err != nil {
		return nil, err
	}

	if role != "OWNER" {
		var userEstabelecimentoId int64
		err := db.DB.QueryRow("SELECT estabelecimento_id FROM users WHERE id = $1", userId).Scan(&userEstabelecimentoId)
		if err != nil {
			return nil, fmt.Errorf("não foi possível obter estabelecimento do usuário: %w", err)
		}

		if userEstabelecimentoId != product.EstabelecimentoID {
			return nil, errors.New("acesso negado: produto não pertence ao estabelecimento do usuário")
		}
	}

	return &product, nil
}

func (p *Product) Update(role string) error {
	if role != "OWNER" {
		var currentEstabID int64
		err := db.DB.QueryRow("SELECT estabelecimento_id FROM products WHERE id = $1", p.ID).Scan(&currentEstabID)
		if err != nil {
			return err
		}
		p.EstabelecimentoID = currentEstabID
	}

	query := `UPDATE products
	SET nome = $1, sku = $2, descricao = $3, valor = $4, estoque = $5, updated_at = $6, estabelecimento_id = $7
	WHERE id = $8`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(p.Nome, p.SKU, p.Descricao, p.Valor, p.Estoque, p.UpdatedAt, p.EstabelecimentoID, p.ID)
	if err != nil {
		return err
	}

	return nil
}

func (p *Product) Delete() error {
	query := "DELETE FROM products WHERE id = $1"

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(p.ID)

	if err != nil {
		return err
	}

	return nil
}
