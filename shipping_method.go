package main

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type ShippingMethod struct {
	Id         int
	MethodName string
	Cost       float64
}

func AllShippingMethods(db *pgx.Conn) ([]ShippingMethod, error) {
	var shippingMethods []ShippingMethod
	rows, err := db.Query(context.Background(), "SELECT * FROM shipping_method")
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
