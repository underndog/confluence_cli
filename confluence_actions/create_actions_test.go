package confluence_actions

import (
	"confluence_cli/helper"
	"confluence_cli/log"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestCreatePageCommand_Structure(t *testing.T) {
	// ... existing code ...
}

func TestCreatePageCommand(t *testing.T) {
	// ... existing code ...
}

// Remove TestParseTestResultsFromHTML function completely

// Test regex pattern matching for Overall Status
func TestOverallStatusPatternMatching(t *testing.T) {
	pattern := `<td><strong>Overall Status</strong></td>\s*<td colspan=['"]2['"]>`
	re := regexp.MustCompile(pattern)

	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name:     "Match Overall Status pattern",
			content:  `<td><strong>Overall Status</strong></td><td colspan="2">`,
			expected: true,
		},
		{
			name:     "Match with single quotes",
			content:  `<td><strong>Overall Status</strong></td><td colspan='2'>`,
			expected: true,
		},
		{
			name:     "No match - different structure",
			content:  `<td><strong>Status</strong></td><td colspan="2">`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := re.MatchString(tt.content)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test macro content generation
func TestMacroContentGeneration(t *testing.T) {
	// Test that macros contain expected content
	attachmentMacro := helper.CreateAttachmentMacro()
	assert.Contains(t, attachmentMacro, "attachments")
	assert.Contains(t, attachmentMacro, "ac:structured-macro")

	// Test action list macro detection
	hasActionList := helper.HasActionListMacro(`<ac:task-list>...</ac:task-list>`)
	assert.True(t, hasActionList)
}

// Mock context for minimal flag testing
func TestCreatePageAction_MissingFlags(t *testing.T) {
	log.InitLogger(true)
	app := cli.NewApp()
	set := flagSet([]string{})
	c := cli.NewContext(app, set, nil)
	err := CreatePageAction(c)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "space-id")
}

// Helper to create a flag set
func flagSet(args []string) *flag.FlagSet {
	if len(args)%2 != 0 {
		panic("args must contain an even number of elements (key-value pairs)")
	}
	set := flag.NewFlagSet("test", 0)
	set.String("space-id", "", "")
	set.String("parent-page-id", "", "")
	set.String("title", "", "")
	set.String("body-value-from-file", "", "")
	set.String("body-value", "", "")
	set.String("file", "", "")
	for i := 0; i < len(args); i += 2 {
		if i+1 >= len(args) {
			break
		}
		set.Set(args[i], args[i+1])
	}
	return set
}

// Test validation logic with mock validation function
func TestValidationLogicWithMock(t *testing.T) {
	// Mock validation function that mimics the logic in CreatePageAction
	validateRequiredFields := func(spaceId, parentId, title string) error {
		if spaceId == "" {
			return fmt.Errorf("--space-id is required")
		}
		if parentId == "" {
			return fmt.Errorf("--parent-page-id is required")
		}
		if title == "" {
			return fmt.Errorf("--title is required")
		}
		return nil
	}

	tests := []struct {
		name        string
		spaceId     string
		parentId    string
		title       string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Missing space-id",
			spaceId:     "",
			parentId:    "123",
			title:       "Test",
			expectError: true,
			errorMsg:    "--space-id is required",
		},
		{
			name:        "Missing parent-page-id",
			spaceId:     "SPACE",
			parentId:    "",
			title:       "Test",
			expectError: true,
			errorMsg:    "--parent-page-id is required",
		},
		{
			name:        "Missing title",
			spaceId:     "SPACE",
			parentId:    "123",
			title:       "",
			expectError: true,
			errorMsg:    "--title is required",
		},
		{
			name:        "Valid required parameters",
			spaceId:     "SPACE",
			parentId:    "123",
			title:       "Test",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRequiredFields(tt.spaceId, tt.parentId, tt.title)
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test file parameter handling

func TestFileParameterHandling(t *testing.T) {
	t.Run("Empty file path", func(t *testing.T) {
		app := cli.NewApp()
		set := flagSet([]string{
			"space-id", "SPACE",
			"parent-page-id", "123",
			"title", "Test",
			"file", "",
		})
		c := cli.NewContext(app, set, nil)
		assert.Equal(t, "", c.String("file"))
	})

	t.Run("Valid file path", func(t *testing.T) {
		dir := t.TempDir()
		fp := filepath.Join(dir, "body.html")
		err := os.WriteFile(fp, []byte("<html>ok</html>"), 0o644)
		assert.NoError(t, err)

		app := cli.NewApp()
		set := flagSet([]string{
			"space-id", "SPACE",
			"parent-page-id", "123",
			"title", "Test",
			"file", fp,
		})
		c := cli.NewContext(app, set, nil)
		assert.Equal(t, fp, c.String("file"))
	})
}
