package main

import (
	"context"
	"errors"

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

func AuthorBySearchTerm(db *pgx.Conn, searchTerm, searchValue string) ([]Author, error) {
	var authors []Author
	var sql string
	switch searchTerm {
	case "name":
		sql = "SELECT * FROM author WHERE author.author_name=$1"
	default:
		return authors, errors.New("invalid search term")
	}

	rows, err := db.Query(context.Background(), sql, searchValue)
	if err != nil {
		return authors, err
	}

	for rows.Next() {
		var a Author
		err := rows.Scan(&a.Id, &a.AuthorName)
		if err != nil {
			return authors, err
		}
		authors = append(authors, a)
	}
	return authors, nil
}
