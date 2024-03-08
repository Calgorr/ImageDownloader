package scaleops

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Calgorr/ImageDownloader/ImageDownloader/handler"
	"github.com/Calgorr/ImageDownloader/ImageDownloader/pkg/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

func connectDB() (*pgxpool.Pool, error) {
	maxConnections, _ := strconv.Atoi(os.Getenv("DB_MAX_CONNECTIONS"))
	fmt.Println(maxConnections)
	config, _ := pgxpool.ParseConfig("")
	config.ConnConfig.Host = os.Getenv("DB_HOST")
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	config.ConnConfig.Port = uint16(port)
	config.ConnConfig.Database = os.Getenv("DB_NAME")
	config.ConnConfig.User = os.Getenv("DB_USER")
	config.ConnConfig.Password = os.Getenv("DB_PASSWORD")
	config.MaxConns = int32(maxConnections)

	return pgxpool.NewWithConfig(context.Background(), config)
}

func Start() {
	pgxPool, err := connectDB()
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
	defer pgxPool.Close()

	imageRepo := repository.NewImageRepository(pgxPool, true)
	fmt.Println("Connected to database")

	processImages(imageRepo)
}

func processImages(imageRepo repository.ImageRepository) {
	imageNumber, _ := strconv.Atoi(os.Args[1])
	queryString := strings.Replace(os.Args[2], " ", "+", -1)
	apiKey := os.Getenv("API_KEY")
	serpApiUrl := os.Getenv("SERP_API_URL")
	fetcher := handler.NewUrlFetcher(imageNumber, queryString, apiKey, serpApiUrl)
	downloader := handler.NewImageDownloader()
	resizer := handler.NewImageResizer(500, 500)
	saver := handler.NewImageSaver(imageRepo)

	urls := fetcher.FetchUrls()
	images := downloader.DownloadImage(urls)
	resizedImages := resizer.ResizeImage(images)
	successful := saver.SaveImage(resizedImages)

	for s := range successful {
		if !s {
			fmt.Println("Error saving images")
		} else {
			fmt.Println("Images saved successfully")
		}
	}
}
