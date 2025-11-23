package testhelper

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
)

var (
	testDB *sqlx.DB
)

func initDB(db *sqlx.DB) error {
	err := runMigrates(db)
	if err != nil {
		return err
	}
	testDB = db
	return nil
}

func GetTestDB() *sqlx.DB {
	return testDB
}

func CleanUp() {
	var tables = []string{"users", "todos", "projects"}
	var builder = strings.Builder{}
	for _, tableName := range tables {
		builder.WriteString("TRUNCATE TABLE ")
		builder.WriteString(tableName + " RESTART IDENTITY CASCADE;")
	}
	_, err := testDB.Exec(builder.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, "DB clean up failed:\n%v", err)
		os.Exit(1)
	}
	fmt.Fprint(os.Stdout, "✓ Tables Truncated\n")
}

func runMigrates(db *sqlx.DB) error {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	migrationPath := filepath.Join(dir, "../..", "migrations")
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

	_, err = db.Exec(builder.String())
	if err != nil {
		return err
	}
	return nil
}
