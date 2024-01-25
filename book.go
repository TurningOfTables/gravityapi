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

// AllBooks returns all books from the database as []Book
// []Book is returned in all cases, so requires a check for error being nil
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

// BooksBySearchTerm returns []Book from the database where searchTerm = searchValue
// To avoid unparameterised user input, only defined search terms are handled, otherwise in 'invalid search term' error  is returned.
// []Book is returned in all cases, so requires a check for error being nil
func BooksBySearchTerm(db *pgx.Conn, searchTerm, searchValue string) ([]Book, error) {

	var books []Book
	var sql string
	switch searchTerm {
	case "title":
		sql = "SELECT * from book WHERE book.title=$1"
	case "isbn":
		sql = "SELECT * from book WHERE book.isbn13=$1"
	case "author":
		sql = `SELECT book.book_id, book.title, book.isbn13, book.language_id, book.num_pages, book.publication_date, book.publisher_id 
		FROM book
		JOIN book_author ON book_author.book_id = book.book_id
		JOIN author ON author.author_id = book_author.author_id
		AND author.author_name=$1`
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
