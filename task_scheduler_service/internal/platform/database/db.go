package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func ConnectToDB(dsn string, ctx context.Context) (*bun.DB, error) {

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	fmt.Println("Connected to database!")
	return db, nil

}

func RegisterModels(db *bun.DB, ctx context.Context, models ...any) error {
	for _, model := range models {
		_, err := db.NewCreateTable().
			Model(model).
			IfNotExists().
			Exec(ctx)

		if err != nil {
			return fmt.Errorf("error creating table for %T: %w", model, err)
		}

		log.Printf("Table ready for model: %T", model)
	}
	return nil
}
