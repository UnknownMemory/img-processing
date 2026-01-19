package aws

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	Client *s3.Client
}

func NewS3Client() *S3Client {
	accessKey := os.Getenv("BUCKET_ID_KEY")
	secretKey := os.Getenv("BUCKET_SECRET_KEY")

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(os.Getenv("S3_REGION")),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		log.Fatal(err)
	}

	client := s3.NewFromConfig(cfg, func(options *s3.Options) {
		options.BaseEndpoint = aws.String(os.Getenv("S3_ENDPOINT"))
		options.UsePathStyle = true
	})

	return &S3Client{
		Client: client,
	}
}
