package models

import (
	"database/sql"
	"errors"
	"log"
	"regexp"
	"time"

	"github.com/AntonioGuilhermeDev/InventoryHubApis/db"
	"github.com/klassmann/cpfcnpj"
)

type Establishment struct {
	ID          int64     `json:"id"`
	RazaoSocial string    `json:"razao_social" binding:"required"`
	CPFCNPJ     string    `json:"cpf_cnpj" binding:"required"`
	EnderecoID  int64     `json:"endereco_id"`
	Endereco    Address   `json:"endereco" binding:"required"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (e *Establishment) Save(tx *sql.Tx) error {
	query := `INSERT INTO estabelecimentos(razao_social, cpf_cnpj, endereco_id)
	VALUES($1, $2, $3)
	RETURNING id, created_at, updated_at
`
	err := tx.QueryRow(query, e.RazaoSocial, e.CPFCNPJ, e.EnderecoID).Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt)

	log.Println(err)
	return err

}

func GetAllEstablishments() ([]Establishment, error) {
	query := `SELECT e.id, e.razao_social, e.cpf_cnpj, e.endereco_id, e.created_at, e.updated_at,
a.logradouro, a.complemento, a.numero, a.bairro, a.cidade, a.uf, a.cep FROM estabelecimentos e
JOIN enderecos a ON a.id = e.endereco_id`

	rows, err := db.DB.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var establishments []Establishment

	for rows.Next() {
		var est Establishment
		var addr Address

		err := rows.Scan(
			&est.ID, &est.RazaoSocial, &est.CPFCNPJ, &est.EnderecoID, &est.CreatedAt, &est.UpdatedAt,
			&addr.Logradouro, &addr.Complemento, &addr.Numero, &addr.Bairro, &addr.Cidade, &addr.UF, &addr.CEP,
		)

		if err != nil {
			return nil, err
		}

		est.Endereco = addr
		establishments = append(establishments, est)
	}

	return establishments, nil
}

func GetEstablishmentByID(id int64) (*Establishment, error) {
	query := `SELECT e.id, e.razao_social, e.cpf_cnpj, e.endereco_id, e.created_at, e.updated_at,
       a.logradouro, a.complemento, a.numero, a.bairro, a.cidade, a.uf, a.cep
FROM estabelecimentos e
JOIN enderecos a ON a.id = e.endereco_id
WHERE e.id = $1`

	row := db.DB.QueryRow(query, id)

	var est Establishment
	var addr Address

	err := row.Scan(
		&est.ID, &est.RazaoSocial, &est.CPFCNPJ, &est.EnderecoID, &est.CreatedAt, &est.UpdatedAt,
		&addr.Logradouro, &addr.Complemento, &addr.Numero, &addr.Bairro, &addr.Cidade, &addr.UF, &addr.CEP,
	)

	if err != nil {
		return nil, err
	}

	est.Endereco = addr

	return &est, nil
}

func (e *Establishment) Update(tx *sql.Tx) error {
	query := `UPDATE estabelecimentos
	SET razao_social = $1, cpf_cnpj = $2, updated_at = $3, endereco_id = $4
	WHERE id = $5`

	stmt, err := tx.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(e.RazaoSocial, e.CPFCNPJ, e.UpdatedAt, e.EnderecoID, e.ID)

	if err != nil {
		return err
	}

	return nil

}

func (e *Establishment) Delete(tx *sql.Tx) error {
	query := "DELETE FROM estabelecimentos WHERE id = $1"

	stmt, err := tx.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(e.ID)

	if err != nil {
		return err
	}

	return nil
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
