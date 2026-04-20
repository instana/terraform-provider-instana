package provider

import (
	"testing"
)

func TestSanitizeSensitiveData(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "sanitize lowercase token",
			input:    "API request with token: sk_live_123456",
			expected: "API request with token: [REDACTED]",
		},
		{
			name:     "sanitize uppercase TOKEN",
			input:    "Config: TOKEN: abc123",
			expected: "Config: TOKEN: [REDACTED]",
		},
		{
			name:     "sanitize mixed case ApiToken",
			input:    "Using ApiToken: secret123",
			expected: "Using ApiToken: [REDACTED]",
		},
		{
			name:     "sanitize authorization header",
			input:    "Authorization: Bearer xyz789",
			expected: "Authorization: [REDACTED] xyz789",
		},
		{
			name:     "sanitize password",
			input:    "Login failed with password: mypass123",
			expected: "Login failed with password: [REDACTED]",
		},
		{
			name:     "sanitize API_TOKEN with underscore",
			input:    "Config {API_TOKEN: secret, URL: https://example.com}",
			expected: "Config {API_TOKEN: [REDACTED], URL: https://example.com}",
		},
		{
			name:     "sanitize secret key",
			input:    "Using secret: mysecret and key: mykey",
			expected: "Using secret: [REDACTED] and key: mykey",
		},
		{
			name:     "no sensitive data",
			input:    "Normal log message without sensitive data",
			expected: "Normal log message without sensitive data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeSensitiveData(tt.input)
			if result != tt.expected {
				t.Errorf("sanitizeSensitiveData() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSanitizeLogFields(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		expected map[string]interface{}
	}{
		{
			name: "sanitize token field",
			input: map[string]interface{}{
				"token":    "secret123",
				"endpoint": "https://example.com",
			},
			expected: map[string]interface{}{
				"token":    "[REDACTED]",
				"endpoint": "https://example.com",
			},
		},
		{
			name: "sanitize apiToken field",
			input: map[string]interface{}{
				"apiToken": "abc123",
				"url":      "https://api.example.com",
			},
			expected: map[string]interface{}{
				"apiToken": "[REDACTED]",
				"url":      "https://api.example.com",
			},
		},
		{
			name: "sanitize password field",
			input: map[string]interface{}{
				"username": "admin",
				"password": "secret",
			},
			expected: map[string]interface{}{
				"username": "admin",
				"password": "[REDACTED]",
			},
		},
		{
			name: "sanitize string value containing token",
			input: map[string]interface{}{
				"message": "Using token: abc123",
				"status":  "success",
			},
			expected: map[string]interface{}{
				"message": "Using token: [REDACTED]",
				"status":  "success",
			},
		},
		{
			name: "no sensitive data",
			input: map[string]interface{}{
				"endpoint": "https://example.com",
				"status":   "ok",
			},
			expected: map[string]interface{}{
				"endpoint": "https://example.com",
				"status":   "ok",
			},
		},
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeLogFields(tt.input)

			if tt.input == nil {
				if result != nil {
					t.Errorf("sanitizeLogFields() = %v, want nil", result)
				}
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("sanitizeLogFields() returned %d fields, want %d", len(result), len(tt.expected))
			}

			for key, expectedValue := range tt.expected {
				if result[key] != expectedValue {
					t.Errorf("sanitizeLogFields()[%s] = %v, want %v", key, result[key], expectedValue)
				}
			}
		})
	}
}
