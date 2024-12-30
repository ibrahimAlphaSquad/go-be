package deformer

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/joho/godotenv"
)

func setupTestEnvironment(t *testing.T) {
    // First try to load from .env.test if it exists
    if err := godotenv.Load(".env.test"); err != nil {
        // If no test-specific .env exists, set the key directly
        err = os.Setenv("PAYLOAD_NOISE_KEY", "test-key-12345678")
        if err != nil {
            t.Fatalf("Failed to set environment variable: %v", err)
        }
    }
}

func cleanupTestEnvironment(t *testing.T) {
    // Only unset if we're using the direct environment variable
    if os.Getenv("PAYLOAD_NOISE_KEY") == "test-key-12345678" {
        err := os.Unsetenv("PAYLOAD_NOISE_KEY")
        if err != nil {
            t.Fatalf("Failed to unset environment variable: %v", err)
        }
    }
}

func TestNewPayloadNoise(t *testing.T) {
	tests := []struct {
		name      string
		envKey    string
		expectErr bool
	}{
		{
			name:      "valid key",
			envKey:    "test-key-12345678",
			expectErr: false,
		},
		{
			name:      "empty key",
			envKey:    "",
			expectErr: true,
		},
		{
			name:      "short key",
			envKey:    "short",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := os.Setenv("PAYLOAD_NOISE_KEY", tt.envKey); err != nil {
				t.Fatalf("Failed to set environment variable: %v", err)
			}
			defer os.Unsetenv("PAYLOAD_NOISE_KEY")

			noise, err := NewPayloadNoise()
			if tt.expectErr {
				if err == nil {
					t.Error("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if noise == nil {
					t.Error("Expected PayloadNoise instance but got nil")
				}
			}
		})
	}
}

func TestPayloadNoise_EncodeAndDecode(t *testing.T) {
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)

	noise, err := NewPayloadNoise()
	if err != nil {
		t.Fatalf("Failed to create PayloadNoise: %v", err)
	}

	testCases := []struct {
		name    string
		payload map[string]interface{}
	}{
		{
			name: "basic types",
			payload: map[string]interface{}{
				"string": "hello",
				"number": 42,
				"bool":   true,
				"null":   nil,
			},
		},
		{
			name: "nested objects",
			payload: map[string]interface{}{
				"nested": map[string]interface{}{
					"foo": "bar",
					"num": 123,
				},
				"array": []interface{}{1, 2, 3},
			},
		},
		{
			name: "complex structure",
			payload: map[string]interface{}{
				"user": map[string]interface{}{
					"id":       12345,
					"username": "john_doe",
					"settings": map[string]interface{}{
						"theme":         "dark",
						"language":      "en",
						"notifications": true,
					},
				},
				"permissions": []interface{}{
					"read",
					"write",
					map[string]interface{}{
						"resource": "users",
						"actions":  []interface{}{"create", "delete"},
					},
				},
			},
		},
		{
			name:    "empty object",
			payload: map[string]interface{}{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Encode
			noisy, err := noise.Encode(tc.payload)
			if err != nil {
				t.Fatalf("Encode failed: %v", err)
			}

			// Validate encoded structure
			validateNoisyPayload(t, noisy)

			// Test JSON marshaling/unmarshaling
			jsonBytes, err := json.Marshal(noisy)
			if err != nil {
				t.Fatalf("JSON marshal failed: %v", err)
			}

			var unmarshaled NoisyPayload
			if err := json.Unmarshal(jsonBytes, &unmarshaled); err != nil {
				t.Fatalf("JSON unmarshal failed: %v", err)
			}

			// Decode
			decoded, err := noise.Decode(&unmarshaled)
			if err != nil {
				t.Fatalf("Decode failed: %v", err)
			}

			// Compare original and decoded payloads
			if !reflect.DeepEqual(tc.payload, decoded) {
				t.Errorf("Payload mismatch.\nExpected: %+v\nGot: %+v", tc.payload, decoded)
			}
		})
	}
}

func TestPayloadNoise_InvalidInputs(t *testing.T) {
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)

	noise, err := NewPayloadNoise()
	if err != nil {
		t.Fatalf("Failed to create PayloadNoise: %v", err)
	}

	// Test encoding nil payload
	if _, err := noise.Encode(nil); err == nil {
		t.Error("Expected error when encoding nil payload")
	}

	// Test invalid noisy payloads for decoding
	invalidPayloads := []*NoisyPayload{
		nil,
		{}, // Empty payload
		{
			Version:   "1.0",
			Timestamp: 123,
			Salt:      "invalid",
			IV:        "invalid",
			Hash:      "invalid",
			Data:      map[string]string{},
		},
		{
			Version:   "",
			Timestamp: 0,
			Salt:      "",
			IV:        "",
			Hash:      "",
			Data:      nil,
		},
	}

	for i, payload := range invalidPayloads {
		if _, err := noise.Decode(payload); err == nil {
			t.Errorf("Case %d: Expected error for invalid payload", i)
		}
	}
}

func TestPayloadNoise_HashVerification(t *testing.T) {
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)

	noise, err := NewPayloadNoise()
	if err != nil {
		t.Fatalf("Failed to create PayloadNoise: %v", err)
	}

	// Create a test payload
	payload := map[string]interface{}{
		"test": "data",
	}

	// Encode the payload
	noisy, err := noise.Encode(payload)
	if err != nil {
		t.Fatalf("Failed to encode payload: %v", err)
	}

	// Test with valid hash
	if _, err := noise.Decode(noisy); err != nil {
		t.Errorf("Hash verification failed for valid payload: %v", err)
	}

	// Test with tampered hash
	originalHash := noisy.Hash
	noisy.Hash = "tampered_hash"
	if _, err := noise.Decode(noisy); err == nil {
		t.Error("Expected error for tampered hash")
	}

	// Test with tampered data but original hash
	noisy.Hash = originalHash
	noisy.Data["tampered_key"] = "tampered_value"
	if _, err := noise.Decode(noisy); err == nil {
		t.Error("Expected error for tampered data")
	}
}

func TestPayloadNoise_KeyRecovery(t *testing.T) {
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)

	noise, err := NewPayloadNoise()
	if err != nil {
		t.Fatalf("Failed to create PayloadNoise: %v", err)
	}

	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    "key_12345678",
			expected: "key",
		},
		{
			input:    "complex.key.with.dots_abcdef12",
			expected: "complex.key.with.dots",
		},
		{
			input:    "with_underscore_in_key_12345678",
			expected: "with_underscore_in_key",
		},
		{
			input:    "simple",
			expected: "simple",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := noise.recoverOriginalKey(tc.input)
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

// Helper function to validate NoisyPayload structure
func validateNoisyPayload(t *testing.T, payload *NoisyPayload) {
	t.Helper()

	if payload.Version != "1.0" {
		t.Errorf("Expected version 1.0, got %s", payload.Version)
	}
	if payload.Timestamp == 0 {
		t.Error("Timestamp should not be 0")
	}
	if payload.Salt == "" {
		t.Error("Salt should not be empty")
	}
	if payload.IV == "" {
		t.Error("IV should not be empty")
	}
	if payload.Hash == "" {
		t.Error("Hash should not be empty")
	}
	if payload.Data == nil {
		t.Error("Data should not be nil")
	}
}
