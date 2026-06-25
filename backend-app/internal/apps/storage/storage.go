package storage

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"nescloud/backend-app/configs"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Storage struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	bucket        string
}

func NewStorage(s3cfg *configs.S3Config) *Storage {
	bucket := os.Getenv("S3_BUCKET")
	if strings.TrimSpace(bucket) == "" {
		bucket = "default"
	}

	return &Storage{
		client:        s3cfg.Client,
		presignClient: s3cfg.PresignClient,
		bucket:        bucket,
	}
}

func (s *Storage) GenerateUploadURL(ctx context.Context, storageKey string, expire time.Duration) (string, time.Time, error) {
	expiredAt := time.Now().Add(expire)

	req, err := s.presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(storageKey),
	}, s3.WithPresignExpires(expire))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("generate upload url: %w", err)
	}

	return req.URL, expiredAt, nil
}

func (s *Storage) GenerateDownloadURL(ctx context.Context, storageKey string, expire time.Duration) (string, error) {
	req, err := s.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(storageKey),
	}, s3.WithPresignExpires(expire))
	if err != nil {
		return "", fmt.Errorf("generate download url: %w", err)
	}

	return req.URL, nil
}

func (s *Storage) ObjectExists(ctx context.Context, storageKey string) (bool, error) {
	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(storageKey),
	})
	if err != nil {
		return false, nil
	}

	return true, nil
}

func (s *Storage) DeleteObject(ctx context.Context, storageKey string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(storageKey),
	})
	if err != nil {
		return fmt.Errorf("delete object: %w", err)
	}

	return nil
}
