package backend

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	testDB      *sql.DB
	pgContainer *postgres.PostgresContainer
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	pgContainer, err = postgres.Run(
		ctx,
		"postgres:15",
		postgres.WithDatabase("gtd_todo"),
		postgres.WithUsername("gtduser"),
		postgres.WithPassword("password1234"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Postgresql container initialize failed\n%v", err)
		os.Exit(1)
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get connection string.\n%v", err)
		pgContainer.Terminate(ctx)
		os.Exit(1)
	}

	testDB, err = sql.Open("postgres", connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get connection.\n%v", err)
		pgContainer.Terminate(ctx)
		os.Exit(1)
	}

	err = runMigrates()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to run migrates.\n%v", err)
		pgContainer.Terminate(ctx)
		os.Exit(1)
	}

	code := m.Run()

	pgContainer.Terminate(ctx)
	os.Exit(code)
}

func runMigrates() error {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	migrationPath := filepath.Join(dir, "..", "migrations")
	migrationEntries, err := os.ReadDir(migrationPath)
	if err != nil {
		return err
	}

	var sqlFileNames []string
	for _, entry := range migrationEntries {
		if strings.HasSuffix(entry.Name(), ".up.sql") {
			sqlFileNames = append(sqlFileNames, entry.Name())
		}
	}

	sort.Strings(sqlFileNames)

	var builder strings.Builder
	for _, fileName := range sqlFileNames {
		sqlFilePath := filepath.Join(migrationPath, fileName)
		bytes, err := os.ReadFile(sqlFilePath)
		if err != nil {
			return err
		}
		builder.Write(bytes)
		fmt.Printf("✓ Loaded: %s\n", fileName)
	}

	_, err = testDB.Exec(builder.String())
	if err != nil {
		return err
	}

	return nil
}
