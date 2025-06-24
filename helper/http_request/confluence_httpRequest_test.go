package http_request

import (
	"confluence_cli/log"
	"os"
	"testing"
)

func setupTestEnvironment(t *testing.T) func() {
	// Skip tests if running in CI or if no local server
	if os.Getenv("CI") != "" || os.Getenv("SKIP_HTTP_TESTS") != "" {
		t.Skip("Skipping HTTP integration tests")
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
	// Skip this test as it requires a real HTTP server
	t.Skip("Skipping HTTP integration test - requires real server")

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
	// Skip this test as it requires a real HTTP server
	t.Skip("Skipping HTTP integration test - requires real server")

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
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Test that environment variables are set correctly
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
