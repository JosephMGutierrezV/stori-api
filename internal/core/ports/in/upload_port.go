package in

import (
	"context"
)

type CSVUploadRequest struct {
	RawBody     []byte
	ContentType string
}

type CSVUploadResult struct {
	Bucket string
	Key    string
}

type CSVUploadPort interface {
	UploadCSV(ctx context.Context, req CSVUploadRequest) (CSVUploadResult, error)
}
