package main

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type Customer struct {
	Id        int
	FirstName string
	LastName  string
	Email     string
}

func AllCustomers(db *pgx.Conn) ([]Customer, error) {
	var customers []Customer
	rows, err := db.Query(context.Background(), "SELECT * FROM customer")
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

func CustomersBySearchTerm(db *pgx.Conn, searchTerm, searchValue string) ([]Customer, error) {
	var customers []Customer
	var sql string
	switch searchTerm {
	case "email":
		sql = "SELECT * from customer WHERE customer.email=$1"
	default:
		return customers, errors.New("invalid search term")
	}

	rows, err := db.Query(context.Background(), sql, searchValue)
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
