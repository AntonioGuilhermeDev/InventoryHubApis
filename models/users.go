package models

import (
	"errors"
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

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type PublicUser struct {
	ID                int64     `json:"id"`
	Nome              string    `json:"nome"`
	Sobrenome         string    `json:"sobrenome"`
	Email             string    `json:"email"`
	Role              string    `json:"role"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	EstabelecimentoID int64     `json:"estabelecimento_id"`
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

func (u *User) ValidateCredentials() error {
	query := "SELECT id, password, role FROM users WHERE email = $1"
	row := db.DB.QueryRow(query, u.Email)

	var retrievedPassword string
	var role string

	err := row.Scan(&u.ID, &retrievedPassword, &role)
	if err != nil {
		return errors.New("credenciais inválidas")
	}

	if !utils.CheckPasswordHash(u.Password, retrievedPassword) {
		return errors.New("credenciais inválidas")
	}

	u.Role = role
	return nil
}

func GetAllUsers() ([]PublicUser, error) {
	query := "SELECT id, nome, sobrenome, email, created_at, updated_at, role, estabelecimento_id FROM users"

	rows, err := db.DB.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []PublicUser

	for rows.Next() {
		var user PublicUser
		err := rows.Scan(&user.ID, &user.Nome, &user.Sobrenome, &user.Email, &user.CreatedAt, &user.UpdatedAt, &user.Role, &user.EstabelecimentoID)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil

}
