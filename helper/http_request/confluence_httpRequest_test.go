package http_request

import (
	"confluence_cli/log"
	"os"
	"testing"
)

func setupTestEnvironment(t *testing.T) func() {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping integration test in CI environment")
	}
	log.InitLogger(true)
	os.Setenv("EMAIL", "test@example.com")
	os.Setenv("API_TOKEN", "testtoken")
	os.Setenv("CONFLUENCE_URL", "http://localhost")

	return func() {
		os.Unsetenv("EMAIL")
		os.Unsetenv("API_TOKEN")
		os.Unsetenv("CONFLUENCE_URL")
	}
}

func TestCreateConfluencePage(t *testing.T) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	resp, err := CreateConfluencePage("SPACE", "PARENT", "TestTitle", "TestBody")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if resp == nil {
		t.Error("expected response, got nil")
	}
}

func TestGetConfluencePagesByTitle(t *testing.T) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	resp, err := GetConfluencePagesByTitle("TestTitle")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if resp == nil {
		t.Error("expected response, got nil")
	}
}

// Test environment setup without HTTP calls
func TestEnvironmentSetup(t *testing.T) {
	// Save original env
	originalEmail := os.Getenv("EMAIL")
	originalToken := os.Getenv("API_TOKEN")
	originalURL := os.Getenv("CONFLUENCE_URL")

	// Restore after test
	defer func() {
		os.Setenv("EMAIL", originalEmail)
		os.Setenv("API_TOKEN", originalToken)
		os.Setenv("CONFLUENCE_URL", originalURL)
	}()

	// Set test values
	os.Setenv("EMAIL", "test@example.com")
	os.Setenv("API_TOKEN", "testtoken")
	os.Setenv("CONFLUENCE_URL", "http://localhost")

	// Verify
	if os.Getenv("EMAIL") != "test@example.com" {
		t.Error("EMAIL environment variable not set correctly")
	}
	if os.Getenv("API_TOKEN") != "testtoken" {
		t.Error("API_TOKEN environment variable not set correctly")
	}
	if os.Getenv("CONFLUENCE_URL") != "http://localhost" {
		t.Error("CONFLUENCE_URL environment variable not set correctly")
	}
}
