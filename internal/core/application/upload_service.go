package application

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	inport "stori-api/internal/core/ports/in"
	"stori-api/internal/core/ports/out"

	"github.com/google/uuid"
)

type CSVUploadService struct {
	storage  out.ObjectStorage
	bucket   string
	basePath string
}

func NewCSVUploadService(storage out.ObjectStorage, bucket string, basePath string) *CSVUploadService {
	return &CSVUploadService{
		storage:  storage,
		bucket:   bucket,
		basePath: strings.TrimSuffix(basePath, "/") + "/",
	}
}

var _ inport.CSVUploadPort = (*CSVUploadService)(nil)

func (s *CSVUploadService) UploadCSV(
	ctx context.Context,
	req inport.CSVUploadRequest,
) (inport.CSVUploadResult, error) {
	if len(req.RawBody) == 0 {
		return inport.CSVUploadResult{}, errors.New("body is empty")
	}

	if err := validateTransactionsCSV(req.RawBody); err != nil {
		return inport.CSVUploadResult{}, fmt.Errorf("invalid CSV: %w", err)
	}

	key := fmt.Sprintf("%s%s.csv", s.basePath, uuid.New().String())

	if err := s.storage.PutObject(ctx, s.bucket, key, "text/csv", req.RawBody); err != nil {
		return inport.CSVUploadResult{}, fmt.Errorf("upload to storage: %w", err)
	}

	return inport.CSVUploadResult{
		Bucket: s.bucket,
		Key:    key,
	}, nil
}

func validateTransactionsCSV(data []byte) error {
	r := csv.NewReader(bytes.NewReader(data))

	header, err := r.Read()
	if err != nil {
		return fmt.Errorf("reading header: %w", err)
	}

	if len(header) < 3 {
		return errors.New("header must have at least 3 columns")
	}

	if strings.TrimSpace(header[0]) != "Id" ||
		strings.TrimSpace(header[1]) != "Date" ||
		strings.TrimSpace(header[2]) != "Transaction" {
		return fmt.Errorf("invalid header, expected 'Id,Date,Transaction', got %q", header)
	}

	rowCount := 0

	for {
		record, err := r.Read()
		if err != nil {
			if errors.Is(err, csv.ErrFieldCount) {
				return fmt.Errorf("row %d has invalid field count", rowCount+1)
			}
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("reading row %d: %w", rowCount+1, err)
		}
		rowCount++

		if len(record) < 3 {
			return fmt.Errorf("row %d has less than 3 columns", rowCount)
		}

		idStr := strings.TrimSpace(record[0])
		dateStr := strings.TrimSpace(record[1])
		amountStr := strings.TrimSpace(record[2])

		if _, err := strconv.Atoi(idStr); err != nil {
			return fmt.Errorf("row %d: invalid Id %q", rowCount, idStr)
		}

		if _, err := time.Parse("1/2", dateStr); err != nil {
			return fmt.Errorf("row %d: invalid Date %q, expected M/D like 7/15", rowCount, dateStr)
		}

		if _, err := strconv.ParseFloat(amountStr, 64); err != nil {
			return fmt.Errorf("row %d: invalid Transaction %q", rowCount, amountStr)
		}
	}

	if rowCount == 0 {
		return errors.New("CSV has no data rows")
	}

	return nil
}
