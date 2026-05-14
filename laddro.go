package laddro

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const DefaultBaseURL = "https://api.laddro.com"

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client

	Templates    *TemplatesService
	Resumes      *ResumesService
	Tailor       *TailorService
	CoverLetters *CoverLettersService
	Export       *ExportService
	Settings     *SettingsService
}

type Option func(*Client)

func WithBaseURL(url string) Option {
	return func(c *Client) { c.baseURL = strings.TrimRight(url, "/") }
}

func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.httpClient = hc }
}

func New(apiKey string, opts ...Option) *Client {
	c := &Client{
		baseURL:    DefaultBaseURL,
		apiKey:     apiKey,
		httpClient: http.DefaultClient,
	}
	for _, opt := range opts {
		opt(c)
	}
	c.Templates = &TemplatesService{c}
	c.Resumes = &ResumesService{c}
	c.Tailor = &TailorService{c}
	c.CoverLetters = &CoverLettersService{c}
	c.Export = &ExportService{c}
	c.Settings = &SettingsService{c}
	return c
}

func (c *Client) doJSON(ctx context.Context, method, path string, body any, result any) error {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.apiKey != "" {
		req.Header.Set("x-api-key", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return parseError(resp)
	}

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}

func (c *Client) doBinary(ctx context.Context, method, path string, body any) ([]byte, error) {
	resp, err := c.doBinaryDetailed(ctx, method, path, body)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func (c *Client) doBinaryDetailed(ctx context.Context, method, path string, body any) (*BinaryResponse, error) {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.apiKey != "" {
		req.Header.Set("x-api-key", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, parseError(resp)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &BinaryResponse{
		Data:     data,
		Metadata: artifactMetadata(resp.Header),
	}, nil
}

func (c *Client) doMultipartBinary(ctx context.Context, path string, fields map[string]string, fileName string, fileData io.Reader) ([]byte, error) {
	resp, err := c.doMultipartBinaryDetailed(ctx, path, fields, fileName, fileData)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func (c *Client) doMultipartBinaryDetailed(ctx context.Context, path string, fields map[string]string, fileName string, fileData io.Reader) (*BinaryResponse, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	for k, v := range fields {
		_ = w.WriteField(k, v)
	}

	if fileData != nil {
		part, err := w.CreateFormFile("file", fileName)
		if err != nil {
			return nil, err
		}
		if _, err := io.Copy(part, fileData); err != nil {
			return nil, err
		}
	}
	w.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	if c.apiKey != "" {
		req.Header.Set("x-api-key", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, parseError(resp)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &BinaryResponse{
		Data:     data,
		Metadata: artifactMetadata(resp.Header),
	}, nil
}

func (c *Client) doSSE(ctx context.Context, method, path string, body any) (<-chan SSEEvent, error) {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "text/event-stream")
	if c.apiKey != "" {
		req.Header.Set("x-api-key", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		return nil, parseError(resp)
	}

	ch := make(chan SSEEvent, 16)
	go func() {
		defer resp.Body.Close()
		defer close(ch)

		buf := make([]byte, 4096)
		var remainder string
		var currentEvent string

		for {
			n, err := resp.Body.Read(buf)
			if n > 0 {
				remainder += string(buf[:n])
				for {
					idx := strings.Index(remainder, "\n")
					if idx == -1 {
						break
					}
					line := remainder[:idx]
					remainder = remainder[idx+1:]

					if strings.HasPrefix(line, "event: ") {
						currentEvent = strings.TrimSpace(line[7:])
					} else if strings.HasPrefix(line, "data: ") {
						data := line[6:]
						if currentEvent != "" {
							select {
							case ch <- SSEEvent{Event: currentEvent, Data: data}:
							case <-ctx.Done():
								return
							}
							currentEvent = ""
						}
					}
				}
			}
			if err != nil {
				return
			}
		}
	}()

	return ch, nil
}

func addQuery(path string, params map[string]string) string {
	if len(params) == 0 {
		return path
	}
	v := url.Values{}
	for k, val := range params {
		if val != "" {
			v.Set(k, val)
		}
	}
	if len(v) == 0 {
		return path
	}
	return path + "?" + v.Encode()
}

func artifactMetadata(header http.Header) ArtifactMetadata {
	mimeType := header.Get("Content-Type")
	if mediaType, _, err := mime.ParseMediaType(mimeType); err == nil {
		mimeType = mediaType
	}

	filename := ""
	if _, params, err := mime.ParseMediaType(header.Get("Content-Disposition")); err == nil {
		filename = params["filename"]
	}

	return ArtifactMetadata{
		ResumeID:      header.Get("X-Resume-Id"),
		CoverLetterID: header.Get("X-Cover-Letter-Id"),
		Filename:      filename,
		MimeType:      mimeType,
	}
}

func intParam(n int) string {
	if n == 0 {
		return ""
	}
	return strconv.Itoa(n)
}
