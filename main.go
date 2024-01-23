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
		if c.Query("name") != "" {
			return handleAuthorsByName(c, db, c.Query("name"))
		}

		return handleAllAuthors(c, db)
	})
	v1.Get("/books", func(c fiber.Ctx) error {
		if c.Query("author") != "" {
			return handleBooksByAuthor(c, db, c.Query("author"))
		}

		if c.Query("title") != "" {
			return handleBooksByTitle(c, db, c.Query("title"))
		}

		return handleAllBooks(c, db)
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
	author, err := AuthorByName(db, c.Query("name"))

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

// GET /v1/books?author=Agatha Christie
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

// GET /v1/books?title=The Tempest
func handleBooksByTitle(c fiber.Ctx, db *pgx.Conn, bookTitle string) error {
	booksByTitle, err := BooksByTitle(db, bookTitle)
	if len(booksByTitle) == 0 {
		return fiber.NewError(fiber.ErrNotFound.Code, "No books found by that title")
	}

	if err != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, fmt.Sprintf("Error retrieving books by title: %v", err.Error()))
	}

	return c.JSON(booksByTitle)
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
