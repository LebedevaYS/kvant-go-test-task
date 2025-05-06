package config

import (
	"crypto/rand"
	"log"
)

var SecretKey []byte

func Init() {
	key := make([]byte, 32) // 32 байта для HS256
	if _, err := rand.Read(key); err != nil {
		log.Fatal("Failed to generate secret key")
	}
	SecretKey = key
}
