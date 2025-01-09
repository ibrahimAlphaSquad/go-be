// package main

// import (
// 	"fmt"
// 	"log"
// 	// "strings"

// 	"go-be-encryption/deformer"
// 	"go-be-encryption/encryption"
// )

// func main() {
// 	// Test Encryption Package
// 	// testEncryption()
// 	// fmt.Println("\n" + strings.Repeat("-", 80) + "\n")

// 	// Test Deformer Package
// 	testDeformer()
// }

// func testEncryption() {
// 	fmt.Println("Testing Encryption Package:")
// 	fmt.Println("---------------------------")

// 	// Create an encryption instance
// 	enc, err := encryption.NewOrgEncryption()
// 	if err != nil {
// 		log.Fatalf("Failed to create encryption: %v", err)
// 	}

// 	// Example data for encryption
// 	originalData := map[string]interface{}{
// 		"userId":   12345,
// 		"username": "john_doe",
// 		"settings": map[string]interface{}{
// 			"theme":    "dark",
// 			"language": "en",
// 		},
// 	}

// 	// Encrypt data
// 	encryptedData, err := enc.Encrypt(originalData)
// 	if err != nil {
// 		log.Fatalf("Failed to encrypt: %v", err)
// 	}
// 	fmt.Printf("Encrypted Data:\n%+v\n\n", encryptedData)

// 	// Decrypt data
// 	var decryptedData map[string]interface{}
// 	err = enc.Decrypt(encryptedData, &decryptedData)
// 	if err != nil {
// 		log.Fatalf("Failed to decrypt: %v", err)
// 	}
// 	fmt.Printf("Decrypted Data:\n%+v\n", decryptedData)
// }

// func testDeformer() {
// 	fmt.Println("Testing Deformer Package:")
// 	fmt.Println("-------------------------")

// 	// // Create a PayloadNoise instance
// 	// noise, err := deformer.NewPayloadNoise()
// 	// if err != nil {
// 	// 	log.Fatalf("Failed to create PayloadNoise: %v", err)
// 	// }

// 	// // Example payload for encoding
// 	// payload := map[string]interface{}{
// 	// 	"userId":   12345,
// 	// 	"username": "john_doe",
// 	// 	"action":   "login",
// 	// 	"metadata": map[string]interface{}{
// 	// 		"device":  "mobile",
// 	// 		"version": "2.0.1",
// 	// 	},
// 	// }

// 	// fmt.Printf("Original Payload:\n%+v\n\n", payload)

// 	// // Encode payload
// 	// encoded, err := noise.Encode(payload)
// 	// if err != nil {
// 	// 	log.Fatalf("Failed to encode payload: %v", err)
// 	// }
// 	// fmt.Printf("Encoded Payload:\n%+v\n\n", encoded)

// 	// // Decode payload
// 	// decoded, err := noise.Decode(encoded)
// 	// if err != nil {
// 	// 	log.Fatalf("Failed to decode payload: %v", err)
// 	// }
// 	// fmt.Printf("Decoded Payload:\n%+v\n", decoded)

// 	// Create PayloadNoise instance
// 	noise, err := deformer.NewPayloadNoise()
// 	if err != nil {
// 		log.Fatalf("Failed to create PayloadNoise: %v", err)
// 	}

// 	// Test Rain Payload
// 	rainPayload := &deformer.NoisyPayload{
// 		Version:   "1.0",
// 		Timestamp: 1734614571089,
// 		Salt:      "7286bc529aaf8e44196972cc00dd6319",
// 		IV:        "9598bd7cad5e06bf0758e5d7fdef8e05",
// 		Data: map[string]string{
// 			"tt_0787ade0":  "QNzDEHYDlAGTsdV5pzTpUg==",
// 			"tor_b1d772bb": "9EcP7Rync6MbMMHS1xF9TQ==",
// 			"tk_26563096":  "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",
// 		},
// 		Hash: "76a68f18f7d6955a1bc988ff6e21712252512a5f32db552553054ac56976f505",
// 	}

// 	fmt.Println("Decoding Rain Payload:")
// 	decodedRain, err := noise.Decode(rainPayload)
// 	if err != nil {
// 		log.Printf("Failed to decode Rain payload: %v", err)
// 	} else {
// 		fmt.Printf("%+v\n\n", decodedRain)
// 	}

// 	// Test Tip Payload
// 	tipPayload := &deformer.NoisyPayload{
// 		Version:   "1.0",
// 		Timestamp: 1734614571089,
// 		Salt:      "7286bc529aaf8e44196972cc00dd6319",
// 		IV:        "9598bd7cad5e06bf0758e5d7fdef8e05",
// 		Data: map[string]string{
// 			"tt_0787ade0":  "QNzDEHYDlAGTsdV5pzTpUg==",
// 			"tor_b1d772bb": "9EcP7Rync6MbMMHS1xF9TQ==",
// 			"tk_26563096":  "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",
// 		},
// 		Hash: "76a68f18f7d6955a1bc988ff6e21712252512a5f32db552553054ac56976f505",
// 	}

// 	fmt.Println("Decoding Rain Payload:")
// 	decodedTip, err := noise.Decode(tipPayload)
// 	if err != nil {
// 		log.Printf("Failed to decode Rain payload: %v", err)
// 	} else {
// 		fmt.Printf("%+v\n\n", decodedTip)
// 	}
// }

package main

import (
	"fmt"
	"log"
	"os"

	"go-be-encryption/deformer"
)

func main() {
	// First unset any existing value
	os.Unsetenv("PAYLOAD_NOISE_KEY")

	// Set the exact same key as used in JavaScript
	os.Setenv("PAYLOAD_NOISE_KEY", "AVwade_3e5hZa9kVwsSt_p6yTmEYgUDnL4iAcX9QlPWvrN0oG7jFbC2Hwp6yTmEYgUDnL4iAcX9Ql")

	// Create PayloadNoise instance
	noise, err := deformer.NewPayloadNoise()
	if err != nil {
		log.Fatalf("Failed to create PayloadNoise: %v", err)
	}

	// Create payload matching JavaScript output
	jsPayload := &deformer.NoisyPayload{
		Version:   "1.0",
		Timestamp: 1735570785873,
		Salt:      "42a0ecdcb3f939df05d8c8bedcff9b76",
		IV:        "5f9cb6c8f569c551f51286abfb07c7df",
		Data: map[string]string{
			"username_06e3eea5":  "RUOq8lAbyRa0PYYVeGbqMQ==",
			"action_a51c5a86":    "eYsjc0NpK3bgvqxAraJF7Q==",
			"timestamp_e2692e63": "h6lpZQVGY3r3PwATk8DSUA==",
		},
		Hash: "5e28e143f80a38f774c498ea653924322bb1e07cfb3053715279aadb8995f4bf",
	}

	// Attempt to decode
	decoded, err := noise.Decode(jsPayload)
	if err != nil {
		log.Printf("Failed to decode payload: %v", err)
	} else {
		fmt.Printf("Decoded JS payload: %+v\n", decoded)
	}

	// Let's also try encoding some data
	testPayload := map[string]interface{}{
		"username":  "john_doe",
		"action":    "login",
		"timestamp": 1234567890,
	}

	encoded, err := noise.Encode(testPayload)
	if err != nil {
		log.Printf("Failed to encode payload: %v", err)
	} else {
		fmt.Printf("\nEncoded test payload: %+v\n", encoded)
	}
}
