package main

import (
	"fmt"
	"log"

	"go-be-encryption/encryption"
)

func main() {
	// Create a new encryption instance
	enc, err := encryption.NewOrgEncryption()
	if err != nil {
		log.Fatalf("Failed to create encryption: %v", err)
	}

	// The encrypted data you want to decrypt
	encryptedData := &encryption.EncryptedData{
		IV:            "Zin/emP/wx4ixLGT",
		EncryptedData: "OZH88Q9Cid1ZGFuCUkj7h1/0FKmV80e19U4oJ9qTkyGQ+59Ce3Kde6TY9Fm8+BSJzgT6TAbMoWAWyDebB4pNKQ==",
	}

	// Create a map to store the decrypted data
	// Since we don't know the exact structure, we'll use a map
	var decryptedData map[string]interface{}

	// Decrypt the data
	err = enc.Decrypt(encryptedData, &decryptedData)
	if err != nil {
		log.Fatalf("Failed to decrypt: %v", err)
	}

	// Print the decrypted data
	fmt.Printf("Decrypted data: %+v\n", decryptedData)
}
