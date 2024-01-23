package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/jackc/pgx/v5"
)

var dbString string = "postgres://postgres:password@localhost:5432/gravity_books"

func main() {
	r := initRouter()

	r.Listen("127.0.0.1:3000")
}

func initRouter() *fiber.App {
	r := fiber.New(fiber.Config{AppName: "Gravity API"})
	r.Use(logger.New())
	db := connectToDb()

	r.Get("/ping", func(c fiber.Ctx) error {
		return c.SendString("pong")
	})

	v1 := r.Group("/v1")
	v1.Get("/countries", func(c fiber.Ctx) error {
		return handleAllCountries(c, db)
	})
	v1.Get("/authors", func(c fiber.Ctx) error {
		return handleAllAuthors(c, db)
	})

	v1.Get("/authors/search", func(c fiber.Ctx) error {
		var validAuthorSearchTerms = []string{"name"}
		if c.Query("name") != "" {
			return handleAuthorsByName(c, db, c.Query("name"))
		}
		return fiber.NewError(fiber.ErrBadRequest.Code, fmt.Sprintf("No valid search term provided. Valid search terms: %v", validAuthorSearchTerms))
	})

	v1.Get("/books", func(c fiber.Ctx) error {
		return handleAllBooks(c, db)
	})
	v1.Get("/books/search", func(c fiber.Ctx) error {
		var validBookSearchTerms = []string{"author", "title", "isbn"}

		// Author has a different handler function for now because of the additional columns returned by
		// querying across mapping tables like book_author
		if c.Query("author") != "" {
			return handleBooksByAuthor(c, db, c.Query("author"))
		}

		// TO DO - potential improvement here as we're checking for valid search terms twice
		// Once here in the routing, and once in the book.go handleBookSearch function
		if c.Query("title") != "" {
			return handleBookSearch(c, db, "title", c.Query("title"))
		}

		if c.Query("isbn") != "" {
			return handleBookSearch(c, db, "isbn", c.Query("isbn"))
		}

		return fiber.NewError(fiber.ErrBadRequest.Code, fmt.Sprintf("No valid search term provided. Valid search terms: %v", validBookSearchTerms))
	})

	v1.Get("/customers", func(c fiber.Ctx) error {
		return handleAllCustomers(c, db)
	})
	v1.Get("/publishers", func(c fiber.Ctx) error {
		return handleAllPublishers(c, db)
	})
	v1.Get("/shipping-methods", func(c fiber.Ctx) error {
		return handleAllShippingMethods(c, db)
	})

	return r
}

func connectToDb() *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), dbString)
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
	} else {
		log.Println("Connected to database successfully")
	}

	return conn
}

// /v1/countries

// GET /v1/countries
func handleAllCountries(c fiber.Ctx, db *pgx.Conn) error {
	countries, err := AllCountries(db)
	if err != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, fmt.Sprintf("Error retrieving countries: %v", err.Error()))
	}
	return c.JSON(countries)
}

// /v1/authors

// GET /v1/authors
func handleAllAuthors(c fiber.Ctx, db *pgx.Conn) error {
	authors, err := AllAuthors(db)
	if err != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, fmt.Sprintf("Error retrieving authors: %v", err.Error()))
	}
	return c.JSON(authors)
}

// GET /v1/authors?name=Agatha Christia
func handleAuthorsByName(c fiber.Ctx, db *pgx.Conn, authorName string) error {
	author, err := AuthorBySearchTerm(db, "name", c.Query("name"))

	if err == pgx.ErrNoRows {
		return fiber.NewError(fiber.ErrNotFound.Code, "No author found by that name")
	}

	if err != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, fmt.Sprintf("Error retrieving author by name: %v", err.Error()))
	}

	return c.JSON(author)
}

// /v1/books

// GET /v1/books
func handleAllBooks(c fiber.Ctx, db *pgx.Conn) error {
	books, err := AllBooks(db)
	if err != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, fmt.Sprintf("Error retrieving books: %v", err.Error()))
	}
	return c.JSON(books)
}

// GET /v1/books/search?author=Agatha Christie
func handleBooksByAuthor(c fiber.Ctx, db *pgx.Conn, authorName string) error {
	booksByAuthor, err := BooksByAuthor(db, authorName)
	if len(booksByAuthor) == 0 {
		return fiber.NewError(fiber.ErrNotFound.Code, "No books found by that author")
	}

	if err != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, fmt.Sprintf("Error retrieving books by author: %v", err.Error()))
	}

	return c.JSON(booksByAuthor)
}

// GET /v1/books/search?title=The Tempest
func handleBookSearch(c fiber.Ctx, db *pgx.Conn, searchTerm, searchValue string) error {
	booksBySearchTerm, err := BooksBySearchTerm(db, searchTerm, searchValue)
	if len(booksBySearchTerm) == 0 {
		return fiber.NewError(fiber.ErrNotFound.Code, "No books found by that title")
	}
	if err != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, fmt.Sprintf("Error retrieving books by search: %v", err.Error()))
	}

	return c.JSON(booksBySearchTerm)
}

// /v1/customers

// GET /v1/customers
func handleAllCustomers(c fiber.Ctx, db *pgx.Conn) error {
	customers, err := AllCustomers(db)
	if err != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, fmt.Sprintf("Error retrieving customers: %v", err.Error()))
	}
	return c.JSON(customers)
}

// /v1/publishers

// GET /v1/publishers
func handleAllPublishers(c fiber.Ctx, db *pgx.Conn) error {
	publishers, err := AllPublishers(db)
	if err != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, fmt.Sprintf("Error retrieving publishers: %v", err.Error()))
	}
	return c.JSON(publishers)
}

// /v1/shipping-methods

// GET /v1/shipping-methods
func handleAllShippingMethods(c fiber.Ctx, db *pgx.Conn) error {
	shippingMethods, err := AllShippingMethods(db)
	if err != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, fmt.Sprintf("Error retrieving shipping methods: %v", err.Error()))
	}
	return c.JSON(shippingMethods)
}
