package repository

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/Calgorr/ImageDownloader/ImageDownloader/pkg/model"
	"github.com/c2fo/testify/assert"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pgxz *pgxpool.Pool
)

func TestMain(m *testing.M) {
	maxConnections, _ := strconv.Atoi(os.Getenv("DB_MAX_CONNECTIONS"))
	config, _ := pgxpool.ParseConfig("")
	config.ConnConfig.Host = os.Getenv("DB_HOST")
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	config.ConnConfig.Port = uint16(port)
	config.ConnConfig.Database = os.Getenv("DB_NAME")
	config.ConnConfig.User = os.Getenv("DB_USER")
	config.ConnConfig.Password = os.Getenv("DB_PASSWORD")
	config.MaxConns = int32(maxConnections)

	pgxz, _ = pgxpool.NewWithConfig(context.Background(), config)
}

func TestInsertImage(t *testing.T) {

	// Initialize ImageRepositoryImpl with mock values
	imageRepo := &ImageRepositoryImpl{
		Pgx:       pgxz,
		BatchSize: 10,
		Batch:     &pgx.Batch{},
	}

	// Mock image data
	testImage := &model.Image{
		Url:   "test_url",
		Image: nil,
	}

	// Create a context
	ctx := context.Background()

	// Test inserting a single image
	err := imageRepo.InsertImage(ctx, []*model.Image{testImage})
	assert.NoError(t, err)
}
