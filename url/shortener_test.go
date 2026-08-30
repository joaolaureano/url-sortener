package url

import (
	"errors"
	"math"
	"strings"
	"testing"
)

func TestCodeSpaceMatchesAlphabet(t *testing.T) {
	t.Parallel()

	want := uint64(math.Pow(float64(len(alphabet)), CodeLength))
	if codeSpace != want {
		t.Fatalf("codeSpace = %d, want %d^%d = %d", uint64(codeSpace), len(alphabet), CodeLength, want)
	}
}

func TestShortCodeIsDeterministicAndWellFormed(t *testing.T) {
	t.Parallel()

	const originalURL = "https://www.mercadolivre.com.br/"

	got, err := ShortCode(originalURL, 0)
	if err != nil {
		t.Fatalf("ShortCode() error = %v, want nil", err)
	}

	again, _ := ShortCode(originalURL, 0)
	if got != again {
		t.Fatalf("ShortCode() = %q then %q, want a stable code", got, again)
	}

	if len(got) != CodeLength {
		t.Fatalf("len(ShortCode()) = %d, want %d", len(got), CodeLength)
	}

	for _, r := range got {
		if !strings.ContainsRune(alphabet, r) {
			t.Fatalf("ShortCode() = %q contains %q, outside the base62 alphabet", got, r)
		}
	}
}

func TestShortCodeAttemptsDiffer(t *testing.T) {
	t.Parallel()

	canonical, _ := ShortCode("https://example.com", 0)

	alternate, err := ShortCode("https://example.com", 1)
	if err != nil {
		t.Fatalf("ShortCode() error = %v, want nil", err)
	}

	if canonical == alternate {
		t.Fatalf("ShortCode(attempt=1) = %q, want a code different from %q", alternate, canonical)
	}
}

func TestShortCodeRejectsInvalidURL(t *testing.T) {
	t.Parallel()

	if _, err := ShortCode("not a url", 0); !errors.Is(err, ErrInvalidURL) {
		t.Fatalf("ShortCode() error = %v, want ErrInvalidURL", err)
	}
}

func TestBuildURLExtractCodeRoundTrip(t *testing.T) {
	t.Parallel()

	const code = "aB3xY9z"

	got, err := ExtractCode(BuildURL(code))
	if err != nil {
		t.Fatalf("ExtractCode() error = %v, want nil", err)
	}

	if got != code {
		t.Fatalf("ExtractCode() = %q, want %q", got, code)
	}
}

func TestExtractCodeRejectsForeignURL(t *testing.T) {
	t.Parallel()

	tests := map[string]string{
		"other service": "https://bit.ly/abc",
		"empty code":    Domain,
		"bare code":     "aB3xY9z",
	}

	for name, shortURL := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if _, err := ExtractCode(shortURL); !errors.Is(err, ErrInvalidShortURL) {
				t.Fatalf("ExtractCode(%q) error = %v, want ErrInvalidShortURL", shortURL, err)
			}
		})
	}
}
