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

func GetAllProducts() ([]Product, error) {
	query := "SELECT * FROM products"

	rows, err := db.DB.Query(query)

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

func GetProduct(id int64) (*Product, error) {
	query := "SELECT * FROM products WHERE id = $1"
	row := db.DB.QueryRow(query, id)

	var product Product

	err := row.Scan(&product.ID, &product.Nome, &product.SKU, &product.Descricao, &product.Valor, &product.Estoque, &product.CreatedAt, &product.UpdatedAt, &product.EstabelecimentoID)

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (p *Product) Update() error {
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
