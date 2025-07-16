package pip

import (
	"fmt"
	"strings"
)

// ErrorType represents different types of pip errors
type ErrorType string

const (
	ErrorTypePipNotInstalled    ErrorType = "pip_not_installed"
	ErrorTypePythonNotFound     ErrorType = "python_not_found"
	ErrorTypePackageNotFound    ErrorType = "package_not_found"
	ErrorTypeVenvNotFound       ErrorType = "venv_not_found"
	ErrorTypeVenvAlreadyExists  ErrorType = "venv_already_exists"
	ErrorTypeInvalidPackageSpec ErrorType = "invalid_package_spec"
	ErrorTypeCommandFailed      ErrorType = "command_failed"
	ErrorTypePermissionDenied   ErrorType = "permission_denied"
	ErrorTypeNetworkError       ErrorType = "network_error"
	ErrorTypeUnsupportedOS      ErrorType = "unsupported_os"
	ErrorTypeFileNotFound       ErrorType = "file_not_found"
	ErrorTypeInvalidPath        ErrorType = "invalid_path"
	ErrorTypeFeatureDisabled    ErrorType = "feature_disabled"
	ErrorTypeTimeout            ErrorType = "timeout"
	ErrorTypeInvalidConfig      ErrorType = "invalid_config"
)

// PipErrorDetails provides additional context for errors
type PipErrorDetails struct {
	Type        ErrorType         `json:"type"`
	Message     string            `json:"message"`
	Command     string            `json:"command,omitempty"`
	Output      string            `json:"output,omitempty"`
	ExitCode    int               `json:"exit_code,omitempty"`
	Suggestions []string          `json:"suggestions,omitempty"`
	Context     map[string]string `json:"context,omitempty"`
	Cause       error             `json:"cause,omitempty"`
}

// Error implements the error interface
func (e *PipErrorDetails) Error() string {
	var parts []string

	parts = append(parts, fmt.Sprintf("[%s] %s", e.Type, e.Message))

	if e.Command != "" {
		parts = append(parts, fmt.Sprintf("Command: %s", e.Command))
	}

	if e.ExitCode != 0 {
		parts = append(parts, fmt.Sprintf("Exit Code: %d", e.ExitCode))
	}

	if e.Output != "" && len(e.Output) < 200 {
		parts = append(parts, fmt.Sprintf("Output: %s", e.Output))
	}

	if len(e.Suggestions) > 0 {
		parts = append(parts, fmt.Sprintf("Suggestions: %s", strings.Join(e.Suggestions, ", ")))
	}

	return strings.Join(parts, " | ")
}

// Unwrap returns the underlying cause
func (e *PipErrorDetails) Unwrap() error {
	return e.Cause
}

// Is checks if the error is of a specific type
func (e *PipErrorDetails) Is(target error) bool {
	if t, ok := target.(*PipErrorDetails); ok {
		return e.Type == t.Type
	}
	return false
}

// NewPipError creates a new pip error with details
func NewPipError(errorType ErrorType, message string) *PipErrorDetails {
	return &PipErrorDetails{
		Type:    errorType,
		Message: message,
		Context: make(map[string]string),
	}
}

// NewCommandError creates a new command execution error
func NewCommandError(command string, output string, exitCode int, cause error) *PipErrorDetails {
	err := &PipErrorDetails{
		Type:     ErrorTypeCommandFailed,
		Message:  "command execution failed",
		Command:  command,
		Output:   output,
		ExitCode: exitCode,
		Cause:    cause,
		Context:  make(map[string]string),
	}

	// Add suggestions based on common error patterns
	err.addSuggestionsFromOutput(output)

	return err
}

// WithSuggestion adds a suggestion to the error
func (e *PipErrorDetails) WithSuggestion(suggestion string) *PipErrorDetails {
	e.Suggestions = append(e.Suggestions, suggestion)
	return e
}

// WithContext adds context information to the error
func (e *PipErrorDetails) WithContext(key, value string) *PipErrorDetails {
	if e.Context == nil {
		e.Context = make(map[string]string)
	}
	e.Context[key] = value
	return e
}

// WithCause sets the underlying cause of the error
func (e *PipErrorDetails) WithCause(cause error) *PipErrorDetails {
	e.Cause = cause
	return e
}

// addSuggestionsFromOutput analyzes command output and adds relevant suggestions
func (e *PipErrorDetails) addSuggestionsFromOutput(output string) {
	output = strings.ToLower(output)

	if strings.Contains(output, "permission denied") {
		e.Suggestions = append(e.Suggestions, "Try running with elevated privileges (sudo on Unix, Run as Administrator on Windows)")
		e.Suggestions = append(e.Suggestions, "Consider using a virtual environment to avoid permission issues")
	}

	if strings.Contains(output, "no module named pip") {
		e.Suggestions = append(e.Suggestions, "Install pip using: python -m ensurepip --upgrade")
		e.Suggestions = append(e.Suggestions, "Or download get-pip.py and run: python get-pip.py")
	}

	if strings.Contains(output, "could not find a version") {
		e.Suggestions = append(e.Suggestions, "Check if the package name is spelled correctly")
		e.Suggestions = append(e.Suggestions, "Try searching for the package on PyPI: https://pypi.org")
		e.Suggestions = append(e.Suggestions, "Check if you need to specify a different Python version")
	}

	if strings.Contains(output, "network") || strings.Contains(output, "connection") {
		e.Suggestions = append(e.Suggestions, "Check your internet connection")
		e.Suggestions = append(e.Suggestions, "Try using a different package index: --index-url")
		e.Suggestions = append(e.Suggestions, "Consider using a proxy if you're behind a corporate firewall")
	}

	if strings.Contains(output, "timeout") {
		e.Suggestions = append(e.Suggestions, "Increase timeout value: --timeout")
		e.Suggestions = append(e.Suggestions, "Try again later as the server might be temporarily unavailable")
	}

	if strings.Contains(output, "disk space") || strings.Contains(output, "no space") {
		e.Suggestions = append(e.Suggestions, "Free up disk space")
		e.Suggestions = append(e.Suggestions, "Clean pip cache: pip cache purge")
	}

	if strings.Contains(output, "requirement already satisfied") {
		e.Suggestions = append(e.Suggestions, "Use --upgrade flag to upgrade to the latest version")
		e.Suggestions = append(e.Suggestions, "Use --force-reinstall to reinstall the package")
	}
}

// Common error instances
var (
	ErrPipNotInstalled = NewPipError(ErrorTypePipNotInstalled, "pip is not installed").
				WithSuggestion("Install pip using your system package manager or download get-pip.py")

	ErrPythonNotFound = NewPipError(ErrorTypePythonNotFound, "Python interpreter not found").
				WithSuggestion("Install Python from https://python.org or use your system package manager")

	ErrPackageNotFound = NewPipError(ErrorTypePackageNotFound, "package not found").
				WithSuggestion("Check the package name spelling and availability on PyPI")

	ErrVenvNotFound = NewPipError(ErrorTypeVenvNotFound, "virtual environment not found").
			WithSuggestion("Create a virtual environment first using CreateVenv()")

	ErrVenvAlreadyExists = NewPipError(ErrorTypeVenvAlreadyExists, "virtual environment already exists").
				WithSuggestion("Use a different path or remove the existing environment")

	ErrInvalidPackageSpec = NewPipError(ErrorTypeInvalidPackageSpec, "invalid package specification").
				WithSuggestion("Provide a valid package name")

	ErrPermissionDenied = NewPipError(ErrorTypePermissionDenied, "permission denied").
				WithSuggestion("Run with elevated privileges or use a virtual environment")

	ErrNetworkError = NewPipError(ErrorTypeNetworkError, "network error").
			WithSuggestion("Check your internet connection and try again")

	ErrUnsupportedOS = NewPipError(ErrorTypeUnsupportedOS, "unsupported operating system").
				WithSuggestion("This operation is not supported on your operating system")

	ErrFeatureDisabled = NewPipError(ErrorTypeFeatureDisabled, "feature is disabled").
				WithSuggestion("This feature has been disabled or is not available")
)

// IsErrorType checks if an error is of a specific type
func IsErrorType(err error, errorType ErrorType) bool {
	if pipErr, ok := err.(*PipErrorDetails); ok {
		return pipErr.Type == errorType
	}
	return false
}

// GetErrorType returns the error type if it's a PipErrorDetails, otherwise returns empty string
func GetErrorType(err error) ErrorType {
	if pipErr, ok := err.(*PipErrorDetails); ok {
		return pipErr.Type
	}
	return ""
}

// WrapError wraps a generic error into a PipErrorDetails
func WrapError(err error, errorType ErrorType, message string) *PipErrorDetails {
	return &PipErrorDetails{
		Type:    errorType,
		Message: message,
		Cause:   err,
		Context: make(map[string]string),
	}
}

// ErrorHandler provides centralized error handling and logging
type ErrorHandler struct {
	logger *Logger
}

// NewErrorHandler creates a new error handler
func NewErrorHandler(logger *Logger) *ErrorHandler {
	return &ErrorHandler{
		logger: logger,
	}
}

// Handle processes an error and logs it appropriately
func (eh *ErrorHandler) Handle(err error, context string) error {
	if err == nil {
		return nil
	}

	if pipErr, ok := err.(*PipErrorDetails); ok {
		eh.handlePipError(pipErr, context)
		return pipErr
	}

	// Wrap generic errors
	wrappedErr := WrapError(err, ErrorTypeCommandFailed, err.Error())
	eh.handlePipError(wrappedErr, context)
	return wrappedErr
}

// handlePipError handles PipErrorDetails specifically
func (eh *ErrorHandler) handlePipError(err *PipErrorDetails, context string) {
	if eh.logger == nil {
		return
	}

	// Log based on error type
	switch err.Type {
	case ErrorTypePipNotInstalled, ErrorTypePythonNotFound:
		eh.logger.Error("System dependency missing in %s: %s", context, err.Message)
	case ErrorTypePermissionDenied:
		eh.logger.Error("Permission error in %s: %s", context, err.Message)
	case ErrorTypeNetworkError:
		eh.logger.Warn("Network issue in %s: %s", context, err.Message)
	case ErrorTypeTimeout:
		eh.logger.Warn("Timeout in %s: %s", context, err.Message)
	default:
		eh.logger.Error("Error in %s: %s", context, err.Message)
	}

	// Log suggestions if available
	if len(err.Suggestions) > 0 {
		eh.logger.Info("Suggestions: %s", strings.Join(err.Suggestions, "; "))
	}

	// Log additional context
	if err.Command != "" {
		eh.logger.Debug("Failed command: %s", err.Command)
	}

	if err.Output != "" {
		eh.logger.Debug("Command output: %s", err.Output)
	}
}
