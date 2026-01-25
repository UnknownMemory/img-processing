package aws

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	Client          *s3.Client
	BucketURL       string
	BucketPublicURL string
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
		log.Println(err)
		return nil
	}

	client := s3.NewFromConfig(cfg, func(options *s3.Options) {
		options.BaseEndpoint = aws.String(os.Getenv("S3_ENDPOINT"))
		options.UsePathStyle = true
	})

	return &S3Client{
		Client:          client,
		BucketURL:       fmt.Sprintf("%s/%s", os.Getenv("S3_ENDPOINT"), os.Getenv("S3_BUCKET")),
		BucketPublicURL: fmt.Sprintf("%s", os.Getenv("S3_PUBLIC_ENDPOINT")),
	}
}

func (s3Client *S3Client) Upload(objectName string, file io.ReadSeeker, contentType string) (*s3.PutObjectOutput, error) {
	object, err := s3Client.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(os.Getenv("S3_BUCKET")),
		Key:         aws.String(objectName),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return object, nil
}
