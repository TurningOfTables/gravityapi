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
	LanguageId      int `json:"-"`
	NumPages        int
	PublicationDate time.Time
	PublisherId     int `json:"-"`
	Publisher       Publisher
	Language        Language
	BookId          int `json:"-"`
	AuthorId        int `json:"-"`
	Author          Author
}

type Language struct {
	Id           int
	LanguageCode string
	LanguageName string
}

func AllBooks(db *pgx.Conn) ([]Book, error) {
	var books []Book
	rows, err := db.Query(context.Background(),
		`SELECT * FROM book, publisher, book_language, book_author, author 
	WHERE book.publisher_id = publisher.publisher_id 
	AND book.language_id = book_language.language_id 
	AND book.book_id = book_author.book_id AND book_author.author_id = author.author_id`)
	if err != nil {
		return books, err
	}
	defer rows.Close()

	for rows.Next() {
		var b Book
		err := rows.Scan(&b.Id, &b.Title, &b.Isbn, &b.LanguageId, &b.NumPages, &b.PublicationDate, &b.PublisherId, &b.Publisher.Id, &b.Publisher.PublisherName, &b.Language.Id, &b.Language.LanguageCode, &b.Language.LanguageName, &b.BookId, &b.AuthorId, &b.Author.Id, &b.Author.AuthorName)
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
