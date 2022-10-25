package models

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

var ctx = context.Background()

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *pgxpool.Pool
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {

	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES ($1, $2, current_timestamp, current_timestamp + interval '1 days' * $3) returning id`

	var id int
	err := m.DB.QueryRow(ctx, stmt, title, content, expires).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {

	s := &Snippet{}

	stmt := `SELECT id, title, content, created, expires FROM snippets where expires > current_timestamp and id = $1`

	row := m.DB.QueryRow(ctx, stmt, id)

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	//shortly version
	//err := m.DB.QueryRow(ctx, stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	//if err != nil {
	//	if errors.Is(err, sql.ErrNoRows) {
	//		return nil, ErrNoRecord
	//	} else {
	//		return nil, err
	//	}
	//}
	return s, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {

	stmt := `SELECT id, title, content, created, expires FROM snippets
				WHERE expires > current_timestamp ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}

	snippets := []*Snippet{}

	defer rows.Close()

	for rows.Next() {
		// Create a pointer to a new zeroed Snippet struct.
		s := &Snippet{}
		// Use rows.Scan() to copy the values from each field in the row to the
		// new Snippet object that we created. Again, the arguments to row.Scan()
		// must be pointers to the place you want to copy the data into, and the
		// number of arguments must be exactly the same as the number of
		// columns returned by your statement.
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		// Append it to the slice of snippets.
		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
