package http_request

import (
	"os"
	"testing"
	"confluence_cli/log"
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

	resp, err = CreateConfluencePage("SPACE", "PARENT", "TestTitle", "TestBody")
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
