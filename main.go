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
	r := fiber.New()
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
	v1.Get("/books", func(c fiber.Ctx) error {
		return handleAllBooks(c, db)
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

func handleAllCountries(c fiber.Ctx, db *pgx.Conn) error {
	countries, err := AllCountries(db)
	if err != nil {
		return fiber.NewError(500, fmt.Sprintf("Error retrieving countries: %v", err.Error()))
	}
	return c.JSON(countries)
}

func handleAllAuthors(c fiber.Ctx, db *pgx.Conn) error {
	authors, err := AllAuthors(db)
	if err != nil {
		return fiber.NewError(500, fmt.Sprintf("Error retrieving authors: %v", err.Error()))
	}
	return c.JSON(authors)
}

func handleAllBooks(c fiber.Ctx, db *pgx.Conn) error {
	books, err := AllBooks(db)
	if err != nil {
		return fiber.NewError(500, fmt.Sprintf("Error retrieving books: %v", err.Error()))
	}
	return c.JSON(books)
}
