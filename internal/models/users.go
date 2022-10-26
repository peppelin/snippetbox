package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

// UserModel type for db connection
type UserModel struct {
	DB *sql.DB
}

// Insert will add new users
func (m *UserModel) Insert(name, email, password string) error {
	return nil
}

// Authenticate will check if the user exists and if the provided password
// is correct. If it's ok, will return the user ID
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// Exists will check if a user existst in the user's table
func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
