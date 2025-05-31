package models

import (
	"database/sql"
	"time"
)

type Address struct {
	ID          int64     `json:"-"`
	Logradouro  string    `json:"logradouro" binding:"required"`
	Complemento string    `json:"complemento"`
	Numero      int64     `json:"numero" binding:"required"`
	Bairro      string    `json:"bairro" binding:"required"`
	Cidade      string    `json:"cidade" binding:"required"`
	UF          string    `json:"uf" binding:"required,oneof=AC AL AP AM BA CE DF ES GO MA MT MS MG PA PB PE PI PR RJ RN RO RR RS SC SE SP TO"`
	CEP         string    `json:"cep" binding:"required"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (a *Address) Save(tx *sql.Tx) error {
	query := `INSERT INTO enderecos(logradouro, complemento, numero, bairro, cidade, uf, cep)
	VALUES($1, $2, $3, $4, $5, $6, $7)
	RETURNING id, created_at, updated_at`

	err := tx.QueryRow(query, a.Logradouro, a.Complemento, a.Numero, a.Bairro, a.Cidade, a.UF, a.CEP).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)

	return err
}

func (a *Address) Update(tx *sql.Tx) error {
	query := `UPDATE enderecos
	SET logradouro = $1, complemento = $2, numero = $3, bairro = $4, cidade = $5, uf = $6, cep = $7, updated_at = $8
	WHERE id = $9`

	stmt, err := tx.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(a.Logradouro, a.Complemento, a.Numero, a.Bairro, a.Cidade, a.UF, a.CEP, a.UpdatedAt, a.ID)

	if err != nil {
		return err
	}

	return nil
}

func (a *Address) Delete(tx *sql.Tx) error {
	query := "DELETE FROM enderecos WHERE id = $1"

	stmt, err := tx.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(a.ID)

	if err != nil {
		return err
	}

	return nil
}
