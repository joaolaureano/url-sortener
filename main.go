package main

import (
	"awesomeProject1/operator"
	"log"
)

func main() {
	url := "https://www.mercadolivre.com.br/"

	shortUrl := operator.CreateNewShortURL(url)
	log.Printf("Current URL : %s\n", url)
	log.Printf("New URL : %s\n", shortUrl)

	recoveredUrl := operator.RecoverOriginalURL(shortUrl)

	log.Printf("Short URL : %s\n", shortUrl)
	log.Printf("Recovered URL : %s\n", recoveredUrl)
}
