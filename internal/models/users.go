package models

import (
	"database/sql"
	"time"
)

type User struct {
	Name           string
	Email          string
	Created        time.Time
	ID             int
	HashedPassword []byte
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	return nil
}

func (m *UserModel) Authenticate(name, email, password string) (int, error) {
	return 0, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
