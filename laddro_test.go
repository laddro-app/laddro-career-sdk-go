package laddro

import (
	"context"
	"testing"
)

func TestPublicEndpoints(t *testing.T) {
	client := New("")

	ctx := context.Background()

	t.Run("list templates", func(t *testing.T) {
		templates, err := client.Templates.List(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if len(templates) < 20 {
			t.Fatalf("expected 20+ templates, got %d", len(templates))
		}
		if templates[0].ID == "" {
			t.Fatal("template missing ID")
		}
		t.Logf("→ %d templates", len(templates))
	})

	t.Run("get template GRAPHITE", func(t *testing.T) {
		detail, err := client.Templates.Get(ctx, "GRAPHITE")
		if err != nil {
			t.Fatal(err)
		}
		if detail.ID != "GRAPHITE" {
			t.Fatalf("expected GRAPHITE, got %s", detail.ID)
		}
		if len(detail.AvailableColors) == 0 {
			t.Fatal("no colors")
		}
		t.Logf("→ %d colors, %d fonts", len(detail.AvailableColors), len(detail.AvailableFonts))
	})

	t.Run("list fonts", func(t *testing.T) {
		fonts, err := client.Templates.Fonts(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if len(fonts) < 20 {
			t.Fatalf("expected 20+ fonts, got %d", len(fonts))
		}
		t.Logf("→ %d fonts", len(fonts))
	})

	t.Run("list languages", func(t *testing.T) {
		languages, err := client.Templates.Languages(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if len(languages) != 14 {
			t.Fatalf("expected 14 languages, got %d", len(languages))
		}
		t.Logf("→ %d languages", len(languages))
	})

	t.Run("list models", func(t *testing.T) {
		models, err := client.Templates.Models(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if len(models) < 10 {
			t.Fatalf("expected 10+ providers, got %d", len(models))
		}
		t.Logf("→ %d providers", len(models))
	})
}
