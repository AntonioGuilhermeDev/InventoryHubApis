package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "inventory_hub"
)

var DB *sql.DB

func InitDB() {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	var err error
	DB, err = sql.Open("postgres", psqlInfo) // usa o global, sem :=

	if err != nil {
		panic(fmt.Sprintf("Erro ao abrir conexão com o banco: %v", err))
	}

	err = DB.Ping()
	if err != nil {
		panic(fmt.Sprintf("Erro ao conectar no banco: %v", err))
	}

	fmt.Println("Conexão realizada com sucesso!")

	err = createTables()
	if err != nil {
		panic(fmt.Sprintf("Erro ao criar tabelas: %v", err))
	}
}

func createTables() error {
	createEnderecosTable := `
	CREATE TABLE IF NOT EXISTS enderecos (
		id SERIAL PRIMARY KEY,
		logradouro VARCHAR(100) NOT NULL,
		complemento VARCHAR(50),
		numero INTEGER NOT NULL,
		bairro VARCHAR(50) NOT NULL,
		cidade VARCHAR(50) NOT NULL,
		uf CHAR(2) NOT NULL,
		cep VARCHAR(9) NOT NULL
	);
	`
	_, err := DB.Exec(createEnderecosTable)

	if err != nil {
		return err
	}

	log.Println("Tabela 'enderecos' criada com sucesso.")

	createEcsTable := `
	CREATE TABLE IF NOT EXISTS estabelecimentos (
		id SERIAL PRIMARY KEY,
		razao_social VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
		endereco_id INTEGER NOT NULL,
		FOREIGN KEY (endereco_id) REFERENCES enderecos(id) ON DELETE CASCADE
	);
	`
	_, err = DB.Exec(createEcsTable)

	if err != nil {
		return err
	}

	log.Println("Tabela 'estabelecimentos' criada com sucesso.")

	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email TEXT NOT NULL UNIQUE,
		nome VARCHAR(50) NOT NULL,
		sobrenome VARCHAR(50) NOT NULL,
		password TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
		role VARCHAR(20) NOT NULL CHECK (role IN ('OWNER', 'MANAGER', 'SELLER')),
		estabelecimento_id INTEGER NOT NULL,
		FOREIGN KEY (estabelecimento_id) REFERENCES estabelecimentos(id) ON DELETE CASCADE
	);
	`
	_, err = DB.Exec(createUsersTable)

	if err != nil {
		return err
	}

	log.Println("Tabela 'users' criada com sucesso.")

	return nil
}
