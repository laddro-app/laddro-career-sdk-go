# laddro-career-sdk-go

Go SDK for the [Laddro Career API](https://api.laddro.com/reference).

## Install

```bash
go get github.com/laddro-app/laddro-career-sdk-go
```

## Usage

```go
package main

import (
	"context"
	"fmt"
	"os"

	laddro "github.com/laddro-app/laddro-career-sdk-go"
)

func main() {
	client := laddro.New("laddro_live_...")

	ctx := context.Background()

	// List resumes
	list, _ := client.Resumes.List(ctx, 20, 0)
	for _, r := range list.Items {
		fmt.Println(r.Title)
	}

	// Tailor a resume
	pdf, _ := client.Tailor.Run(ctx, laddro.TailorRequest{
		PositionName: "Senior Frontend Engineer",
		JobURL:       "https://jobs.example.com/senior-frontend",
	})
	os.WriteFile("tailored.pdf", pdf, 0644)

	// Stream progress
	events, _ := client.Tailor.Stream(ctx, laddro.TailorRequest{
		PositionName:   "Backend Engineer",
		JobDescription: "We are looking for...",
	})
	for event := range events {
		fmt.Printf("[%s] %s\n", event.Event, event.Data)
	}

	// Generate cover letter
	cl, _ := client.CoverLetters.Generate(ctx, laddro.GenerateCoverLetterRequest{
		PositionName: "Product Manager",
		JobURL:       "https://jobs.example.com/pm",
	})
	os.WriteFile("cover-letter.pdf", cl, 0644)

	// Browse templates (no auth needed)
	public := laddro.New("")
	templates, _ := public.Templates.List(ctx)
	for _, t := range templates {
		fmt.Printf("%s (ATS: %d)\n", t.Name, t.ATSScore)
	}
}
```

## Links

- [Laddro](https://laddro.com) — AI-powered career tools
- [API Reference](https://api.laddro.com/reference) — Interactive docs
- [Documentation](https://docs.laddro.com) — Guides and tutorials
- [GitHub](https://github.com/laddro-app) — All SDKs and tools

## License

MIT
