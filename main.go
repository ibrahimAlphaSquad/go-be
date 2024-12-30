package main

import (
	"fmt"
	"log"

	"go-be-encryption/deformer"
	"go-be-encryption/encryption"
)

func main() {
	// Create an encryption instance
	enc, err := encryption.NewOrgEncryption()
	if err != nil {
		log.Fatalf("Failed to create encryption: %v", err)
	}

	// Example: Decrypting data
	encryptedData := &encryption.EncryptedData{
		IV:            "ZTw8qcuUVeUI382i",
		EncryptedData: "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",
	}

	var decryptedData map[string]interface{}
	err = enc.Decrypt(encryptedData, &decryptedData)
	if err != nil {
		log.Fatalf("Failed to decrypt: %v", err)
	}
	fmt.Printf("\nDecrypted data: %+v\n", decryptedData)

	// Create a PayloadNoise instance
	noise, err := deformer.NewPayloadNoise()
	if err != nil {
		log.Fatalf("\n\n\n\n\n\n\n\nFailed to create PayloadNoise: %v", err)
	}

	// Example: Encoding data
	payload := map[string]interface{}{
		"userId":   12345,
		"username": "john_doe",
	}

	fmt.Printf("\n\n\n\nPayload to be encrypted: %+v\n\n", payload)

	encoded, err := noise.Encode(payload)
	if err != nil {
		log.Fatalf("Failed to encode payload: %v", err)
	}
	fmt.Printf("\n\n\n\nEncoded payload: %+v\n\n", encoded)

	// Example: Decoding hardcoded data
	hardcoded := &deformer.NoisyPayload{
		Version:   "1.0",
		Timestamp: 1734604976778,
		Salt:      "2cf0f9685e242a9a1c0a35992ce596d0",
		IV:        "effe536170cb2d4d1933d146433ed558",
		Data: map[string]string{
			"tt_b960d363": "QD2izyJSxW0c7zMJwMsS5A==",
			"ct_012722a9": "QD2izyJSxW0c7zMJwMsS5A==",
			"tk_fa53fab9": "a0yyiqfoW7Klrcv5Dx+xdCVaUCONMpHvUkqvC4PYtT+WD6gcweVc3lV5VEZ8aHAdKS2iimtilE7EWQTphgwiIAXmd8dXalxfC2DK/I7SfFGfnjJoKefbZyVZ9gWVs8K1...",
		},
		Hash: "200861202e69701ca63ba8c065afd36a7d9a5ad6b348bd0da93483ca3a4e06e7",
	}

	decoded, err := noise.Decode(hardcoded)
	if err != nil {
		log.Fatalf("\n\n\n\nFailed to decode hardcoded payload: %v", err)
	}
	fmt.Printf("\n\n\n\nDecoded hardcoded data: %+v\n", decoded)

	// Test round-trip encryption/decryption
	roundTripDecoded, err := noise.Decode(encoded)
	if err != nil {
		log.Fatalf("Failed to decode round-trip payload: %v", err)
	}
	fmt.Printf("\nRound-trip decoded data: %+v\n", roundTripDecoded)
}
