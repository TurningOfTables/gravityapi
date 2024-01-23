package main

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Author struct {
	Id         int
	AuthorName string
}

func AllAuthors(db *pgx.Conn) ([]Author, error) {
	var authors []Author
	rows, err := db.Query(context.Background(), "SELECT * FROM author")
	if err != nil {
		return authors, err
	}
	defer rows.Close()

	for rows.Next() {
		var a Author
		err := rows.Scan(&a.Id, &a.AuthorName)
		if err != nil {
			return authors, err
		}
		authors = append(authors, a)
	}

	if err = rows.Err(); err != nil {
		return authors, err
	}
	return authors, nil
}

func AuthorByName(db *pgx.Conn, authorName string) (Author, error) {
	var author Author

	row := db.QueryRow(context.Background(), "SELECT * FROM author WHERE author.author_name=$1", authorName)
	if err := row.Scan(&author.Id, &author.AuthorName); err != nil {
		return author, err
	}

	return author, nil
}
