package models

import (
	"database/sql"
	"time"
)

// Define the information our snippets will held.
// The fields of the struct correspond to the one of the mysql table
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// wrapping the sql.DB connection pool
type SnippetModel struct {
	DB *sql.DB
}

// SQL actions
// Add a new snippet to the database
func (m *SnippetModel) Insert(title string, content string, expires string) (int, error) {
	return 0, nil
}

// Get a specific Snippet
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	return nil, nil
}

// Get the last 10 snippets created.
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}
