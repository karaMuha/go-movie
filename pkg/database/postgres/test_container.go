package postgres

import (
	"context"
	"database/sql"
	"path/filepath"
	"time"

	_ "github.com/lib/pq"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func CreateTestContainer(ctx context.Context, dbName string, initScriptPath string) (*sql.DB, error) {
	pgContainer, err := postgres.Run(ctx,
		"postgres:16.3-alpine",
		postgres.WithInitScripts(filepath.Join("..", "..", "..", "dbscripts", "public_schema.sql")),
		postgres.WithDatabase(dbName),
		postgres.WithUsername("admin"),
		postgres.WithPassword("secret"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(10*time.Second),
		),
	)

	if err != nil {
		return nil, err
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")

	if err != nil {
		return nil, err
	}

	dbHandler, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	err = dbHandler.Ping()

	if err != nil {
		return nil, err
	}

	return dbHandler, nil
}
