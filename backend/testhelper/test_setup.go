package testhelper

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	db, dbContainer, err := prepareTestDatabase(ctx)
	
	if err != nil {
		log.Fatalf("Failed to prepare test database.\n%v", err)
	}

	testDB = db
	code := m.Run()
	terminateDb(db)
	terminateContainer(ctx, dbContainer)
	os.Exit(code)
}

func prepareTestDatabase(ctx context.Context) (*sqlx.DB, testcontainers.Container, error) {
	pgContainer := startPostgresContainer(ctx)

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		terminateContainer(ctx, pgContainer)
		return nil, nil, fmt.Errorf("failed to get connection string: %w", err)
	}

	sqlxDB, err := sqlx.Open("postgres", connStr)
	if err != nil {
		terminateDb(sqlxDB)
		terminateContainer(ctx, pgContainer)
		return nil, nil, fmt.Errorf("failed to open connection: %w", err)
	}

	if err := initDB(sqlxDB); err != nil {
		terminateDb(sqlxDB)
		terminateContainer(ctx, pgContainer)
		return nil, nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	os.Setenv("JWT_SECRET_KEY", "TEST_JWT_SECRET_KEY_EXAMPLE")

	return sqlxDB, pgContainer, nil
}

func terminateDb(db *sqlx.DB) {
	if db == nil {
		return
	}
	if err := db.Close(); err != nil {
		log.Printf("Failed to close test database connection.\n%v", err)
	}
}

func terminateContainer(ctx context.Context, container testcontainers.Container) {
	if container == nil {
		return
	}
	if err := container.Terminate(ctx); err != nil {
		log.Printf("Failed to terminate container.\n%v", err)
	}
}

func startPostgresContainer(ctx context.Context) *postgres.PostgresContainer {
	pgContainer, err := postgres.Run(
		ctx,
		"postgres:15",
		postgres.WithDatabase("gtd_todo"),
		postgres.WithUsername("gtd_test_user"),
		postgres.WithPassword("gtd_test_password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)

	if err != nil {
		log.Fatalf("Postgresql container initialize failed\n%v", err)
	}
	return pgContainer
}
