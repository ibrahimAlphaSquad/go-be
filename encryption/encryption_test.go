package encryption

import (
	"os"
	"testing"
)

type TestData struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

func TestNewOrgEncryption(t *testing.T) {
	// Test with no environment variable
	os.Unsetenv("CHAT_SECRET")
	_, err := NewOrgEncryption()
	if err == nil {
		t.Error("Expected error when CHAT_SECRET is not set")
	}

	// Test with short key
	os.Setenv("CHAT_SECRET", "short")
	_, err = NewOrgEncryption()
	if err == nil {
		t.Error("Expected error when CHAT_SECRET is too short")
	}

	// Test with valid key
	os.Setenv("CHAT_SECRET", "valid-secret-key-16")
	enc, err := NewOrgEncryption()
	if err != nil {
		t.Errorf("Unexpected error with valid key: %v", err)
	}
	if enc == nil {
		t.Error("Expected valid OrgEncryption instance")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	// Setup
	os.Setenv("CHAT_SECRET", "valid-secret-key-16")
	enc, err := NewOrgEncryption()
	if err != nil {
		t.Fatalf("Failed to create encryption instance: %v", err)
	}

	// Test data
	originalData := TestData{
		Name:    "Test Name",
		Message: "Hello, World!",
	}

	// Encrypt
	encrypted, err := enc.Encrypt(originalData)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}
	if encrypted.IV == "" || encrypted.EncryptedData == "" {
		t.Error("Expected non-empty IV and EncryptedData")
	}

	// Decrypt
	var decryptedData TestData
	err = enc.Decrypt(encrypted, &decryptedData)
	if err != nil {
		t.Fatalf("Failed to decrypt: %v", err)
	}

	// Compare
	if decryptedData.Name != originalData.Name || decryptedData.Message != originalData.Message {
		t.Error("Decrypted data doesn't match original")
	}
}

func TestDecryptInvalidData(t *testing.T) {
	// Setup
	os.Setenv("CHAT_SECRET", "valid-secret-key-16")
	enc, err := NewOrgEncryption()
	if err != nil {
		t.Fatalf("Failed to create encryption instance: %v", err)
	}

	// Test invalid base64
	invalidData := &EncryptedData{
		IV:            "invalid-base64",
		EncryptedData: "invalid-base64",
	}
	var target TestData
	err = enc.Decrypt(invalidData, &target)
	if err == nil {
		t.Error("Expected error when decrypting invalid base64")
	}
}
