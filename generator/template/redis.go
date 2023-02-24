package template

var Redis = `package postgres

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"{{.Vendor}}{{.Service}}/postgres/sqlc"
)

type DB struct {
	conn *sqlc.Queries
}

func NewDB(connString string) (*DB, func(), error) {
	// Do not use main context since some business logic closing down might still
	//  need to commit to database. Be sure to defer pool.Close in main.
	pool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to create pgx connection pool")
	}

	db := DB{
		conn: sqlc.New(pool),
	}
	return &db, pool.Close, nil
}

// QueryExample is a example of how you can use sqlc to create your database layer
func (db *DB) QueryExample() (int, error) {
	i, err := db.conn.SampleQuery(context.Background())
	if err != nil {
		return int(i), errors.Wrap(err, "Failed to query SampleQuery")
	}
	return int(i), nil
}
`
