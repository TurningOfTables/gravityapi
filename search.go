package main

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
)

// Defines the possible models we allow to be searchable
type Searchable interface {
	Author | Book | Customer
}

// HandleSearch gives us a more generic way to perform searches without defining a handler for each model.
// It takes validSearchTerms, a Searchable model and a search function, then returns the resulting []s Searchable
func HandleSearch[s Searchable](db *pgx.Conn, c fiber.Ctx, validSearchTerms []string, searchModel s, searchFunc func(db *pgx.Conn, searchTerm, searchValue string) ([]s, error)) ([]s, error) {
	var results []s

	if len(c.Queries()) > 1 {
		return results, fiber.NewError(fiber.ErrBadRequest.Code, "Multiple search terms not supported")
	}

	for _, searchTerm := range validSearchTerms {
		if c.Query(searchTerm) != "" {
			results, err := searchFunc(db, searchTerm, c.Query(searchTerm))
			if len(results) == 0 && err == nil {
				return results, fiber.NewError(fiber.ErrNotFound.Code, "No results found")
			}
			if err != nil {
				return results, fiber.NewError(fiber.ErrInternalServerError.Code, fmt.Sprintf("Error retrieving by search %v", err.Error()))
			}
			return results, nil
		}
	}

	return results, fiber.NewError(fiber.ErrBadRequest.Code, fmt.Sprintf("No valid search term / value found. Valid search terms: %v", validSearchTerms))

}
