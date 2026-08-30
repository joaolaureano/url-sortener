// Package url turns original URLs into short codes and assembles the public
// short URLs that carry them. It owns the link format; no other layer should
// build or parse a short URL by hand.
package url

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strings"

	"github.com/asaskevich/govalidator"

	"urlshortener/hash"
)

const (
	// Domain is the prefix of every short URL this service issues.
	Domain = "https://me.li/"

	// CodeLength is the number of base62 characters in a short code.
	CodeLength = 7

	// alphabet is the base62 digit set, ordered so codes sort lexicographically.
	alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	// codeSpace is len(alphabet)^CodeLength, the number of distinct codes.
	// At 62^7 the first expected collision sits around 2.6M links, versus
	// ~1.2K for a 5-character hex code. TestCodeSpaceMatchesAlphabet keeps
	// this literal honest.
	codeSpace = 3_521_614_606_208
)

var (
	// ErrInvalidURL reports that the input is not a well-formed URL.
	ErrInvalidURL = errors.New("url: invalid url")

	// ErrInvalidShortURL reports that a string was not issued by this service.
	ErrInvalidShortURL = errors.New("url: not a short url of this service")
)

// ShortCode returns the short code for originalURL.
//
// attempt selects among deterministic candidates for the same URL: 0 yields
// the canonical code, and higher values yield alternates that callers can fall
// back to when a code is already taken. The result is stable across processes,
// so the same (URL, attempt) pair always maps to the same code.
//
// It returns ErrInvalidURL if originalURL is not a valid URL.
func ShortCode(originalURL string, attempt int) (string, error) {
	if !govalidator.IsURL(originalURL) {
		return "", fmt.Errorf("%w: %q", ErrInvalidURL, originalURL)
	}

	seed := originalURL
	if attempt > 0 {
		seed = fmt.Sprintf("%s#%d", originalURL, attempt)
	}

	return base62(binary.BigEndian.Uint64(hash.Digest(seed)[:8])), nil
}

// BuildURL assembles the public short URL carrying code.
func BuildURL(code string) string {
	return Domain + code
}

// ExtractCode is the inverse of BuildURL. It returns ErrInvalidShortURL if
// shortURL was not issued by this service or carries no code.
func ExtractCode(shortURL string) (string, error) {
	code, found := strings.CutPrefix(shortURL, Domain)
	if !found || code == "" {
		return "", fmt.Errorf("%w: %q", ErrInvalidShortURL, shortURL)
	}

	return code, nil
}

// base62 renders n as a zero-padded, fixed-width base62 code.
func base62(n uint64) string {
	n %= codeSpace

	buf := make([]byte, CodeLength)
	for i := CodeLength - 1; i >= 0; i-- {
		buf[i] = alphabet[n%uint64(len(alphabet))]
		n /= uint64(len(alphabet))
	}

	return string(buf)
}
