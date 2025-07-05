package models

import (
	"errors"
	"fmt"
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

func GetUserById(id int64, role, userId string) (*PublicUser, error) {
	query := "SELECT id, nome, sobrenome, email, created_at, updated_at, role, estabelecimento_id FROM users WHERE id = $1"

	row := db.DB.QueryRow(query, id)

	var user PublicUser

	err := row.Scan(&user.ID, &user.Nome, &user.Sobrenome, &user.Email, &user.CreatedAt, &user.UpdatedAt, &user.Role, &user.EstabelecimentoID)

	if err != nil {
		return nil, err
	}

	if role != "OWNER" {
		var userEstabelecimentoId int64
		err := db.DB.QueryRow("SELECT estabelecimento_id FROM users WHERE id = $1", userId).Scan(&userEstabelecimentoId)
		if err != nil {
			return nil, fmt.Errorf("não foi possível obter estabelecimento do usuário: %w", err)
		}

		if userEstabelecimentoId != user.EstabelecimentoID {
			return nil, errors.New("acesso negado: usuario não pertence ao estabelecimento do usuário")
		}
	}

	return &user, nil
}

func (u *PublicUser) Update(role string) error {
	if role != "OWNER" {
		var currentEstabID int64
		err := db.DB.QueryRow("SELECT estabelecimento_id FROM users WHERE id = $1", u.ID).Scan(&currentEstabID)
		if err != nil {
			return err
		}
		u.EstabelecimentoID = currentEstabID
	}

	query := `UPDATE users
	SET nome = $1, sobrenome = $2, email = $3, updated_at = $4, role = $5
	WHERE id = $6`

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(u.Nome, u.Sobrenome, u.Email, u.UpdatedAt, u.Role, u.ID)

	if err != nil {
		return err
	}

	return nil
}

func (u *User) Delete() error {
	query := "DELETE FROM users WHERE id = $1"

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(u.ID)

	if err != nil {
		return err
	}

	return nil
}
