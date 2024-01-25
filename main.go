package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

var responseSizeLimit int = 100

func main() {
	dockerMode := flag.Bool("docker", false, "Set to true to use docker env file")
	flag.Parse()

	if *dockerMode {
		log.Print("Starting in docker mode")
		godotenv.Load(".env.docker")
	} else {
		log.Print("Starting in local mode")
		godotenv.Load(".env.local")
	}

	r := initRouter()
	r.Listen(os.Getenv("GRAVITY_API_APP_HOST"))
}

// initRouter sets up a new Fiber app, connects to the database and registers API routes to handler functions
// It returns the resulting *fiber.App instance
func initRouter() *fiber.App {

	templateEngine := html.New("./views", ".html")

	r := fiber.New(fiber.Config{AppName: "Gravity API", Views: templateEngine})
	r.Use(logger.New())
	r.Use(parseLimitOffset)
	db := connectToDb()

	r.Get("/", func(c fiber.Ctx) error {
		return c.Render("index", fiber.Map{"routes": r.GetRoutes()})
	})

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

// connectToDb connects to the database specified by env var GRAVITY_API_DB_CONNECTION_STRING
// If the var isn't set, it falls back to the .env.local file
// It returns the resulting database connection
func connectToDb() *pgx.Conn {
	// Using initRouter() directly like we do when testing means
	// the env file hasn't been read, so if it's empty we read it in
	// TO DO: Avoid this check if we can refactor the order of app initialising
	if os.Getenv("GRAVITY_API_DB_CONNECTION_STRING") == "" {
		godotenv.Load(".env.local")
	}
	var dbString string = os.Getenv("GRAVITY_API_DB_CONNECTION_STRING")
	log.Printf("Connecting to db at : %v", dbString)
	conn, err := pgx.Connect(context.Background(), dbString)
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
	} else {
		log.Println("Connected to database successfully")
	}

	return conn
}

// /v1/countries

// handleAllCountries handles GET /v1/countries
func handleAllCountries(c fiber.Ctx, db *pgx.Conn) error {
	countries, err := AllCountries(db, c)
	if err != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, fmt.Sprintf("Error retrieving countries: %v", err.Error()))
	}
	return c.JSON(countries)
}

// /v1/authors

// handleAllAuthors handles GET /v1/authors
func handleAllAuthors(c fiber.Ctx, db *pgx.Conn) error {

	authors, err := AllAuthors(db, c)
	if err != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, fmt.Sprintf("Error retrieving authors: %v", err.Error()))
	}
	return c.JSON(authors)
}

// handleAuthorsSearch handles GET /v1/authors/search?<searchTerm>=<searchValue>
// Valid query params for search terms are defined within the function
func handleAuthorsSearch(c fiber.Ctx, db *pgx.Conn) error {
	var validAuthorSearchTerms = []string{"name"}

	res, err := HandleSearch(db, c, validAuthorSearchTerms, Author{}, AuthorsBySearchTerm)
	if err != nil {
		return err
	}

	return c.JSON(res)
}

// /v1/books

// handleAllBooks handles GET /v1/books
func handleAllBooks(c fiber.Ctx, db *pgx.Conn) error {
	books, err := AllBooks(db, c)
	if err != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, fmt.Sprintf("Error retrieving books: %v", err.Error()))
	}
	return c.JSON(books)
}

// handleBooksSearch handles GET /v1/books/search?<searchTerm>=<searchValue>
// Valid query params for search terms are defined within the function
func handleBooksSearch(c fiber.Ctx, db *pgx.Conn) error {
	var validBookSearchTerms = []string{"title", "isbn", "author"}

	res, err := HandleSearch(db, c, validBookSearchTerms, Book{}, BooksBySearchTerm)
	if err != nil {
		return err
	}

	return c.JSON(res)
}

// /v1/customers

// handleAllCustomers handles GET /v1/customers
func handleAllCustomers(c fiber.Ctx, db *pgx.Conn) error {
	customers, err := AllCustomers(db, c)
	if err != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, fmt.Sprintf("Error retrieving customers: %v", err.Error()))
	}
	return c.JSON(customers)
}

// handleCustomerSearch handles GET /v1/customers/search?<searchTerm>=<searchValue>
// Valid query params for search terms are defined within the function
func handleCustomersSearch(c fiber.Ctx, db *pgx.Conn) error {
	var validCustomerSearchTerms = []string{"email"}

	res, err := HandleSearch(db, c, validCustomerSearchTerms, Customer{}, CustomersBySearchTerm)
	if err != nil {
		return err
	}

	return c.JSON(res)
}

// /v1/publishers

// handleAllPublishers handles GET /v1/publishers
func handleAllPublishers(c fiber.Ctx, db *pgx.Conn) error {
	publishers, err := AllPublishers(db, c)
	if err != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, fmt.Sprintf("Error retrieving publishers: %v", err.Error()))
	}
	return c.JSON(publishers)
}

// /v1/shipping-methods

// handleAllShippingMethods handles GET /v1/shipping-methods
func handleAllShippingMethods(c fiber.Ctx, db *pgx.Conn) error {
	shippingMethods, err := AllShippingMethods(db, c)
	if err != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, fmt.Sprintf("Error retrieving shipping methods: %v", err.Error()))
	}
	return c.JSON(shippingMethods)
}

// parseLimitOffset checks for query params 'limit' and 'offset'
// Sets a c.Locals for future handlers if they are valid, otherwise sets them to defaults
// Used as LIMIT and OFFSET in subsequent SQL queries
func parseLimitOffset(c fiber.Ctx) error {
	var limit int
	var offset int

	requestedLimit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || c.Query("limit") == "" || requestedLimit > responseSizeLimit {
		limit = responseSizeLimit
	} else {
		limit = requestedLimit
	}

	requestedOffset, err := strconv.Atoi(c.Query("offset"))
	if err != nil || c.Query("offset") == "" {
		offset = 0
	} else {
		offset = requestedOffset
	}

	c.Locals("limit", fmt.Sprintf("%d", limit))
	c.Locals("offset", fmt.Sprintf("%d", offset))
	return c.Next()
}
