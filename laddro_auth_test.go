package laddro

import (
	"context"
	"os"
	"testing"
)

func TestProtectedEndpoints(t *testing.T) {
	apiKey := os.Getenv("LADDRO_API_KEY")
	if apiKey == "" {
		t.Skip("LADDRO_API_KEY not set")
	}

	client := New(apiKey)
	ctx := context.Background()

	t.Run("list resumes", func(t *testing.T) {
		list, err := client.Resumes.List(ctx, 5, 0)
		if err != nil {
			t.Fatal(err)
		}
		if list.Items == nil {
			t.Fatal("items is nil")
		}
		t.Logf("→ %d resumes (total: %d)", len(list.Items), list.Total)
	})

	t.Run("get settings", func(t *testing.T) {
		settings, err := client.Settings.Get(ctx)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("→ ai configured: %v", settings.AI != nil)
	})

	t.Run("list cover letters", func(t *testing.T) {
		list, err := client.CoverLetters.List(ctx, 5, 0)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("→ %d cover letters", len(list.Items))
	})

	t.Run("get resume by id", func(t *testing.T) {
		list, err := client.Resumes.List(ctx, 1, 0)
		if err != nil {
			t.Fatal(err)
		}
		if len(list.Items) == 0 {
			t.Skip("no resumes")
		}
		resume, err := client.Resumes.Get(ctx, list.Items[0].ResumeID)
		if err != nil {
			t.Fatal(err)
		}
		if resume.ResumeID != list.Items[0].ResumeID {
			t.Fatalf("id mismatch: %s != %s", resume.ResumeID, list.Items[0].ResumeID)
		}
		t.Logf("→ %s", resume.Title)
	})

	t.Run("auth error on bad key", func(t *testing.T) {
		bad := New("laddro_live_invalid")
		_, err := bad.Resumes.List(ctx, 1, 0)
		if err == nil {
			t.Fatal("expected error")
		}
		if !IsAuthError(err) {
			t.Fatalf("expected auth error, got: %v", err)
		}
	})
}
