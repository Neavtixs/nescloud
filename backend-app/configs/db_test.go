package configs

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConnection(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	envPath := filepath.Join(filepath.Dir(filename), "..", ".env")

	LoadEnv(envPath)
	db := GetConnection()
	defer db.Close()

	err := db.Ping()
	assert.NoError(t, err)
}
