package operator

import (
	"errors"
	"testing"

	"urlshortener/storage"
	"urlshortener/url"
)

// isolate points the package-level store at a fresh instance for one test and
// restores the previous one afterwards.
func isolate(t *testing.T) {
	t.Helper()

	previous := storage.Default
	storage.Default = storage.NewMemoryStore()

	t.Cleanup(func() { storage.Default = previous })
}

func TestCreateThenRecoverRoundTrip(t *testing.T) {
	isolate(t)

	const originalURL = "https://www.mercadolivre.com.br/"

	shortURL, err := CreateNewShortURL(originalURL)
	if err != nil {
		t.Fatalf("CreateNewShortURL() error = %v, want nil", err)
	}

	got, err := RecoverOriginalURL(shortURL)
	if err != nil {
		t.Fatalf("RecoverOriginalURL() error = %v, want nil", err)
	}

	if got != originalURL {
		t.Fatalf("RecoverOriginalURL() = %q, want %q", got, originalURL)
	}
}

func TestCreateNewShortURLIsIdempotent(t *testing.T) {
	isolate(t)

	const originalURL = "https://www.example.org/stable"

	first, err := CreateNewShortURL(originalURL)
	if err != nil {
		t.Fatalf("CreateNewShortURL() error = %v, want nil", err)
	}

	second, err := CreateNewShortURL(originalURL)
	if err != nil {
		t.Fatalf("CreateNewShortURL() error = %v, want nil", err)
	}

	if first != second {
		t.Fatalf("CreateNewShortURL() = %q then %q, want the same short url", first, second)
	}
}

func TestCreateNewShortURLRejectsInvalidURL(t *testing.T) {
	isolate(t)

	if _, err := CreateNewShortURL("nope"); !errors.Is(err, url.ErrInvalidURL) {
		t.Fatalf("CreateNewShortURL() error = %v, want url.ErrInvalidURL", err)
	}
}

func TestCreateNewShortURLSurvivesCollision(t *testing.T) {
	isolate(t)

	const originalURL = "https://www.example.org/collides"

	// Squat on the canonical code with a different URL.
	canonical, err := url.ShortCode(originalURL, 0)
	if err != nil {
		t.Fatalf("ShortCode() error = %v, want nil", err)
	}

	storage.CreateReference(canonical, "https://squatter.example")

	shortURL, err := CreateNewShortURL(originalURL)
	if err != nil {
		t.Fatalf("CreateNewShortURL() error = %v, want nil", err)
	}

	if shortURL == url.BuildURL(canonical) {
		t.Fatal("CreateNewShortURL() reused the squatted code, want an alternate")
	}

	got, err := RecoverOriginalURL(shortURL)
	if err != nil {
		t.Fatalf("RecoverOriginalURL() error = %v, want nil", err)
	}

	if got != originalURL {
		t.Fatalf("RecoverOriginalURL() = %q, want %q", got, originalURL)
	}

	squatted, _ := storage.RecoverReference(canonical)
	if want := "https://squatter.example"; squatted != want {
		t.Fatalf("squatted code now holds %q, want %q untouched", squatted, want)
	}
}

func TestRecoverOriginalURLUnknownCode(t *testing.T) {
	isolate(t)

	if _, err := RecoverOriginalURL(url.BuildURL("0000000")); !errors.Is(err, storage.ErrNotFound) {
		t.Fatalf("RecoverOriginalURL() error = %v, want storage.ErrNotFound", err)
	}
}

func TestRecoverOriginalURLForeignLink(t *testing.T) {
	isolate(t)

	if _, err := RecoverOriginalURL("https://bit.ly/abc"); !errors.Is(err, url.ErrInvalidShortURL) {
		t.Fatalf("RecoverOriginalURL() error = %v, want url.ErrInvalidShortURL", err)
	}
}

func TestDistinctURLsGetDistinctShortURLs(t *testing.T) {
	isolate(t)

	seen := make(map[string]string)

	for _, originalURL := range []string{
		"https://example.com/a", "https://example.com/b", "https://example.com/c",
		"https://example.com/d", "https://example.com/e",
	} {
		shortURL, err := CreateNewShortURL(originalURL)
		if err != nil {
			t.Fatalf("CreateNewShortURL(%q) error = %v, want nil", originalURL, err)
		}

		if previous, ok := seen[shortURL]; ok {
			t.Fatalf("%q and %q share the short url %q", previous, originalURL, shortURL)
		}

		seen[shortURL] = originalURL
	}
}
