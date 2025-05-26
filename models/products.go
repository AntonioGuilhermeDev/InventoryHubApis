package models

import (
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

func SKUExists(sku string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM products WHERE sku = $1)`
	err := db.DB.QueryRow(query, sku).Scan(&exists)
	return exists, err
}
