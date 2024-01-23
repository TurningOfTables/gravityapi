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

type Language struct {
	Id           int
	LanguageCode string
	LanguageName string
}

func AllBooks(db *pgx.Conn) ([]Book, error) {
	var books []Book
	rows, err := db.Query(context.Background(),
		`SELECT * FROM book`)
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

func BooksByAuthor(db *pgx.Conn, authorName string) ([]Book, error) {
	var books []Book
	author, err := AuthorByName(db, authorName)
	if err != nil {
		return books, err
	}

	rows, err := db.Query(context.Background(), `SELECT * from book 
	JOIN book_author ON book_author.book_id = book.book_id
	JOIN author ON author.author_id = book_author.author_id
	AND author.author_id=$1`, author.Id)
	if err != nil {
		return books, err
	}

	for rows.Next() {
		var b Book
		err := rows.Scan(&b.Id, &b.Title, &b.Isbn, &b.LanguageId, &b.NumPages, &b.PublicationDate, &b.PublisherId, nil, nil, nil, nil) // ignoring joined author info with nil. SQL could be improved I think.
		if err != nil {
			return books, err
		}
		books = append(books, b)
	}
	return books, nil
}

func BooksByTitle(db *pgx.Conn, bookTitle string) ([]Book, error) {
	var books []Book
	rows, err := db.Query(context.Background(), `SELECT * from book WHERE book.title=$1`, bookTitle)
	if err != nil {
		return books, err
	}

	for rows.Next() {
		var b Book
		err := rows.Scan(&b.Id, &b.Title, &b.Isbn, &b.LanguageId, &b.NumPages, &b.PublicationDate, &b.PublisherId)
		if err != nil {
			return books, err
		}
		books = append(books, b)
	}
	return books, nil
}
