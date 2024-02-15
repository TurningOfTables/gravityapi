package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
)

type Book struct {
	Id              int       `json:"id"`
	Title           string    `json:"title"`
	Isbn            string    `json:"isbn"`
	LanguageId      int       `json:"languageId"`
	NumPages        int       `json:"numPages"`
	PublicationDate time.Time `json:"publicationDate"`
	PublisherId     int       `json:"publisherId"`
}

type Language struct {
	Id           int    `json:"id"`
	LanguageCode string `json:"languageCode"`
	LanguageName string `json:"languageName"`
}

// AllBooks returns all books from the database as []Book
// []Book is returned in all cases, so requires a check for error being nil
func AllBooks(db *pgx.Conn, c fiber.Ctx) ([]Book, error) {
	var books []Book
	rows, err := db.Query(context.Background(),
		`SELECT * FROM book LIMIT $1 OFFSET $2`, c.Locals("limit"), c.Locals("offset"))
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
func BooksBySearchTerm(db *pgx.Conn, c fiber.Ctx, searchTerm, searchValue string) ([]Book, error) {

	var books []Book
	var sql string
	switch searchTerm {
	case "title":
		sql = "SELECT * from book WHERE book.title=$1 LIMIT $2 OFFSET $3"
	case "isbn":
		sql = "SELECT * from book WHERE book.isbn13=$1 LIMIT $2 OFFSET $3"
	case "author":
		sql = `SELECT book.book_id, book.title, book.isbn13, book.language_id, book.num_pages, book.publication_date, book.publisher_id 
		FROM book
		JOIN book_author ON book_author.book_id = book.book_id
		JOIN author ON author.author_id = book_author.author_id
		AND author.author_name=$1 LIMIT $2 OFFSET $3`
	default:
		return books, errors.New("invalid search term")
	}

	log.Print(searchTerm)
	log.Print(searchValue)
	log.Print(c.Locals("limit"))
	log.Print(c.Locals("offset"))
	log.Print(sql)

	rows, err := db.Query(context.Background(), sql, searchValue, c.Locals("limit"), c.Locals("offset"))
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
