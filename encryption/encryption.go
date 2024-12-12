package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// EncryptedData represents the structure of encrypted data
type EncryptedData struct {
	IV            string `json:"iv"`
	EncryptedData string `json:"encryptedData"`
}

// OrgEncryption handles encryption/decryption operations
type OrgEncryption struct {
	orgID string
}

// init loads the environment variables from .env file
func init() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// We don't return an error here as the .env file is optional
		// Environment variables might be set through other means
		fmt.Printf("Warning: No .env file found: %v\n", err)
	}
}

// NewOrgEncryption creates a new instance of OrgEncryption
func NewOrgEncryption() (*OrgEncryption, error) {
	orgID := os.Getenv("CHAT_SECRET")
	if orgID == "" {
		return nil, fmt.Errorf("CHAT_SECRET environment variable is not set")
	}

	if len(orgID) < 16 {
		return nil, fmt.Errorf("CHAT_SECRET must be at least 16 characters long for security")
	}

	return &OrgEncryption{
		orgID: orgID,
	}, nil
}

// generateKey creates an encryption key from orgID
func (o *OrgEncryption) generateKey() []byte {
	// Use SHA-256 to get a consistent key length
	hash := sha256.Sum256([]byte(o.orgID))
	return hash[:]
}

// Encrypt encrypts data using AES-GCM
func (o *OrgEncryption) Encrypt(data interface{}) (*EncryptedData, error) {
	// Convert data to JSON bytes
	plaintext, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	// Generate key
	key := o.generateKey()

	// Create cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate nonce (IV)
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt data
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	// Create response
	encryptedData := &EncryptedData{
		IV:            base64.StdEncoding.EncodeToString(nonce),
		EncryptedData: base64.StdEncoding.EncodeToString(ciphertext),
	}

	return encryptedData, nil
}

// Decrypt decrypts data using AES-GCM
func (o *OrgEncryption) Decrypt(encryptedData *EncryptedData, target interface{}) error {
	// Decode base64 strings
	nonce, err := base64.StdEncoding.DecodeString(encryptedData.IV)
	if err != nil {
		return fmt.Errorf("failed to decode IV: %w", err)
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData.EncryptedData)
	if err != nil {
		return fmt.Errorf("failed to decode encrypted data: %w", err)
	}

	// Generate key
	key := o.generateKey()

	// Create cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to create GCM: %w", err)
	}

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("failed to decrypt: %w", err)
	}

	// Unmarshal into target
	if err := json.Unmarshal(plaintext, target); err != nil {
		return fmt.Errorf("failed to unmarshal decrypted data: %w", err)
	}

	return nil
}
