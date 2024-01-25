package main

import (
	"context"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
)

type ShippingMethod struct {
	Id         int
	MethodName string
	Cost       float64
}

// AllShippingMethods returns all shipping methods from the database as []ShippingMethod
// []ShoppingMethod is returned in all cases, so requires a check for error being nil
func AllShippingMethods(db *pgx.Conn, c fiber.Ctx) ([]ShippingMethod, error) {
	var shippingMethods []ShippingMethod
	rows, err := db.Query(context.Background(), "SELECT * FROM shipping_method LIMIT $1 OFFSET $2", c.Locals("limit"), c.Locals("offset"))
	if err != nil {
		return shippingMethods, err
	}
	defer rows.Close()

	for rows.Next() {
		var s ShippingMethod
		err := rows.Scan(&s.Id, &s.MethodName, &s.Cost)
		if err != nil {
			return shippingMethods, err
		}
		shippingMethods = append(shippingMethods, s)
	}

	if err = rows.Err(); err != nil {
		return shippingMethods, err
	}
	return shippingMethods, nil
}
