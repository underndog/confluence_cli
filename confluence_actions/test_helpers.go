package confluence_actions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type ValidationTestCase struct {
	Name        string
	Args        []string
	ExpectError bool
	ErrorMsg    string
}

// RunValidationTable runs a table of validation test cases for a given validation function.
func RunValidationTable(t *testing.T, fn func(args ...string) error, cases []ValidationTestCase) {
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			err := fn(tc.Args...)
			if tc.ExpectError {
				assert.Error(t, err)
				if tc.ErrorMsg != "" {
					assert.Contains(t, err.Error(), tc.ErrorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
