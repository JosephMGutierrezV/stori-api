package bootstrap

import (
	"context"
	"stori-api/internal/core/application"
	"stori-api/internal/infra/aws/s3client"
	"stori-api/internal/infra/config"
	"stori-api/internal/interfaces/in/apigw"
	"stori-api/internal/interfaces/out/s3uploader"
)

type UploadAPIContext struct {
	Handler *apigw.UploadHandler
}

func InitializeUploadAPI(cfg *config.Config) (*UploadAPIContext, error) {
	awsS3Client, err := s3client.NewS3Client(context.Background(), cfg)
	if err != nil {
		return nil, err
	}

	storage := s3uploader.NewS3Uploader(awsS3Client)

	useCase := application.NewCSVUploadService(
		storage,
		cfg.S3BucketName,
		"uploads",
	)

	handler := apigw.NewUploadHandler(useCase)

	return &UploadAPIContext{
		Handler: handler,
	}, nil
}
