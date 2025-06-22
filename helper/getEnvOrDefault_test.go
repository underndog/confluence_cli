package helper

import (
	"os"
	"testing"
)

func TestGetEnvOrDefault(t *testing.T) {
	os.Setenv("TEST_KEY", "test_value")
	defer os.Unsetenv("TEST_KEY")

	if val := GetEnvOrDefault("TEST_KEY", "default"); val != "test_value" {
		t.Errorf("expected 'test_value', got '%s'", val)
	}
	if val := GetEnvOrDefault("NON_EXISTENT_KEY", "default"); val != "default" {
		t.Errorf("expected 'default', got '%s'", val)
	}
}
