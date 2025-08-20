# Test Commands for Confluence CLI

Below are the commands to run tests for the Confluence CLI project:

## 🚀 Quick Start Commands

### 1. Run all tests
```bash
go test ./...
```

### 2. Run tests with detailed output
```bash
go test -v ./...
```

### 3. Run tests with coverage
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## 📦 Package-Specific Tests

### 4. Test helper functions
```bash
go test -v ./helper/
```

### 5. Test model structures
```bash
go test -v ./model/
```

### 6. Test commands structure
```bash
go test -v ./confluence_commands/
```

### 7. Test actions (only run logic tests, don't run CLI context tests)
```bash
go test -v -run TestValidationLogic ./confluence_actions/
go test -v -run TestValidationLogicWithMock ./confluence_actions/
go test -v -run TestFileValidationLogic ./confluence_actions/
go test -v -run TestUpdatePageValidationLogic ./confluence_actions/
go test -v -run TestUpdatePageFileValidationLogic ./confluence_actions/
go test -v -run TestUpdatePageValidationLogicDirect ./confluence_actions/
go test -v -run TestUploadAttachmentValidationLogic ./confluence_actions/
go test -v -run TestUploadAttachmentFileValidationLogic ./confluence_actions/
go test -v -run TestUploadAttachmentValidationLogicDirect ./confluence_actions/
```

## 🎯 Specific Test Functions

### 8. Test validation logic (passed)
```bash
go test -v -run TestValidationLogic ./confluence_actions/
```

### 9. Test file validation logic (passed)
```bash
go test -v -run TestFileValidationLogic ./confluence_actions/
```

### 10. Test update page validation logic (passed)
```bash
go test -v -run TestUpdatePageValidationLogic ./confluence_actions/
go test -v -run TestUpdatePageFileValidationLogic ./confluence_actions/
go test -v -run TestUpdatePageValidationLogicDirect ./confluence_actions/
```

### 11. Test upload attachment validation logic (passed)
```bash
go test -v -run TestUploadAttachmentValidationLogic ./confluence_actions/
go test -v -run TestUploadAttachmentFileValidationLogic ./confluence_actions/
go test -v -run TestUploadAttachmentValidationLogicDirect ./confluence_actions/
```

### 12. Test environment variables
```bash
go test -v -run TestGetEnvOrDefault ./helper/
```

### 13. Test JSON marshaling
```bash
go test -v -run TestPage_JSONMarshal ./model/req
```

## 📝 Notes

- Only run logic validation, file validation, update page validation logic tests in `confluence_actions/`
- Action tests depending on CLI context, logging, HTTP client have been deleted
- Tests in `helper/` and `model/` work well
- Tests in `confluence_commands/` test command structure, not logic
- To test completely, need to mock HTTP client and initialize logging system

## 🔄 Continuous Integration

To run tests as in CI:
```bash
# Clean and run tests
go clean -testcache
go mod tidy
go test -v -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
``` 