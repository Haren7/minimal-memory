package blobstore

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type BlobStoreInterface interface {
	Store(ctx context.Context, bucket string, path string, files []os.File) error
	Retrieve(ctx context.Context, bucket string, path string) (map[string]io.Reader, error)
}

type s3Store struct {
	s3Client *s3.Client
}

func NewS3Store(s3Client *s3.Client) BlobStoreInterface {
	return &s3Store{
		s3Client: s3Client,
	}
}

func (r *s3Store) Store(ctx context.Context, bucket string, path string, files []os.File) error {
	for _, file := range files {
		key := filepath.Join(path, file.Name())
		_, err := r.s3Client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: &bucket,
			Key:    &key,
			Body:   &file,
		})
		if err != nil {
			return fmt.Errorf("error uploading file %s: %w", key, err)
		}
	}
	return nil
}

func (r *s3Store) Retrieve(ctx context.Context, bucket string, path string) (map[string]io.Reader, error) {
	// get files at path
	files, err := r.s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: &bucket,
		Prefix: &path,
	})
	if err != nil {
		return nil, fmt.Errorf("error listing objects: %w", err)
	}

	filesMap := make(map[string]io.Reader)
	for _, object := range files.Contents {
		file, err := r.s3Client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: &bucket,
			Key:    object.Key,
		})
		if err != nil {
			return nil, fmt.Errorf("error getting object: %w", err)
		}
		// map of file name vs reader
		filesMap[strings.TrimPrefix(*object.Key, path)] = file.Body
	}
	return filesMap, nil
}
