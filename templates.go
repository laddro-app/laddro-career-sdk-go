package laddro

import "context"

type TemplatesService struct {
	client *Client
}

func (s *TemplatesService) List(ctx context.Context) ([]Template, error) {
	var resp struct {
		Templates []Template `json:"templates"`
	}
	if err := s.client.doJSON(ctx, "GET", "/v1/templates", nil, &resp); err != nil {
		return nil, err
	}
	return resp.Templates, nil
}

func (s *TemplatesService) Get(ctx context.Context, templateID string) (*TemplateDetail, error) {
	var resp TemplateDetail
	if err := s.client.doJSON(ctx, "GET", "/v1/templates/"+templateID, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *TemplatesService) Fonts(ctx context.Context) ([]TemplateFont, error) {
	var resp struct {
		Fonts []TemplateFont `json:"fonts"`
	}
	if err := s.client.doJSON(ctx, "GET", "/v1/fonts", nil, &resp); err != nil {
		return nil, err
	}
	return resp.Fonts, nil
}

func (s *TemplatesService) Languages(ctx context.Context) ([]Language, error) {
	var resp struct {
		Languages []Language `json:"languages"`
	}
	if err := s.client.doJSON(ctx, "GET", "/v1/languages", nil, &resp); err != nil {
		return nil, err
	}
	return resp.Languages, nil
}

func (s *TemplatesService) Models(ctx context.Context) ([]ModelProvider, error) {
	var resp struct {
		Models []ModelProvider `json:"models"`
	}
	if err := s.client.doJSON(ctx, "GET", "/v1/models", nil, &resp); err != nil {
		return nil, err
	}
	return resp.Models, nil
}
