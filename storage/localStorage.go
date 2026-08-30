package storage

import "log"

var table = make(map[string]string)

func CreateReference(key string, value string) bool {
	table[key] = value
	log.Printf("Created Reference. Key is %s. Value is %s", key, value)
	return true
}

func RecoverReference(key string) string {
	value := table[key]
	log.Printf("Recovered Reference. Key is %s. Value is %s", key, value)
	return value
}
