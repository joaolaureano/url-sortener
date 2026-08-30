package operator

import (
	database "awesomeProject1/storage"
	urlEnshorter "awesomeProject1/url"
)

func RecoverOriginalURL(url string) string {
	originalURL := database.RecoverReference(url)
	return originalURL
}
func CreateNewShortURL(url string) string {
	shortUrl := urlEnshorter.GenerateShortURL(url)
	database.CreateReference(shortUrl, url)
	return shortUrl
}
