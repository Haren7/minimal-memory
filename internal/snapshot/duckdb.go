package snapshot

import (
	"context"
	"fmt"
	"memory/internal/blobstore"
	"memory/internal/persistence/rdbms"
)

type duckdbManager struct {
	dir          string
	bucket       string
	s3           blobstore.BlobStoreInterface
	duckdbClient *rdbms.DuckDBClient
}

func NewDuckdbManager(bucket string, s3 blobstore.BlobStoreInterface, duckdbClient *rdbms.DuckDBClient) Manager {
	return &duckdbManager{
		dir:          "/duckdb",
		bucket:       bucket,
		s3:           s3,
		duckdbClient: duckdbClient,
	}
}

func (r *duckdbManager) Store(ctx context.Context) error {
	files, err := r.duckdbClient.Export(r.dir)
	if err != nil {
		return fmt.Errorf("snapshot: error exporting duckdb: %w", err)
	}
	err = r.s3.Store(ctx, r.bucket, r.dir, files)
	if err != nil {
		return fmt.Errorf("snapshot: error storing duckdb: %w", err)
	}
	return nil
}

func (r *duckdbManager) Load(ctx context.Context) error {
	files, err := r.s3.Retrieve(ctx, r.bucket, r.dir)
	if err != nil {
		return fmt.Errorf("snapshot: error retrieving duckdb: %w", err)
	}
	err = r.duckdbClient.Mount(r.dir, files)
	if err != nil {
		return fmt.Errorf("snapshot: error mounting duckdb: %w", err)
	}
	return nil
}
