package snapshot

import (
	"context"
	"fmt"
	"memory/internal/blobstore"
	"memory/internal/persistence/vector"
)

type faissManager struct {
	dir         string
	bucket      string
	s3          blobstore.BlobStoreInterface
	faissClient *vector.FaissClient
}

func NewFaissManager(bucket string, s3 blobstore.BlobStoreInterface, faissClient *vector.FaissClient) Manager {
	return &faissManager{
		dir:         "/faiss",
		bucket:      bucket,
		s3:          s3,
		faissClient: faissClient,
	}
}

func (r *faissManager) Store(ctx context.Context) error {
	files, err := r.faissClient.Export(r.dir)
	if err != nil {
		return fmt.Errorf("snapshot: error exporting faiss: %w", err)
	}
	err = r.s3.Store(ctx, r.bucket, r.dir, files)
	if err != nil {
		return fmt.Errorf("snapshot: error storing faiss: %w", err)
	}
	return nil
}

func (r *faissManager) Load(ctx context.Context) error {
	files, err := r.s3.Retrieve(ctx, r.bucket, r.dir)
	if err != nil {
		return fmt.Errorf("snapshot: error retrieving faiss: %w", err)
	}
	err = r.faissClient.Mount(r.dir, files)
	if err != nil {
		return fmt.Errorf("snapshot: error mounting faiss: %w", err)
	}
	return nil
}
