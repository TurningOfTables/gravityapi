package main

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
)

type Customer struct {
	Id        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

// AllCustomers returns all customers from the database as []Customer
// []Customer is returned in all cases, so requires a check for error being nil
func AllCustomers(db *pgx.Conn, c fiber.Ctx) ([]Customer, error) {
	var customers []Customer
	rows, err := db.Query(context.Background(), "SELECT * FROM customer LIMIT $1 OFFSET $2", c.Locals("limit"), c.Locals("offset"))
	if err != nil {
		return customers, err
	}
	defer rows.Close()

	for rows.Next() {
		var c Customer
		err := rows.Scan(&c.Id, &c.FirstName, &c.LastName, &c.Email)
		if err != nil {
			return customers, err
		}
		customers = append(customers, c)
	}

	if err = rows.Err(); err != nil {
		return customers, err
	}
	return customers, nil
}

// CustomersBySearchTerm returns []Customer from the database where searchTerm = searchValue
// To avoid unparameterised user input, only defined search terms are handled, otherwise in 'invalid search term' error  is returned.
// []Customer is returned in all cases, so requires a check for error being nil
func CustomersBySearchTerm(db *pgx.Conn, c fiber.Ctx, searchTerm, searchValue string) ([]Customer, error) {
	var customers []Customer
	var sql string
	switch searchTerm {
	case "email":
		sql = "SELECT * from customer WHERE customer.email=$1 LIMIT $2 OFFSET $3"
	default:
		return customers, errors.New("invalid search term")
	}

	rows, err := db.Query(context.Background(), sql, searchValue, c.Locals("limit"), c.Locals("offset"))
	if err != nil {
		return customers, err
	}

	for rows.Next() {
		var c Customer
		err := rows.Scan(&c.Id, &c.FirstName, &c.LastName, &c.Email)
		if err != nil {
			return customers, err
		}
		customers = append(customers, c)
	}
	return customers, nil
}
