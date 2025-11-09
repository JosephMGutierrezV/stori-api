package apigw

import (
	"context"
	"encoding/base64"
	"encoding/json"
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

func (h *UploadHandler) Handle(
	ctx context.Context,
	req events.APIGatewayV2HTTPRequest,
) (events.APIGatewayV2HTTPResponse, error) {
	bodyBytes, err := getBodyBytes(req)
	if err != nil {
		return errorResponse(http.StatusBadRequest, "invalid body"), nil
	}
	if len(bodyBytes) == 0 {
		return errorResponse(http.StatusBadRequest, "body is empty"), nil
	}

	contentType := strings.ToLower(req.Headers["content-type"])

	result, err := h.useCase.UploadCSV(ctx, inport.CSVUploadRequest{
		RawBody:     bodyBytes,
		ContentType: contentType,
	})
	if err != nil {
		return errorResponse(http.StatusBadRequest, err.Error()), nil
	}

	resp := map[string]string{
		"message": "file uploaded successfully",
		"bucket":  result.Bucket,
		"key":     result.Key,
	}
	body, _ := json.Marshal(resp)

	return events.APIGatewayV2HTTPResponse{
		StatusCode:      http.StatusCreated,
		Body:            string(body),
		Headers:         map[string]string{"Content-Type": "application/json"},
		IsBase64Encoded: false,
	}, nil
}

func getBodyBytes(req events.APIGatewayV2HTTPRequest) ([]byte, error) {
	if req.Body == "" {
		return nil, nil
	}
	if !req.IsBase64Encoded {
		return []byte(req.Body), nil
	}
	return base64.StdEncoding.DecodeString(req.Body)
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
