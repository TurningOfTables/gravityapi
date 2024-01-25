package main

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
)

type Author struct {
	Id         int
	AuthorName string
}

// AllAuthors returns all authors from the database as []Author
// []Author is returned in all cases, so requires a check for error being nil
func AllAuthors(db *pgx.Conn, c fiber.Ctx) ([]Author, error) {
	var authors []Author
	rows, err := db.Query(context.Background(), "SELECT * FROM author LIMIT $1 OFFSET $2", c.Locals("limit"), c.Locals("offset"))
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

// AuthorsBySearchTerm returns []Author from the database where searchTerm = searchValue
// To avoid unparameterised user input, only defined search terms are handled, otherwise in 'invalid search term' error  is returned.
// []Author is returned in all cases, so requires a check for error being nil
func AuthorsBySearchTerm(db *pgx.Conn, c fiber.Ctx, searchTerm, searchValue string) ([]Author, error) {
	var authors []Author
	var sql string
	switch searchTerm {
	case "name":
		sql = "SELECT * FROM author WHERE author.author_name=$1 LIMIT $2 OFFSET $3"
	default:
		return authors, errors.New("invalid search term")
	}

	rows, err := db.Query(context.Background(), sql, searchValue, c.Locals("limit"), c.Locals("offset"))
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
