package laddro

import "context"

type ExportService struct {
	client *Client
}

func (s *ExportService) PDF(ctx context.Context, req ExportRequest) ([]byte, error) {
	return s.client.doBinary(ctx, "POST", "/v1/export", req)
}

func (s *ExportService) Stream(ctx context.Context, req ExportRequest) (<-chan SSEEvent, error) {
	return s.client.doSSE(ctx, "POST", "/v1/export", req)
}
