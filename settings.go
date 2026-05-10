package laddro

import "context"

type SettingsService struct {
	client *Client
}

func (s *SettingsService) Get(ctx context.Context) (*SettingsResponse, error) {
	var resp SettingsResponse
	if err := s.client.doJSON(ctx, "GET", "/v1/settings", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *SettingsService) UpdateModel(ctx context.Context, req UpdateAISettingsRequest) (*SettingsResponse, error) {
	var resp SettingsResponse
	if err := s.client.doJSON(ctx, "PUT", "/v1/settings/model", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *SettingsService) DeleteModel(ctx context.Context) (*SettingsResponse, error) {
	var resp SettingsResponse
	if err := s.client.doJSON(ctx, "DELETE", "/v1/settings/model", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
