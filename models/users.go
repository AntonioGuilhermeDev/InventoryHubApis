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

func EmailExists(email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	err := db.DB.QueryRow(query, email).Scan(&exists)
	return exists, err
}

func EmailExistsExcludingUser(email string, id int64) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND id != $2)`
	err := db.DB.QueryRow(query, email, id).Scan(&exists)
	return exists, err
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

func GetUserById(id int64) (*PublicUser, error) {
	query := "SELECT id, nome, sobrenome, email, created_at, updated_at, role, estabelecimento_id FROM users WHERE id = $1"

	row := db.DB.QueryRow(query, id)

	var user PublicUser

	err := row.Scan(&user.ID, &user.Nome, &user.Sobrenome, &user.Email, &user.CreatedAt, &user.UpdatedAt, &user.Role, &user.EstabelecimentoID)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *PublicUser) Update() error {
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
