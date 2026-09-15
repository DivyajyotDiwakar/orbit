package utils

import (
	"Orbit/configs"
	"Orbit/internal/db"
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	MaxEventSize   = 5 * 1024 * 1024 // 5MB
	MaxProductSize = 4 * 1024 * 1024 // 4MB
)

type Uploader interface {
	UploadPhoto(ctx context.Context, fileHeader *multipart.FileHeader, id string) (string, error)
	DeletePhoto(ctx context.Context, folder string, filename string) error
	UpdatePhoto(ctx context.Context, fileHeader *multipart.FileHeader, id string, oldFolder string, oldFilename string) (string, error)
}

type Event struct{}
type Product struct{}

func validateAndDecode(fileHeader *multipart.FileHeader, maxSize int64) (image.Image, string, error) {
	if fileHeader == nil {
		return nil, "", errors.New("image file is required")
	}

	if fileHeader.Size > maxSize {
		return nil, "", fmt.Errorf("file too large: maximum allowed size is %dMB", maxSize/(1024*1024))
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer file.Close()

	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, "", fmt.Errorf("failed to read file header: %w", err)
	}

	contentType := http.DetectContentType(buf[:n])
	if contentType != "image/jpeg" && contentType != "image/png" {
		return nil, "", errors.New("invalid format: only JPEG and PNG images are accepted")
	}

	if _, err = file.Seek(0, io.SeekStart); err != nil {
		return nil, "", fmt.Errorf("failed to reset file pointer for dimension check: %w", err)
	}

	cfg, _, err := image.DecodeConfig(file)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read image dimensions: %w", err)
	}
	if cfg.Width > 5000 || cfg.Height > 5000 {
		return nil, "", fmt.Errorf("image too large: dimensions %dx%d exceed the 5000×5000 limit", cfg.Width, cfg.Height)
	}

	if _, err = file.Seek(0, io.SeekStart); err != nil {
		return nil, "", fmt.Errorf("failed to reset file pointer for decode: %w", err)
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, "", errors.New("image is corrupted or could not be decoded")
	}

	return img, contentType, nil
}

func uploadToR2Internal(ctx context.Context, folder, filename, contentType string, body io.Reader) (string, error) {
	r2Cfg := configs.LoadConfig().R2
	client := db.GetR2Client()

	key := fmt.Sprintf("%s/%s", folder, filename)

	_, err := client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r2Cfg.Bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("R2 upload failed: %w", err)
	}

	return fmt.Sprintf("%s/%s", r2Cfg.Domain, key), nil
}

func deleteFromR2Internal(ctx context.Context, folder, filename string) error {
	r2Cfg := configs.LoadConfig().R2
	key := fmt.Sprintf("%s/%s", folder, filename)

	_, err := db.GetR2Client().DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(r2Cfg.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("R2 delete failed for %s: %w", key, err)
	}
	return nil
}

func encodeAndUpload(ctx context.Context, img image.Image, contentType, folder, id string, quality int) (string, error) {
	buf := new(bytes.Buffer)
	var filename string
	var encErr error

	stamp := time.Now().UnixNano()

	if contentType == "image/png" {
		filename = fmt.Sprintf("%s-%d.png", id, stamp)
		encErr = png.Encode(buf, img)
	} else {
		filename = fmt.Sprintf("%s-%d.jpg", id, stamp)
		encErr = jpeg.Encode(buf, img, &jpeg.Options{Quality: quality})
		contentType = "image/jpeg"
	}

	if encErr != nil {
		return "", fmt.Errorf("image encoding failed: %w", encErr)
	}

	return uploadToR2Internal(ctx, folder, filename, contentType, buf)
}

func (e *Event) UploadPhoto(ctx context.Context, fileHeader *multipart.FileHeader, id string) (string, error) {
	img, contentType, err := validateAndDecode(fileHeader, MaxEventSize)
	if err != nil {
		return "", err
	}
	return encodeAndUpload(ctx, img, contentType, "event-banner", id, 85)
}

func (e *Event) DeletePhoto(ctx context.Context, folder string, filename string) error {
	return deleteFromR2Internal(ctx, folder, filename)
}

func (e *Event) UpdatePhoto(ctx context.Context, fileHeader *multipart.FileHeader, id string, oldFolder string, oldFilename string) (string, error) {
	newURL, err := e.UploadPhoto(ctx, fileHeader, id)
	if err != nil {
		return "", err
	}
	if oldFilename != "" && oldFolder != "" {
		if delErr := deleteFromR2Internal(ctx, oldFolder, oldFilename); delErr != nil {
			log.Printf("WARN orphan file — event banner not deleted %s/%s: %v", oldFolder, oldFilename, delErr)
		}
	}
	return newURL, nil
}

func (p *Product) UploadPhoto(ctx context.Context, fileHeader *multipart.FileHeader, id string) (string, error) {
	img, contentType, err := validateAndDecode(fileHeader, MaxProductSize)
	if err != nil {
		return "", err
	}
	return encodeAndUpload(ctx, img, contentType, "product", id, 80)
}

func (p *Product) DeletePhoto(ctx context.Context, folder string, filename string) error {
	return deleteFromR2Internal(ctx, folder, filename)
}

func (p *Product) UpdatePhoto(ctx context.Context, fileHeader *multipart.FileHeader, id string, oldFolder string, oldFilename string) (string, error) {
	newURL, err := p.UploadPhoto(ctx, fileHeader, id)
	if err != nil {
		return "", err
	}
	if oldFilename != "" && oldFolder != "" {
		if delErr := deleteFromR2Internal(ctx, oldFolder, oldFilename); delErr != nil {
			log.Printf("WARN orphan file — product image not deleted %s/%s: %v", oldFolder, oldFilename, delErr)
		}
	}
	return newURL, nil
}
