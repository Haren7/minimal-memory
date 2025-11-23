package blobstore

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewS3Client() *s3.Client {
	config, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("ap-south-1"))
	if err != nil {
		log.Printf("[ERROR] NewS3Client: Failed to load default config - %v", err)
		return nil
	}
	return s3.NewFromConfig(config)
}
