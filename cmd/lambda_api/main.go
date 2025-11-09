package main

import (
	"context"
	"log"

	"stori-api/internal/infra/bootstrap"
	"stori-api/internal/infra/config"
	"stori-api/internal/infra/logger"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
)

var uploadCtx *bootstrap.UploadAPIContext

func init() {
	if err := logger.Init(); err != nil {
		log.Fatalf("error iniciando logger: %v", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Logger.Fatal("error cargando configuración", zap.Error(err))
	}

	logger.Logger.Info("config upload API cargada",
		zap.String("s3_bucket", cfg.S3BucketName),
		zap.String("s3_region", cfg.S3Region),
	)

	uploadCtx, err = bootstrap.InitializeUploadAPI(cfg)
	if err != nil {
		logger.Logger.Fatal("error inicializando upload API", zap.Error(err))
	}
}

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return uploadCtx.Handler.Handle(ctx, req)
}

func main() {
	defer logger.Sync()
	log.Println("Lambda Upload API de Stori iniciando...")
	lambda.Start(handler)
}
