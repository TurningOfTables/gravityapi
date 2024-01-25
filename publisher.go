package main

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Publisher struct {
	Id            int
	PublisherName string
}

// AllPublishers returns all publishers from the database as []Publisher
// []Publisher is returned in all cases, so requires a check for error being nil
func AllPublishers(db *pgx.Conn) ([]Publisher, error) {
	var publishers []Publisher
	rows, err := db.Query(context.Background(), "SELECT * FROM publisher")
	if err != nil {
		return publishers, err
	}
	defer rows.Close()

	for rows.Next() {
		var p Publisher
		err := rows.Scan(&p.Id, &p.PublisherName)
		if err != nil {
			return publishers, err
		}
		publishers = append(publishers, p)
	}

	if err = rows.Err(); err != nil {
		return publishers, err
	}
	return publishers, nil
}
