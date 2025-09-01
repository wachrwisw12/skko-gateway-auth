package utils

import (
	"crypto/rsa"
	"encoding/base64"
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var (
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
)

func LoadKeys() {
	// Load Private Key
	privKeyBase64 := os.Getenv("PRIVATE_KEY")
	privKeyBytes, err := base64.StdEncoding.DecodeString(privKeyBase64)
	if err != nil {
		log.Fatalf("❌ Failed to decode PRIVATE_KEY: %v", err)
	}
	PrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privKeyBytes)
	if err != nil {
		log.Fatalf("❌ Failed to parse PRIVATE_KEY: %v", err)
	}

	// Load Public Key
	pubKeyBase64 := os.Getenv("PUBLIC_KEY")
	pubKeyBytes, err := base64.StdEncoding.DecodeString(pubKeyBase64)
	if err != nil {
		log.Fatalf("❌ Failed to decode PUBLIC_KEY: %v", err)
	}
	PublicKey, err = jwt.ParseRSAPublicKeyFromPEM(pubKeyBytes)
	if err != nil {
		log.Fatalf("❌ Failed to parse PUBLIC_KEY: %v", err)
	}

	log.Println("✅ RSA keys loaded successfully")
}
