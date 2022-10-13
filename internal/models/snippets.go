package models

import (
	"database/sql"
	"errors"
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
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires) 
	VALUES (?,?,UTC_TIMESTAMP(),DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	// Use the Exec() m,ethod to launch the sql statement.
	// If there's any errors, we return it
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// Check for the last insert ID, if its nil something went wrong
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// Get a specific Snippet
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets 
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	// Use the QueryRor() method to launch the sql statement.
	row := m.DB.QueryRow(stmt, id)
	s := &Snippet{}

	// Scan the row for values. If we didn;t retrieve any row, we return the specific error.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// We return a curtom errors from models.errors
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

// Get the last 10 snippets created.
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := []*Snippet{}
	for rows.Next() {
		s := &Snippet{}
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
