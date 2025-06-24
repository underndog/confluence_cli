package helper

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvOrDefault(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "Environment variable exists",
			key:          "TEST_VAR",
			defaultValue: "default",
			envValue:     "env_value",
			expected:     "env_value",
		},
		{
			name:         "Environment variable does not exist",
			key:          "NON_EXISTENT_VAR",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
		{
			name:         "Empty environment variable",
			key:          "EMPTY_VAR",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable if provided
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			result := GetEnvOrDefault(tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetEnvOrDefault_ConfluenceURL(t *testing.T) {
	// Test with actual Confluence environment variables
	tests := []struct {
		name          string
		envURL        string
		envEmail      string
		envToken      string
		expectedURL   string
		expectedEmail string
		expectedToken string
	}{
		{
			name:          "All environment variables set",
			envURL:        "https://test.atlassian.net/",
			envEmail:      "test@example.com",
			envToken:      "test-token",
			expectedURL:   "https://test.atlassian.net/",
			expectedEmail: "test@example.com",
			expectedToken: "test-token",
		},
		{
			name:          "No environment variables set",
			envURL:        "",
			envEmail:      "",
			envToken:      "",
			expectedURL:   "https://nimtechnology.atlassian.net",
			expectedEmail: "dc.nim94@gmail.com",
			expectedToken: "nimtechnology",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			if tt.envURL != "" {
				os.Setenv("CONFLUENCE_URL", tt.envURL)
				defer os.Unsetenv("CONFLUENCE_URL")
			}
			if tt.envEmail != "" {
				os.Setenv("EMAIL", tt.envEmail)
				defer os.Unsetenv("EMAIL")
			}
			if tt.envToken != "" {
				os.Setenv("API_TOKEN", tt.envToken)
				defer os.Unsetenv("API_TOKEN")
			}

			url := GetEnvOrDefault("CONFLUENCE_URL", "https://nimtechnology.atlassian.net")
			email := GetEnvOrDefault("EMAIL", "dc.nim94@gmail.com")
			token := GetEnvOrDefault("API_TOKEN", "nimtechnology")

			assert.Equal(t, tt.expectedURL, url)
			assert.Equal(t, tt.expectedEmail, email)
			assert.Equal(t, tt.expectedToken, token)
		})
	}
}
