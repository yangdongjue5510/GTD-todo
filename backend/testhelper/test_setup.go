package testhelper

import (
	"context"
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

	pgContainer := startPostgresContainer(ctx)

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		terminate(ctx, pgContainer)
		log.Fatalf("Failed to get connection string.\n%v", err)
	}

	sqlxDB, err := sqlx.Open("postgres", connStr)
	if err != nil {
		terminate(ctx, pgContainer)
		log.Fatalf("Failed to get connection.\n%v", err)
	}

	if err := initDB(sqlxDB); err != nil {
		terminate(ctx, pgContainer)
		log.Fatalf("Failed to init DB.\n%v", err)
	}

	os.Setenv("JWT_SECRET_KEY", "TEST_JWT_SECRET_KEY_EXAMPLE")
	code := m.Run()

	terminate(ctx, pgContainer)	
	os.Exit(code)
}

func terminate(ctx context.Context, container testcontainers.Container) {
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
