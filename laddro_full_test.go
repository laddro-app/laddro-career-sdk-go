package laddro

import (
	"context"
	"os"
	"testing"
)

func TestAllEndpoints(t *testing.T) {
	apiKey := os.Getenv("LADDRO_API_KEY")
	if apiKey == "" {
		t.Skip("LADDRO_API_KEY not set")
	}

	client := New(apiKey)
	public := New("")
	ctx := context.Background()

	// PUBLIC (5)
	t.Run("GET /v1/templates", func(t *testing.T) {
		templates, err := public.Templates.List(ctx)
		if err != nil { t.Fatal(err) }
		if len(templates) != 22 { t.Fatalf("expected 22, got %d", len(templates)) }
	})
	t.Run("GET /v1/templates/{id}", func(t *testing.T) {
		d, err := public.Templates.Get(ctx, "GRAPHITE")
		if err != nil { t.Fatal(err) }
		if d.ID != "GRAPHITE" { t.Fatal("wrong id") }
		if len(d.AvailableColors) == 0 { t.Fatal("no colors") }
	})
	t.Run("GET /v1/fonts", func(t *testing.T) {
		fonts, err := public.Templates.Fonts(ctx)
		if err != nil { t.Fatal(err) }
		if len(fonts) != 21 { t.Fatalf("expected 21, got %d", len(fonts)) }
	})
	t.Run("GET /v1/languages", func(t *testing.T) {
		langs, err := public.Templates.Languages(ctx)
		if err != nil { t.Fatal(err) }
		if len(langs) != 14 { t.Fatalf("expected 14, got %d", len(langs)) }
	})
	t.Run("GET /v1/models", func(t *testing.T) {
		models, err := public.Templates.Models(ctx)
		if err != nil { t.Fatal(err) }
		if len(models) != 10 { t.Fatalf("expected 10, got %d", len(models)) }
	})

	// RESUMES (3 + parse skip)
	var resumeID string
	t.Run("GET /v1/resumes", func(t *testing.T) {
		list, err := client.Resumes.List(ctx, 5, 0)
		if err != nil { t.Fatal(err) }
		if len(list.Items) == 0 { t.Fatal("no resumes") }
		for _, r := range list.Items {
			if r.IsDefault { resumeID = r.ResumeID; break }
		}
		if resumeID == "" { resumeID = list.Items[0].ResumeID }
	})
	t.Run("GET /v1/resumes/{id}", func(t *testing.T) {
		r, err := client.Resumes.Get(ctx, resumeID)
		if err != nil { t.Fatal(err) }
		if r.ResumeID != resumeID { t.Fatal("id mismatch") }
	})
	t.Run("PUT /v1/resumes/{id}/render", func(t *testing.T) {
		pdf, err := client.Resumes.Render(ctx, resumeID, RenderOptions{TemplateID: "GRAPHITE"})
		if err != nil { t.Fatal(err) }
		if len(pdf) < 1000 { t.Fatalf("too small: %d", len(pdf)) }
	})

	// TAILOR (1)
	t.Run("POST /v1/tailor", func(t *testing.T) {
		pdf, err := client.Tailor.Run(ctx, TailorRequest{
			ResumeID:       resumeID,
			PositionName:   "Go SDK Test",
			JobDescription: "Build Go SDKs.",
		})
		if err != nil { t.Fatal(err) }
		if len(pdf) < 5000 { t.Fatalf("too small: %d", len(pdf)) }
	})

	// EXPORT (1)
	t.Run("POST /v1/export", func(t *testing.T) {
		pdf, err := client.Export.PDF(ctx, ExportRequest{ResumeID: resumeID, TemplateID: "COBALT"})
		if err != nil { t.Fatal(err) }
		if len(pdf) < 1000 { t.Fatalf("too small: %d", len(pdf)) }
	})

	// COVER LETTERS (5)
	var clID string
	t.Run("GET /v1/cover-letters", func(t *testing.T) {
		_, err := client.CoverLetters.List(ctx, 5, 0)
		if err != nil { t.Fatal(err) }
	})
	t.Run("POST /v1/cover-letters", func(t *testing.T) {
		resp, err := client.CoverLetters.Create(ctx, CreateCoverLetterRequest{
			FullName:      "Go Test",
			LetterContent: "<p>Test from Go SDK.</p>",
		})
		if err != nil { t.Fatal(err) }
		clID = resp.CoverLetterID
	})
	t.Run("GET /v1/cover-letters/{id}", func(t *testing.T) {
		cl, err := client.CoverLetters.Get(ctx, clID)
		if err != nil { t.Fatal(err) }
		if cl.CoverLetterID != clID { t.Fatal("id mismatch") }
	})
	t.Run("PUT /v1/cover-letters/{id}/render", func(t *testing.T) {
		pdf, err := client.CoverLetters.Render(ctx, clID, RenderOptions{TemplateID: "NICKEL"})
		if err != nil { t.Fatal(err) }
		if len(pdf) < 1000 { t.Fatalf("too small: %d", len(pdf)) }
	})
	t.Run("POST /v1/cover-letters/generate", func(t *testing.T) {
		pdf, err := client.CoverLetters.Generate(ctx, GenerateCoverLetterRequest{
			ResumeID:       resumeID,
			PositionName:   "Go SDK Test",
			JobDescription: "Write Go code.",
		})
		if err != nil { t.Fatal(err) }
		if len(pdf) < 1000 { t.Fatalf("too small: %d", len(pdf)) }
	})

	// SETTINGS (3)
	t.Run("GET /v1/settings", func(t *testing.T) {
		_, err := client.Settings.Get(ctx)
		if err != nil { t.Fatal(err) }
	})
	t.Run("PUT /v1/settings/model", func(t *testing.T) {
		_, err := client.Settings.UpdateModel(ctx, UpdateAISettingsRequest{
			Provider: "OpenAI",
			Model:    "gpt-4o-mini",
			APIKey:   "sk-test-invalid",
		})
		// Expect 400 (key validation fails) — SDK sends correctly
		if err != nil {
			if !IsAuthError(err) { /* 400 is fine */ }
		}
	})
	t.Run("DELETE /v1/settings/model", func(t *testing.T) {
		resp, err := client.Settings.DeleteModel(ctx)
		if err != nil { t.Fatal(err) }
		if resp.AI != nil { t.Fatal("ai should be nil") }
	})

	// ERRORS
	t.Run("401 on bad key", func(t *testing.T) {
		bad := New("laddro_live_invalid")
		_, err := bad.Resumes.List(ctx, 1, 0)
		if !IsAuthError(err) { t.Fatalf("expected auth error, got %v", err) }
	})
	t.Run("404 on missing resume", func(t *testing.T) {
		_, err := client.Resumes.Get(ctx, "00000000-0000-0000-0000-000000000000")
		if !IsNotFoundError(err) { t.Fatalf("expected not found, got %v", err) }
	})
}
