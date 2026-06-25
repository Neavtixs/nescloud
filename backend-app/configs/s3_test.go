package configs

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewS3(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	envPath := filepath.Join(filepath.Dir(filename), "..", ".env")
	LoadEnv(envPath)

	s3cfg := NewS3()
	assert.NotNil(t, s3cfg.Client)
	assert.NotNil(t, s3cfg.PresignClient)
	assert.NotEmpty(t, s3cfg.Bucket)
}

func TestS3Connection(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	envPath := filepath.Join(filepath.Dir(filename), "..", ".env")
	LoadEnv(envPath)

	s3cfg := NewS3()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s3cfg.HealthCheck(ctx)
	assert.NoError(t, err)
}
