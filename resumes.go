package laddro

import (
	"context"
	"io"
)

type ResumesService struct {
	client *Client
}

func (s *ResumesService) List(ctx context.Context, limit, offset int) (*PaginatedList[ResumeSummary], error) {
	path := addQuery("/v1/resumes", map[string]string{
		"limit":  intParam(limit),
		"offset": intParam(offset),
	})
	var resp PaginatedList[ResumeSummary]
	if err := s.client.doJSON(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *ResumesService) Get(ctx context.Context, resumeID string) (*ResumeSummary, error) {
	var resp ResumeSummary
	if err := s.client.doJSON(ctx, "GET", "/v1/resumes/"+resumeID, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *ResumesService) Parse(ctx context.Context, file io.Reader, fileName string, fields map[string]string) ([]byte, error) {
	return s.client.doMultipartBinary(ctx, "/v1/resumes/parse", fields, fileName, file)
}

func (s *ResumesService) Render(ctx context.Context, resumeID string, opts RenderOptions) ([]byte, error) {
	return s.client.doBinary(ctx, "PUT", "/v1/resumes/"+resumeID+"/render", opts)
}
