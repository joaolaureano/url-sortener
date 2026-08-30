package url

import (
	hash2 "awesomeProject1/hash"
	"fmt"
	"github.com/asaskevich/govalidator"
	"log"
)

const baseUrl = "https://me.li/%s"

func GenerateShortURL(originalUrl string) string {
	validURL := govalidator.IsURL(originalUrl)
	if !validURL {
		log.Panicf("URL %s is not a valid url.", originalUrl)
	}
	log.Printf("URL %s is a valid url.", originalUrl)
	shortHash := enshortOriginalURL(originalUrl)
	shortURL := fmt.Sprintf(baseUrl, shortHash)
	log.Printf("New enshorted URL is %s.", shortURL)
	return shortURL
}

func enshortOriginalURL(originalURL string) string {
	hash := hash2.Encode(originalURL)
	return hash[:5]
}
