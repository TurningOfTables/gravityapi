package main

import (
	"context"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
)

type Country struct {
	Id          int
	CountryName string
}

// AllCountries returns all countries from the database as []Country
// []Countries is returned in all cases, so requires a check for error being nil
func AllCountries(db *pgx.Conn, c fiber.Ctx) ([]Country, error) {
	var countries []Country
	rows, err := db.Query(context.Background(), "SELECT * FROM country LIMIT $1 OFFSET $2", c.Locals("limit"), c.Locals("offset"))
	if err != nil {
		return countries, err
	}
	defer rows.Close()

	for rows.Next() {
		var c Country
		err := rows.Scan(&c.Id, &c.CountryName)
		if err != nil {
			return countries, err
		}
		countries = append(countries, c)
	}

	if err = rows.Err(); err != nil {
		return countries, err
	}
	return countries, nil
}
