package pip

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestNewPipError(t *testing.T) {
	errorType := ErrorTypePipNotInstalled
	message := "pip is not installed"

	err := NewPipError(errorType, message)

	if err.Type != errorType {
		t.Errorf("Type = %s, want %s", err.Type, errorType)
	}

	if err.Message != message {
		t.Errorf("Message = %s, want %s", err.Message, message)
	}

	if err.Context == nil {
		t.Error("Context should not be nil")
	}
}

func TestNewCommandError(t *testing.T) {
	command := "pip install nonexistent"
	output := "ERROR: Could not find a version"
	exitCode := 1
	cause := errors.New("command failed")

	err := NewCommandError(command, output, exitCode, cause)

	if err.Type != ErrorTypeCommandFailed {
		t.Errorf("Type = %s, want %s", err.Type, ErrorTypeCommandFailed)
	}

	if err.Command != command {
		t.Errorf("Command = %s, want %s", err.Command, command)
	}

	if err.Output != output {
		t.Errorf("Output = %s, want %s", err.Output, output)
	}

	if err.ExitCode != exitCode {
		t.Errorf("ExitCode = %d, want %d", err.ExitCode, exitCode)
	}

	if err.Cause != cause {
		t.Errorf("Cause = %v, want %v", err.Cause, cause)
	}
}

func TestPipErrorDetailsError(t *testing.T) {
	tests := []struct {
		name     string
		err      *PipErrorDetails
		contains []string // strings that should be in the error message
	}{
		{
			name: "basic error",
			err: &PipErrorDetails{
				Type:    ErrorTypePipNotInstalled,
				Message: "pip not found",
			},
			contains: []string{"pip_not_installed", "pip not found"},
		},
		{
			name: "error with command",
			err: &PipErrorDetails{
				Type:    ErrorTypeCommandFailed,
				Message: "command failed",
				Command: "pip install requests",
			},
			contains: []string{"command_failed", "command failed", "pip install requests"},
		},
		{
			name: "error with exit code",
			err: &PipErrorDetails{
				Type:     ErrorTypeCommandFailed,
				Message:  "command failed",
				ExitCode: 1,
			},
			contains: []string{"command_failed", "Exit Code: 1"},
		},
		{
			name: "error with suggestions",
			err: &PipErrorDetails{
				Type:        ErrorTypePermissionDenied,
				Message:     "permission denied",
				Suggestions: []string{"use sudo", "try virtual environment"},
			},
			contains: []string{"permission_denied", "use sudo", "try virtual environment"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errorStr := tt.err.Error()

			for _, contain := range tt.contains {
				if !strings.Contains(errorStr, contain) {
					t.Errorf("Error string should contain '%s', got: %s", contain, errorStr)
				}
			}
		})
	}
}

func TestPipErrorDetailsUnwrap(t *testing.T) {
	cause := errors.New("underlying error")
	err := &PipErrorDetails{
		Type:  ErrorTypeCommandFailed,
		Cause: cause,
	}

	if err.Unwrap() != cause {
		t.Errorf("Unwrap() = %v, want %v", err.Unwrap(), cause)
	}
}

func TestPipErrorDetailsIs(t *testing.T) {
	err1 := &PipErrorDetails{Type: ErrorTypePipNotInstalled}
	err2 := &PipErrorDetails{Type: ErrorTypePipNotInstalled}
	err3 := &PipErrorDetails{Type: ErrorTypePythonNotFound}

	if !err1.Is(err2) {
		t.Error("err1.Is(err2) should be true for same error types")
	}

	if err1.Is(err3) {
		t.Error("err1.Is(err3) should be false for different error types")
	}

	if err1.Is(errors.New("generic error")) {
		t.Error("err1.Is(generic error) should be false")
	}
}

func TestPipErrorDetailsWithMethods(t *testing.T) {
	err := NewPipError(ErrorTypePipNotInstalled, "pip not found")

	// Test WithSuggestion
	err = err.WithSuggestion("install pip")
	if len(err.Suggestions) != 1 || err.Suggestions[0] != "install pip" {
		t.Error("WithSuggestion failed")
	}

	// Test WithContext
	err = err.WithContext("os", "linux")
	if err.Context["os"] != "linux" {
		t.Error("WithContext failed")
	}

	// Test WithCause
	cause := errors.New("underlying cause")
	err = err.WithCause(cause)
	if err.Cause != cause {
		t.Error("WithCause failed")
	}
}

func TestAddSuggestionsFromOutput(t *testing.T) {
	tests := []struct {
		name                string
		output              string
		expectedSuggestions []string
	}{
		{
			name:   "permission denied",
			output: "Permission denied when trying to install",
			expectedSuggestions: []string{
				"Try running with elevated privileges",
				"Consider using a virtual environment",
			},
		},
		{
			name:   "no module named pip",
			output: "No module named pip",
			expectedSuggestions: []string{
				"Install pip using: python -m ensurepip --upgrade",
				"Or download get-pip.py",
			},
		},
		{
			name:   "could not find version",
			output: "Could not find a version that satisfies the requirement",
			expectedSuggestions: []string{
				"Check if the package name is spelled correctly",
				"Try searching for the package on PyPI",
			},
		},
		{
			name:   "network error",
			output: "Network connection failed",
			expectedSuggestions: []string{
				"Check your internet connection",
				"Try using a different package index",
			},
		},
		{
			name:   "timeout",
			output: "Request timeout occurred",
			expectedSuggestions: []string{
				"Increase timeout value",
				"Try again later",
			},
		},
		{
			name:   "disk space",
			output: "No space left on device",
			expectedSuggestions: []string{
				"Free up disk space",
				"Clean pip cache",
			},
		},
		{
			name:   "already satisfied",
			output: "Requirement already satisfied",
			expectedSuggestions: []string{
				"Use --upgrade flag",
				"Use --force-reinstall",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &PipErrorDetails{
				Type:    ErrorTypeCommandFailed,
				Message: "test error",
				Context: make(map[string]string),
			}

			err.addSuggestionsFromOutput(tt.output)

			for _, expectedSuggestion := range tt.expectedSuggestions {
				found := false
				for _, suggestion := range err.Suggestions {
					if strings.Contains(suggestion, expectedSuggestion) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected suggestion containing '%s' not found in %v", expectedSuggestion, err.Suggestions)
				}
			}
		})
	}
}

func TestIsErrorType(t *testing.T) {
	pipErr := NewPipError(ErrorTypePipNotInstalled, "pip not found")
	genericErr := errors.New("generic error")

	if !IsErrorType(pipErr, ErrorTypePipNotInstalled) {
		t.Error("IsErrorType should return true for matching pip error type")
	}

	if IsErrorType(pipErr, ErrorTypePythonNotFound) {
		t.Error("IsErrorType should return false for non-matching pip error type")
	}

	if IsErrorType(genericErr, ErrorTypePipNotInstalled) {
		t.Error("IsErrorType should return false for generic error")
	}
}

func TestGetErrorType(t *testing.T) {
	pipErr := NewPipError(ErrorTypePipNotInstalled, "pip not found")
	genericErr := errors.New("generic error")

	if GetErrorType(pipErr) != ErrorTypePipNotInstalled {
		t.Error("GetErrorType should return correct type for pip error")
	}

	if GetErrorType(genericErr) != "" {
		t.Error("GetErrorType should return empty string for generic error")
	}
}

func TestWrapError(t *testing.T) {
	originalErr := errors.New("original error")
	wrappedErr := WrapError(originalErr, ErrorTypeCommandFailed, "wrapped message")

	if wrappedErr.Type != ErrorTypeCommandFailed {
		t.Error("WrapError should set correct error type")
	}

	if wrappedErr.Message != "wrapped message" {
		t.Error("WrapError should set correct message")
	}

	if wrappedErr.Cause != originalErr {
		t.Error("WrapError should set correct cause")
	}
}

func TestErrorHandler(t *testing.T) {
	// Create a test logger
	logger, err := NewLogger(DefaultLoggerConfig())
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	handler := NewErrorHandler(logger)

	// Test handling nil error
	result := handler.Handle(nil, "test context")
	if result != nil {
		t.Error("Handle(nil) should return nil")
	}

	// Test handling pip error
	pipErr := NewPipError(ErrorTypePipNotInstalled, "pip not found")
	result = handler.Handle(pipErr, "test context")
	if result != pipErr {
		t.Error("Handle(pipErr) should return the same error")
	}

	// Test handling generic error
	genericErr := errors.New("generic error")
	result = handler.Handle(genericErr, "test context")
	if result == nil {
		t.Error("Handle(genericErr) should return wrapped error")
	}

	if _, ok := result.(*PipErrorDetails); !ok {
		t.Error("Handle(genericErr) should return PipErrorDetails")
	}
}

func TestHandlePipErrorDetailed(t *testing.T) {
	// Create a test logger
	logger, err := NewLogger(DefaultLoggerConfig())
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	handler := NewErrorHandler(logger)

	tests := []struct {
		name      string
		pipError  *PipError
		context   string
		expectLog bool
	}{
		{
			name: "system dependency error",
			pipError: &PipError{
				Type:    string(ErrorTypePipNotInstalled),
				Message: "pip not found",
			},
			context:   "system dependency test",
			expectLog: true,
		},
		{
			name: "permission error",
			pipError: &PipError{
				Type:    string(ErrorTypePermissionDenied),
				Message: "permission denied",
			},
			context:   "permission test",
			expectLog: true,
		},
		{
			name: "network error",
			pipError: &PipError{
				Type:    string(ErrorTypeNetworkError),
				Message: "network timeout",
			},
			context:   "network test",
			expectLog: true,
		},
		{
			name: "package not found error",
			pipError: &PipError{
				Type:    string(ErrorTypePackageNotFound),
				Message: "package not found",
			},
			context:   "package test",
			expectLog: true,
		},
		{
			name: "unknown error type",
			pipError: &PipError{
				Type:    "unknown_error",
				Message: "unknown error",
			},
			context:   "unknown test",
			expectLog: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Handle method
			result := handler.Handle(tt.pipError, tt.context)

			// Handle() returns a PipErrorDetails, not the original PipError
			if result == nil {
				t.Errorf("Handle() should not return nil")
			}

			// The function should not panic and should log appropriately
			// We can't easily test the logging output without capturing it,
			// but we can ensure the function completes without error
		})
	}
}

func TestErrorHandlerWithDifferentContexts(t *testing.T) {
	// Create a test logger
	logger, err := NewLogger(DefaultLoggerConfig())
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	handler := NewErrorHandler(logger)

	pipErr := NewPipError(ErrorTypeCommandFailed, "test error")

	contexts := []string{
		"",
		"installation context",
		"virtual environment context",
		"package management context",
		"very long context string that might test buffer limits and edge cases in logging",
	}

	for _, context := range contexts {
		t.Run(fmt.Sprintf("context_%s", context), func(t *testing.T) {
			result := handler.Handle(pipErr, context)
			if result != pipErr {
				t.Errorf("Handle() should return the same error regardless of context")
			}
		})
	}
}

// Benchmark tests
func BenchmarkNewPipError(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewPipError(ErrorTypePipNotInstalled, "pip not found")
	}
}

func BenchmarkPipErrorDetailsError(b *testing.B) {
	err := &PipErrorDetails{
		Type:        ErrorTypeCommandFailed,
		Message:     "command failed",
		Command:     "pip install requests",
		ExitCode:    1,
		Suggestions: []string{"try again", "check network"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkAddSuggestionsFromOutput(b *testing.B) {
	err := &PipErrorDetails{
		Type:    ErrorTypeCommandFailed,
		Message: "test error",
		Context: make(map[string]string),
	}
	output := "Permission denied when trying to install package"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err.Suggestions = nil // Reset suggestions
		err.addSuggestionsFromOutput(output)
	}
}
