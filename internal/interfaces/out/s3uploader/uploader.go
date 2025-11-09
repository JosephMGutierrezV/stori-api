package s3uploader

import (
	"bytes"
	"context"

	"stori-api/internal/core/ports/out"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

type S3Uploader struct {
	client Client
}

var _ out.ObjectStorage = (*S3Uploader)(nil)

func NewS3Uploader(client Client) *S3Uploader {
	return &S3Uploader{client: client}
}

func (u *S3Uploader) PutObject(
	ctx context.Context,
	bucket string,
	key string,
	contentType string,
	data []byte,
) error {
	_, err := u.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	return err
}
