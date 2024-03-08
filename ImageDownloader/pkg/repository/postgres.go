package repository

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image/jpeg"
	"sync"

	"github.com/Calgorr/ImageDownloader/ImageDownloader/pkg/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	InsertQuery = "INSERT INTO images(url, byte_array) VALUES($1, $2) ON CONFLICT DO NOTHING"
	BatchSize   = 10
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

type ImageRepository interface {
	InsertImage(ctx context.Context, image []*model.Image) error
}

type ImageRepositoryImpl struct {
	Pgx       *pgxpool.Pool
	BatchSize int
	Batch     *pgx.Batch
}

func NewImageRepository(pg *pgxpool.Pool, autoCreate bool) ImageRepository {
	if autoCreate {
		_, err := pg.Exec(context.Background(), UrlSchema)
		if err != nil {
			panic(err)
		}
	}
	return &ImageRepositoryImpl{
		Pgx:       pg,
		BatchSize: BatchSize,
		Batch:     &pgx.Batch{},
	}
}

func (i *ImageRepositoryImpl) InsertImage(ctx context.Context, image []*model.Image) error {
	if len(image) == 0 {
		return errors.New("no images to insert")
	}
	conn, err := i.Pgx.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	for _, img := range image {
		var buffer = bufferPool.Get().(*bytes.Buffer)
		err := jpeg.Encode(buffer, img.Image, nil)
		if err != nil {
			return err
		}
		imageByteArray := buffer.Bytes()
		i.Batch.Queue(InsertQuery, img.Url, imageByteArray)
	}

	br := conn.SendBatch(ctx, i.Batch)
	err = br.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
