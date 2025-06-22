package http_request

import (
	"os"
	"testing"
	"confluence_cli/log"
)

func TestCreateConfluencePage(t *testing.T) {
	log.InitLogger(true)
	os.Setenv("EMAIL", "test@example.com")
	os.Setenv("API_TOKEN", "testtoken")
	os.Setenv("CONFLUENCE_URL", "http://localhost")
	defer os.Unsetenv("EMAIL")
	defer os.Unsetenv("API_TOKEN")
	defer os.Unsetenv("CONFLUENCE_URL")

	resp, err := CreateConfluencePage("SPACE", "PARENT", "TestTitle", "TestBody")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if resp == nil {
		t.Error("expected response, got nil")
	}
}

func TestGetConfluencePagesByTitle(t *testing.T) {
	log.InitLogger(true)
	os.Setenv("EMAIL", "test@example.com")
	os.Setenv("API_TOKEN", "testtoken")
	os.Setenv("CONFLUENCE_URL", "http://localhost")
	defer os.Unsetenv("EMAIL")
	defer os.Unsetenv("API_TOKEN")
	defer os.Unsetenv("CONFLUENCE_URL")

	resp, err := GetConfluencePagesByTitle("TestTitle")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if resp == nil {
		t.Error("expected response, got nil")
	}
}
