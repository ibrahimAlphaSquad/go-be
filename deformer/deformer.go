package deformer

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/pbkdf2"
)

// PayloadNoise handles encryption and decryption of payload data
type PayloadNoise struct {
	key []byte
}

// NoisyPayload represents the encrypted payload structure
type NoisyPayload struct {
	Version   string            `json:"_v"`
	Timestamp int64             `json:"_t"`
	Salt      string            `json:"_s"`
	IV        string            `json:"_i"`
	Hash      string            `json:"_h"`
	Data      map[string]string `json:"data"`
}

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: No .env file found: %v\n", err)
	}
}

// NewPayloadNoise creates a new PayloadNoise instance
func NewPayloadNoise() (*PayloadNoise, error) {
	key := os.Getenv("PAYLOAD_NOISE_KEY")
	if key == "" {
		key = "test-encryption-key-2024" // Default key
	}

	fmt.Printf("Using key (string): %s\n", key)
	// Store key as UTF-8 bytes to match CryptoJS behavior
	keyBytes := []byte(key)
	fmt.Printf("Using key (hex): %x\n", keyBytes)

	if len(key) < 16 {
		return nil, fmt.Errorf("PAYLOAD_NOISE_KEY must be at least 16 characters long")
	}

	return &PayloadNoise{
		key: keyBytes,
	}, nil
}

// Encode encrypts a payload map into a NoisyPayload
func (p *PayloadNoise) Encode(payload map[string]interface{}) (*NoisyPayload, error) {
	if payload == nil {
		return nil, fmt.Errorf("payload must not be nil")
	}

	// Generate salt and IV
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %v", err)
	}

	iv := make([]byte, 16)
	if _, err := rand.Read(iv); err != nil {
		return nil, fmt.Errorf("failed to generate IV: %v", err)
	}

	// Derive key using PBKDF2
	derivedKey := pbkdf2.Key(p.key, salt, 1000, 32, sha256.New)

	// Create cipher block
	block, err := aes.NewCipher(derivedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher block: %v", err)
	}

	// Initialize CBC mode
	mode := cipher.NewCBCEncrypter(block, iv)

	// Process each field
	noisyData := make(map[string]string)
	timestamp := time.Now().UnixMilli()

	for key, value := range payload {
		// Convert value to JSON string
		valueBytes, err := json.Marshal(value)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal value for key %s: %v", key, err)
		}

		// Pad data
		paddedData := p.pkcs7Pad(valueBytes, aes.BlockSize)

		// Encrypt value
		encrypted := make([]byte, len(paddedData))
		mode.CryptBlocks(encrypted, paddedData)

		// Generate noisy key
		noisyKey := p.generateNoisyKey(key, timestamp)

		// Store encrypted value as base64 string
		noisyData[noisyKey] = base64.StdEncoding.EncodeToString(encrypted)
	}

	// Create final payload
	finalPayload := &NoisyPayload{
		Version:   "1.0",
		Timestamp: timestamp,
		Salt:      hex.EncodeToString(salt),
		IV:        hex.EncodeToString(iv),
		Data:      noisyData,
	}

	// Add hash last
	finalPayload.Hash = p.generateHash(finalPayload.Data)

	return finalPayload, nil
}

// Decode decrypts a NoisyPayload back into the original payload
func (p *PayloadNoise) Decode(noisyPayload *NoisyPayload) (map[string]interface{}, error) {
	if !p.validatePayload(noisyPayload) {
		return nil, fmt.Errorf("invalid payload structure")
	}

	// Verify hash
	if !p.verifyHash(noisyPayload) {
		return nil, fmt.Errorf("payload integrity check failed")
	}

	// Decode salt and IV
	salt, err := hex.DecodeString(noisyPayload.Salt)
	if err != nil {
		return nil, fmt.Errorf("failed to decode salt: %v", err)
	}

	iv, err := hex.DecodeString(noisyPayload.IV)
	if err != nil {
		return nil, fmt.Errorf("failed to decode IV: %v", err)
	}

	// Derive key using PBKDF2
	derivedKey := pbkdf2.Key(p.key, salt, 1000, 32, sha256.New)

	// Create cipher block
	block, err := aes.NewCipher(derivedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher block: %v", err)
	}

	// Initialize CBC mode
	mode := cipher.NewCBCDecrypter(block, iv)

	originalPayload := make(map[string]interface{})

	// Decrypt each field
	for noisyKey, encryptedValue := range noisyPayload.Data {
		// Decode base64
		encrypted, err := base64.StdEncoding.DecodeString(encryptedValue)
		if err != nil {
			originalPayload[p.recoverOriginalKey(noisyKey)] = nil
			continue
		}

		// Decrypt
		decrypted := make([]byte, len(encrypted))
		mode.CryptBlocks(decrypted, encrypted)

		// Unpad
		unpaddedData, err := p.pkcs7Unpad(decrypted)
		if err != nil {
			originalPayload[p.recoverOriginalKey(noisyKey)] = nil
			continue
		}

		// Parse JSON
		var value interface{}
		if err := json.Unmarshal(unpaddedData, &value); err != nil {
			originalPayload[p.recoverOriginalKey(noisyKey)] = nil
			continue
		}

		originalPayload[p.recoverOriginalKey(noisyKey)] = value
	}

	return originalPayload, nil
}

// Helper methods

func (p *PayloadNoise) generateNoisyKey(key string, timestamp int64) string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%s%d", key, timestamp)))
	hash := hex.EncodeToString(h.Sum(nil))
	return fmt.Sprintf("%s_%s", key, hash[:8])
}

func (p *PayloadNoise) recoverOriginalKey(noisyKey string) string {
	parts := strings.Split(noisyKey, "_")
	if len(parts) < 1 {
		return noisyKey
	}
	return parts[0]
}

func (p *PayloadNoise) generateHash(payload map[string]string) string {
    // Convert the payload to JSON in the same way JavaScript does
    payloadBytes, _ := json.Marshal(payload)
    
    // Debug: Print the JSON string before hashing
    fmt.Printf("JSON string before hashing: %s\n", string(payloadBytes))
    
    // Create HMAC using SHA256
    h := hmac.New(sha256.New, p.key)
    h.Write(payloadBytes)
    
    // Get the hash result
    hashResult := hex.EncodeToString(h.Sum(nil))
    fmt.Printf("Generated hash: %s\n", hashResult)
    
    return hashResult
}

func (p *PayloadNoise) verifyHash(payload *NoisyPayload) bool {
    providedHash := payload.Hash
    calculatedHash := p.generateHash(payload.Data)
    
    fmt.Printf("Provided hash:   %s\n", providedHash)
    fmt.Printf("Calculated hash: %s\n", calculatedHash)
    
    return providedHash == calculatedHash
}

func (p *PayloadNoise) validatePayload(payload *NoisyPayload) bool {
	return payload != nil &&
		payload.Version != "" &&
		payload.Timestamp != 0 &&
		payload.Salt != "" &&
		payload.IV != "" &&
		payload.Hash != "" &&
		payload.Data != nil
}

func (p *PayloadNoise) pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func (p *PayloadNoise) pkcs7Unpad(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, fmt.Errorf("invalid padding: data is empty")
	}

	padding := int(data[length-1])
	if padding > length {
		return nil, fmt.Errorf("invalid padding size")
	}

	for i := length - padding; i < length; i++ {
		if data[i] != byte(padding) {
			return nil, fmt.Errorf("invalid padding values")
		}
	}

	return data[:length-padding], nil
}
