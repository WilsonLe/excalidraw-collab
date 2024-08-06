package models

import (
	"database/sql"
	"errors"

	"github.com/wilsonle/excalidraw-collab/pkg/database"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int
	Username string
	Password string
	Role     sql.NullString
}

func CreateUser(username, password, role string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = database.DB.Exec("INSERT INTO users (username, password, role) VALUES (?, ?, ?)", username, hashedPassword, role)
	return err
}

func GetUserByUsername(username string) (*User, error) {
	user := User{}
	err := database.DB.QueryRow("SELECT id, username, password, role FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func AuthenticateUser(username, password string) (bool, error) {
	user, err := GetUserByUsername(username)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return false, err
	}

	return true, nil
}
