package main

import (
	"context"
	"errors"
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
	author, err := AuthorBySearchTerm(db, "name", authorName)
	if err != nil {
		return books, err
	}

	rows, err := db.Query(context.Background(), `SELECT * from book 
	JOIN book_author ON book_author.book_id = book.book_id
	JOIN author ON author.author_id = book_author.author_id
	AND author.author_id=$1`, author[0].Id) // potential issue here as we assume first author result is the one we want
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

func BooksBySearchTerm(db *pgx.Conn, searchTerm, searchValue string) ([]Book, error) {
	var books []Book
	var sql string
	switch searchTerm {
	case "title":
		sql = "SELECT * from book WHERE book.title=$1"
	case "isbn":
		sql = "SELECT * from book WHERE book.isbn13=$1"
	default:
		return books, errors.New("invalid search term")
	}

	rows, err := db.Query(context.Background(), sql, searchValue)
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
