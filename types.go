package laddro

type ArtifactMetadata struct {
	ResumeID      string
	CoverLetterID string
	Filename      string
	MimeType      string
}

type BinaryResponse struct {
	Data     []byte
	Metadata ArtifactMetadata
}

type ResumeSummary struct {
	ID        string `json:"id"`
	ResumeID  string `json:"resumeId"`
	Title     string `json:"title"`
	IsDefault bool   `json:"isDefault"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type PaginatedList[T any] struct {
	Items  []T `json:"items"`
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type Template struct {
	ID                   string           `json:"id"`
	Name                 string           `json:"name"`
	ATSScore             int              `json:"atsScore"`
	LayoutType           string           `json:"layoutType"`
	SupportsProfileImage bool             `json:"supportsProfileImage"`
	Defaults             TemplateDefaults `json:"defaults"`
}

type TemplateDefaults struct {
	PageSize      string `json:"pageSize"`
	Spacing       int    `json:"spacing"`
	FontSize      int    `json:"fontSize"`
	Font          string `json:"font"`
	PageNumbering string `json:"pageNumbering"`
}

type TemplateColor struct {
	ID                  string `json:"id"`
	BackgroundColor     string `json:"backgroundColor"`
	BackgroundPartColor string `json:"backgroundPartColor,omitempty"`
	UnderlineColor      string `json:"underlineColor,omitempty"`
	Text                string `json:"text,omitempty"`
	TextMuted           string `json:"textMuted,omitempty"`
}

type TemplateFont struct {
	Family string `json:"family"`
	Label  string `json:"label"`
}

type TemplateDetail struct {
	Template
	AvailableColors []TemplateColor `json:"availableColors"`
	AvailableFonts  []TemplateFont  `json:"availableFonts"`
}

type Model struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Recommended bool   `json:"recommended"`
}

type ModelProvider struct {
	Provider  string  `json:"provider"`
	Name      string  `json:"name"`
	BaseURL   string  `json:"baseUrl"`
	Models    []Model `json:"models"`
	KeyPrefix string  `json:"keyPrefix"`
	DocsURL   string  `json:"docsUrl"`
}

type Language struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type RenderOptions struct {
	TemplateID       string  `json:"templateId"`
	Locale           string  `json:"locale,omitempty"`
	ColorID          string  `json:"colorId,omitempty"`
	Font             string  `json:"font,omitempty"`
	Spacing          float64 `json:"spacing,omitempty"`
	Margin           float64 `json:"margin,omitempty"`
	FontSize         float64 `json:"fontSize,omitempty"`
	ShowProfileImage *bool   `json:"showProfileImage,omitempty"`
	ProfileImageURL  string  `json:"profileImageUrl,omitempty"`
	PageNumbering    string  `json:"pageNumbering,omitempty"`
}

type TailorRequest struct {
	ResumeID           string  `json:"resumeId,omitempty"`
	PositionName       string  `json:"positionName"`
	JobDescription     string  `json:"jobDescription,omitempty"`
	JobURL             string  `json:"jobUrl,omitempty"`
	Mode               string  `json:"mode,omitempty"`
	Language           string  `json:"language,omitempty"`
	IncludeCoverLetter *bool   `json:"includeCoverLetter,omitempty"`
	TemplateID         string  `json:"templateId,omitempty"`
	ColorID            string  `json:"colorId,omitempty"`
	Font               string  `json:"font,omitempty"`
	Spacing            float64 `json:"spacing,omitempty"`
	Margin             float64 `json:"margin,omitempty"`
	FontSize           float64 `json:"fontSize,omitempty"`
	PageNumbering      string  `json:"pageNumbering,omitempty"`
}

type CoverLetterSummary struct {
	ID            string         `json:"id"`
	CoverLetterID string         `json:"coverLetterId"`
	Title         string         `json:"title"`
	LetterContent string         `json:"letterContent,omitempty"`
	Data          map[string]any `json:"data,omitempty"`
	CreatedAt     string         `json:"createdAt"`
	UpdatedAt     string         `json:"updatedAt"`
}

type CreateCoverLetterRequest struct {
	Title         string `json:"title,omitempty"`
	FullName      string `json:"fullName"`
	JobTitle      string `json:"jobTitle,omitempty"`
	Address       string `json:"address,omitempty"`
	Email         string `json:"email,omitempty"`
	Phone         string `json:"phone,omitempty"`
	CompanyName   string `json:"companyName,omitempty"`
	HiringManager string `json:"hiringManager,omitempty"`
	LetterContent string `json:"letterContent"`
}

type CreateCoverLetterResponse struct {
	CoverLetterID string `json:"coverLetterId"`
	Title         string `json:"title"`
	Status        string `json:"status"`
}

type GenerateCoverLetterRequest struct {
	ResumeID       string  `json:"resumeId,omitempty"`
	PositionName   string  `json:"positionName"`
	JobDescription string  `json:"jobDescription,omitempty"`
	JobURL         string  `json:"jobUrl,omitempty"`
	Language       string  `json:"language,omitempty"`
	TemplateID     string  `json:"templateId,omitempty"`
	ColorID        string  `json:"colorId,omitempty"`
	Font           string  `json:"font,omitempty"`
	Spacing        float64 `json:"spacing,omitempty"`
	Margin         float64 `json:"margin,omitempty"`
	FontSize       float64 `json:"fontSize,omitempty"`
	PageNumbering  string  `json:"pageNumbering,omitempty"`
}

type ExportRequest struct {
	ResumeID         string  `json:"resumeId"`
	TemplateID       string  `json:"templateId,omitempty"`
	Locale           string  `json:"locale,omitempty"`
	ColorID          string  `json:"colorId,omitempty"`
	Font             string  `json:"font,omitempty"`
	Spacing          float64 `json:"spacing,omitempty"`
	Margin           float64 `json:"margin,omitempty"`
	FontSize         float64 `json:"fontSize,omitempty"`
	ShowProfileImage *bool   `json:"showProfileImage,omitempty"`
	ProfileImageURL  string  `json:"profileImageUrl,omitempty"`
	PageNumbering    string  `json:"pageNumbering,omitempty"`
}

type AISettings struct {
	Provider  string `json:"provider"`
	Model     string `json:"model"`
	BaseURL   string `json:"baseUrl"`
	HasAPIKey bool   `json:"hasApiKey"`
	UpdatedAt string `json:"updatedAt,omitempty"`
}

type SettingsResponse struct {
	AI *AISettings `json:"ai"`
}

type UpdateAISettingsRequest struct {
	Provider string `json:"provider"`
	Model    string `json:"model,omitempty"`
	APIKey   string `json:"apiKey"`
}

type SSEEvent struct {
	Event string
	Data  string
}
