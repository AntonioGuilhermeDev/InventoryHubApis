package models

import (
	"time"

	"github.com/AntonioGuilhermeDev/InventoryHubApis/db"
	"github.com/AntonioGuilhermeDev/InventoryHubApis/utils"
)

type User struct {
	ID                int64     `json:"id"`
	Nome              string    `json:"nome" binding:"required"`
	Sobrenome         string    `json:"sobrenome" binding:"required"`
	Email             string    `json:"email" binding:"required"`
	Password          string    `json:"password" binding:"required"`
	Role              string    `json:"role" binding:"required"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	EstabelecimentoID int64     `json:"estabelecimento_id" binding:"required"`
}

func (u *User) Save() error {
	query := `INSERT INTO users(nome, sobrenome, email, password, role, estabelecimento_id)
	          VALUES ($1, $2, $3, $4, $5, $6)
	          RETURNING id;`

	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	err = db.DB.QueryRow(query,
		u.Nome,
		u.Sobrenome,
		u.Email,
		hashedPassword,
		u.Role,
		u.EstabelecimentoID,
	).Scan(&u.ID)

	if err != nil {
		return err
	}

	return nil
}
