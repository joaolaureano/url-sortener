package hash

import (
	"crypto/md5"
	"encoding/hex"
	"log"
)

type HashOperator interface {
	Encode(content string) string
	Compare(content string, hash string) bool
}

func Encode(content string) string {
	hash := md5.Sum([]byte(content))
	hashString := hex.EncodeToString(hash[:])
	log.Print("Hash for %s is %s", content, hash)
	return hashString
}
