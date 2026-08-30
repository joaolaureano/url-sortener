// Command urlshortener demonstrates the shortening round trip.
package main

import (
	"log/slog"
	"os"

	"urlshortener/operator"
)

func main() {
	if err := run("https://www.mercadolivre.com.br/"); err != nil {
		slog.Error("url shortener failed", "error", err)
		os.Exit(1)
	}
}

func run(originalURL string) error {
	shortURL, err := operator.CreateNewShortURL(originalURL)
	if err != nil {
		return err
	}

	recoveredURL, err := operator.RecoverOriginalURL(shortURL)
	if err != nil {
		return err
	}

	slog.Info("round trip complete",
		"original", originalURL, "short", shortURL, "recovered", recoveredURL)

	return nil
}
