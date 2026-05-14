package laddro

import (
	"context"
	"io"
)

type CoverLettersService struct {
	client *Client
}

func (s *CoverLettersService) List(ctx context.Context, limit, offset int) (*PaginatedList[CoverLetterSummary], error) {
	path := addQuery("/v1/cover-letters", map[string]string{
		"limit":  intParam(limit),
		"offset": intParam(offset),
	})
	var resp PaginatedList[CoverLetterSummary]
	if err := s.client.doJSON(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *CoverLettersService) Get(ctx context.Context, id string) (*CoverLetterSummary, error) {
	var resp CoverLetterSummary
	if err := s.client.doJSON(ctx, "GET", "/v1/cover-letters/"+id, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *CoverLettersService) Create(ctx context.Context, req CreateCoverLetterRequest) (*CreateCoverLetterResponse, error) {
	var resp CreateCoverLetterResponse
	if err := s.client.doJSON(ctx, "POST", "/v1/cover-letters", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *CoverLettersService) Generate(ctx context.Context, req GenerateCoverLetterRequest) ([]byte, error) {
	return s.client.doBinary(ctx, "POST", "/v1/cover-letters/generate", req)
}

func (s *CoverLettersService) GenerateDetailed(ctx context.Context, req GenerateCoverLetterRequest) (*BinaryResponse, error) {
	return s.client.doBinaryDetailed(ctx, "POST", "/v1/cover-letters/generate", req)
}

func (s *CoverLettersService) Upload(ctx context.Context, file io.Reader, fileName string, fields map[string]string) ([]byte, error) {
	return s.client.doMultipartBinary(ctx, "/v1/cover-letters/generate", fields, fileName, file)
}

func (s *CoverLettersService) UploadDetailed(ctx context.Context, file io.Reader, fileName string, fields map[string]string) (*BinaryResponse, error) {
	return s.client.doMultipartBinaryDetailed(ctx, "/v1/cover-letters/generate", fields, fileName, file)
}

func (s *CoverLettersService) Render(ctx context.Context, id string, opts RenderOptions) ([]byte, error) {
	return s.client.doBinary(ctx, "PUT", "/v1/cover-letters/"+id+"/render", opts)
}

func (s *CoverLettersService) GenerateStream(ctx context.Context, req GenerateCoverLetterRequest) (<-chan SSEEvent, error) {
	return s.client.doSSE(ctx, "POST", "/v1/cover-letters/generate", req)
}
