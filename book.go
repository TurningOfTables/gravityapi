package main

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type Book struct {
	Id              int
	Title           string
	Isbn            string
	LanguageId      int
	NumPages        int
	PublicationDate time.Time
	PublisherId     int
}

func AllBooks(db *pgx.Conn) ([]Book, error) {
	var books []Book
	rows, err := db.Query(context.Background(), "SELECT * FROM book")
	if err != nil {
		return books, err
	}
	defer rows.Close()

	for rows.Next() {
		var b Book
		err := rows.Scan(&b.Id, &b.Title, &b.Isbn, &b.LanguageId, &b.NumPages, &b.PublicationDate, &b.PublisherId)
		if err != nil {
			return books, err
		}
		books = append(books, b)
	}

	if err = rows.Err(); err != nil {
		return books, err
	}
	return books, nil
}
