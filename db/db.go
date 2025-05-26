package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	host := os.Getenv("DB_HOST")
	portStr := os.Getenv("DB_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic(fmt.Sprintf("Porta inválida: %v", err))
	}
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	DB, err = sql.Open("postgres", psqlInfo)

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

	createProductsTable := `
	CREATE TABLE IF NOT EXISTS products (
	id SERIAL PRIMARY KEY,
	nome VARCHAR(150) NOT NULL,
	sku VARCHAR(50) NOT NULL UNIQUE,
	descricao TEXT NOT NULL,
	valor NUMERIC(10,2) NOT NULL,
	estoque NUMERIC(10,3) NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	estabelecimento_id BIGINT NOT NULL,
	FOREIGN KEY (estabelecimento_id) REFERENCES estabelecimentos(id) ON DELETE CASCADE
);
	`
	_, err = DB.Exec(createProductsTable)

	if err != nil {
		return err
	}

	log.Println("Tabela 'produtos' criada com sucesso.")

	return nil
}
