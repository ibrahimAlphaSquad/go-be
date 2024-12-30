package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"go-be-encryption/deformer"
	"go-be-encryption/encryption"
)

func main() {
	// Set environment variables for testing
	os.Setenv("CHAT_SECRET", "valid-secret-key-16")
	os.Setenv("PAYLOAD_NOISE_KEY", "test-encryption-key-2024")

	// Test Encryption Package
	testEncryption()
	fmt.Println("\n" + strings.Repeat("-", 80) + "\n")

	// Test Deformer Package
	testDeformer()
}

func testEncryption() {
	fmt.Println("Testing Encryption Package:")
	fmt.Println("---------------------------")

	// Create an encryption instance
	enc, err := encryption.NewOrgEncryption()
	if err != nil {
		log.Fatalf("Failed to create encryption: %v", err)
	}

	// Example data for encryption
	originalData := map[string]interface{}{
		"userId":   12345,
		"username": "john_doe",
		"settings": map[string]interface{}{
			"theme":    "dark",
			"language": "en",
		},
	}

	// Encrypt data
	encryptedData, err := enc.Encrypt(originalData)
	if err != nil {
		log.Fatalf("Failed to encrypt: %v", err)
	}
	fmt.Printf("Encrypted Data:\n%+v\n\n", encryptedData)

	// Decrypt data
	var decryptedData map[string]interface{}
	err = enc.Decrypt(encryptedData, &decryptedData)
	if err != nil {
		log.Fatalf("Failed to decrypt: %v", err)
	}
	fmt.Printf("Decrypted Data:\n%+v\n", decryptedData)
}

func testDeformer() {
	fmt.Println("Testing Deformer Package:")
	fmt.Println("-------------------------")

	// Create a PayloadNoise instance
	noise, err := deformer.NewPayloadNoise()
	if err != nil {
		log.Fatalf("Failed to create PayloadNoise: %v", err)
	}

	// Example payload for encoding
	payload := map[string]interface{}{
		"userId":   12345,
		"username": "john_doe",
		"action":   "login",
		"metadata": map[string]interface{}{
			"device":  "mobile",
			"version": "2.0.1",
		},
	}

	fmt.Printf("Original Payload:\n%+v\n\n", payload)

	// Encode payload
	encoded, err := noise.Encode(payload)
	if err != nil {
		log.Fatalf("Failed to encode payload: %v", err)
	}
	fmt.Printf("Encoded Payload:\n%+v\n\n", encoded)

	// Decode payload
	decoded, err := noise.Decode(encoded)
	if err != nil {
		log.Fatalf("Failed to decode payload: %v", err)
	}
	fmt.Printf("Decoded Payload:\n%+v\n", decoded)
}
