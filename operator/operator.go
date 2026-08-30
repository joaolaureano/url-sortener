// Package operator is the service facade: it wires URL shortening to the
// reference store and is the only entry point callers need.
package operator

import (
	"errors"
	"fmt"
	"log/slog"

	"urlshortener/storage"
	"urlshortener/url"
)

// maxCollisionAttempts caps the probing for a free short code. Each attempt
// costs one hash and one map write, and reaching the cap is astronomically
// unlikely at the current code space, so a small bound is enough.
const maxCollisionAttempts = 10

// ErrCollisionExhausted reports that every candidate code for a URL was taken.
// It signals a saturated code space, not a transient failure.
var ErrCollisionExhausted = errors.New("operator: could not allocate a free short code")

// CreateNewShortURL shortens originalURL and stores the reference, returning
// the public short URL.
//
// The operation is idempotent: shortening the same URL again returns the same
// short URL. If a code is already taken by a different URL, alternates are
// tried before giving up with ErrCollisionExhausted. Invalid input yields
// url.ErrInvalidURL.
func CreateNewShortURL(originalURL string) (string, error) {
	for attempt := range maxCollisionAttempts {
		code, err := url.ShortCode(originalURL, attempt)
		if err != nil {
			return "", err
		}

		created, current := storage.CreateReference(code, originalURL)
		if !created {
			slog.Warn("short code collision, trying next candidate",
				"code", code, "held_by", current, "url", originalURL, "attempt", attempt)

			continue
		}

		slog.Info("shortened url", "code", code, "url", originalURL)

		return url.BuildURL(code), nil
	}

	return "", fmt.Errorf("%w: %q after %d attempts", ErrCollisionExhausted, originalURL, maxCollisionAttempts)
}

// RecoverOriginalURL returns the original URL behind shortURL. It yields
// url.ErrInvalidShortURL for a foreign link and storage.ErrNotFound for an
// unknown code.
func RecoverOriginalURL(shortURL string) (string, error) {
	code, err := url.ExtractCode(shortURL)
	if err != nil {
		return "", err
	}

	return storage.RecoverReference(code)
}
