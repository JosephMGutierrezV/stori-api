package apigw

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"

	inport "stori-api/internal/core/ports/in"

	"github.com/aws/aws-lambda-go/events"
)

type UploadHandler struct {
	useCase inport.CSVUploadPort
}

func NewUploadHandler(useCase inport.CSVUploadPort) *UploadHandler {
	return &UploadHandler{useCase: useCase}
}

func (h *UploadHandler) Handle(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	contentType := req.Headers["content-type"]

	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return errorResponse(http.StatusBadRequest, "invalid content-type"), nil
	}

	if strings.HasPrefix(mediaType, "multipart/") {
		body, err := getBodyBytes(req)
		if err != nil {
			return errorResponse(http.StatusBadRequest, "cannot read body"), nil
		}

		reader := multipart.NewReader(bytes.NewReader(body), params["boundary"])
		part, err := reader.NextPart()
		if err != nil {
			return errorResponse(http.StatusBadRequest, "no file found"), nil
		}

		fileBytes, err := io.ReadAll(part)
		if err != nil {
			return errorResponse(http.StatusInternalServerError, "cannot read file"), nil
		}

		result, err := h.useCase.UploadCSV(ctx, inport.CSVUploadRequest{
			RawBody:     fileBytes,
			ContentType: part.Header.Get("Content-Type"),
		})
		if err != nil {
			return errorResponse(http.StatusBadRequest, err.Error()), nil
		}

		resp := map[string]string{
			"message": "file uploaded successfully",
			"bucket":  result.Bucket,
			"key":     result.Key,
		}
		bodyResp, _ := json.Marshal(resp)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusCreated,
			Body:       string(bodyResp),
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	return h.handleRawCSV(ctx, req)
}

func (h *UploadHandler) handleRawCSV(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	body, err := getBodyBytes(req)
	if err != nil {
		return errorResponse(http.StatusBadRequest, "cannot read body"), nil
	}

	result, err := h.useCase.UploadCSV(ctx, inport.CSVUploadRequest{
		RawBody:     body,
		ContentType: req.Headers["content-type"],
	})
	if err != nil {
		return errorResponse(http.StatusBadRequest, err.Error()), nil
	}

	resp := map[string]string{
		"message": "file uploaded successfully",
		"bucket":  result.Bucket,
		"key":     result.Key,
	}
	bodyResp, _ := json.Marshal(resp)
	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusCreated,
		Body:       string(bodyResp),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func getBodyBytes(req events.APIGatewayV2HTTPRequest) ([]byte, error) {
	if req.Body == "" {
		return nil, nil
	}
	if req.IsBase64Encoded {
		return base64.StdEncoding.DecodeString(req.Body)
	}
	return []byte(req.Body), nil
}

func errorResponse(status int, msg string) events.APIGatewayV2HTTPResponse {
	body, _ := json.Marshal(map[string]string{"error": msg})
	return events.APIGatewayV2HTTPResponse{
		StatusCode:      status,
		Body:            string(body),
		Headers:         map[string]string{"Content-Type": "application/json"},
		IsBase64Encoded: false,
	}
}
