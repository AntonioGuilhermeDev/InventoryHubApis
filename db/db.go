package db

import (
	"database/sql"
	"fmt"

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

}
