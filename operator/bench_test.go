package operator

import (
	"io"
	"log/slog"
	"testing"

	"urlshortener/storage"
)

// Benchmarks for the shortening facade. They reuse one URL per benchmark so the
// store stays a fixed size across iterations: CreateNewShortURL is idempotent
// for a repeated URL, so the map neither grows nor triggers collision probing.

const benchURL = "https://www.mercadolivre.com.br/anuncio/item-12345"

func isolateBench(b *testing.B) {
	b.Helper()

	previousStore := storage.Default
	storage.Default = storage.NewMemoryStore()

	// Silence the operator's per-call logging: it writes to stderr on every
	// iteration and would dominate both the timing and the CPU profile.
	previousLogger := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))

	b.Cleanup(func() {
		storage.Default = previousStore
		slog.SetDefault(previousLogger)
	})
}

func BenchmarkCreateNewShortURL(b *testing.B) {
	isolateBench(b)

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		if _, err := CreateNewShortURL(benchURL); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRecoverOriginalURL(b *testing.B) {
	isolateBench(b)

	shortURL, err := CreateNewShortURL(benchURL)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		if _, err := RecoverOriginalURL(shortURL); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRoundTrip(b *testing.B) {
	isolateBench(b)

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		shortURL, err := CreateNewShortURL(benchURL)
		if err != nil {
			b.Fatal(err)
		}
		if _, err := RecoverOriginalURL(shortURL); err != nil {
			b.Fatal(err)
		}
	}
}
