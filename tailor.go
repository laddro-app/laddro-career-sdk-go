package laddro

import (
	"context"
	"io"
)

type TailorService struct {
	client *Client
}

func (s *TailorService) Run(ctx context.Context, req TailorRequest) ([]byte, error) {
	return s.client.doBinary(ctx, "POST", "/v1/tailor", req)
}

func (s *TailorService) RunDetailed(ctx context.Context, req TailorRequest) (*BinaryResponse, error) {
	return s.client.doBinaryDetailed(ctx, "POST", "/v1/tailor", req)
}

func (s *TailorService) Upload(ctx context.Context, file io.Reader, fileName string, fields map[string]string) ([]byte, error) {
	return s.client.doMultipartBinary(ctx, "/v1/tailor", fields, fileName, file)
}

func (s *TailorService) UploadDetailed(ctx context.Context, file io.Reader, fileName string, fields map[string]string) (*BinaryResponse, error) {
	return s.client.doMultipartBinaryDetailed(ctx, "/v1/tailor", fields, fileName, file)
}

func (s *TailorService) Stream(ctx context.Context, req TailorRequest) (<-chan SSEEvent, error) {
	return s.client.doSSE(ctx, "POST", "/v1/tailor", req)
}
