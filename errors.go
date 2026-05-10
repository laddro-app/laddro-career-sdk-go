package laddro

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type APIError struct {
	Message string `json:"error"`
	Code    string `json:"code,omitempty"`
	Status  int    `json:"-"`
}

func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("laddro: %s (status %d, code %s)", e.Message, e.Status, e.Code)
	}
	return fmt.Sprintf("laddro: %s (status %d)", e.Message, e.Status)
}

func IsAuthError(err error) bool {
	if e, ok := err.(*APIError); ok {
		return e.Status == 401
	}
	return false
}

func IsUsageLimitError(err error) bool {
	if e, ok := err.(*APIError); ok {
		return e.Status == 402
	}
	return false
}

func IsNotFoundError(err error) bool {
	if e, ok := err.(*APIError); ok {
		return e.Status == 404
	}
	return false
}

func parseError(resp *http.Response) error {
	body, _ := io.ReadAll(resp.Body)
	apiErr := &APIError{Status: resp.StatusCode}
	if err := json.Unmarshal(body, apiErr); err != nil {
		apiErr.Message = string(body)
	}
	if apiErr.Message == "" {
		apiErr.Message = resp.Status
	}
	return apiErr
}
