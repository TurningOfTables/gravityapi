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
		return handleAuthorsSearch(c, db)
	})
	v1.Get("/books", func(c fiber.Ctx) error {
		return handleAllBooks(c, db)
	})
	v1.Get("/books/search", func(c fiber.Ctx) error {
		return handleBooksSearch(c, db)
	})
	v1.Get("/customers", func(c fiber.Ctx) error {
		return handleAllCustomers(c, db)
	})
	v1.Get("/customers/search", func(c fiber.Ctx) error {
		return handleCustomersSearch(c, db)
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
func handleAuthorsSearch(c fiber.Ctx, db *pgx.Conn) error {
	var validAuthorSearchTerms = []string{"name"}

	res, err := handleSearch(db, c, validAuthorSearchTerms, Author{}, AuthorsBySearchTerm)
	if err != nil {
		return err
	}

	return c.JSON(res)
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

// GET /v1/books/search?title=The Tempest
func handleBooksSearch(c fiber.Ctx, db *pgx.Conn) error {
	var validBookSearchTerms = []string{"title", "isbn", "author"}

	res, err := handleSearch(db, c, validBookSearchTerms, Book{}, BooksBySearchTerm)
	if err != nil {
		return err
	}

	return c.JSON(res)
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

// GET /v1/customers/search?email=rvatini1@fema.gov
func handleCustomersSearch(c fiber.Ctx, db *pgx.Conn) error {
	var validCustomerSearchTerms = []string{"email"}

	res, err := handleSearch(db, c, validCustomerSearchTerms, Customer{}, CustomersBySearchTerm)
	if err != nil {
		return err
	}

	return c.JSON(res)
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
