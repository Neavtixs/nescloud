package repository

import (
	"database/sql"
	"path/filepath"
	"runtime"
	"testing"

	"nescloud/backend-app/configs"
)

func SetupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	_, filename, _, _ := runtime.Caller(0)
	envPath := filepath.Join(filepath.Dir(filename), "..", "..", "..", "..", ".env")

	configs.LoadEnv(envPath)
	db := configs.GetConnection()
	t.Cleanup(func() { db.Close() })

	return db
}
