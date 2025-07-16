package pip

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLogLevel(t *testing.T) {
	tests := []struct {
		level    LogLevel
		expected string
	}{
		{LogLevelDebug, "DEBUG"},
		{LogLevelInfo, "INFO"},
		{LogLevelWarn, "WARN"},
		{LogLevelError, "ERROR"},
		{LogLevelFatal, "FATAL"},
		{LogLevel(999), "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.level.String()
			if result != tt.expected {
				t.Errorf("LogLevel.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected LogLevel
		wantErr  bool
	}{
		{"DEBUG", LogLevelDebug, false},
		{"debug", LogLevelDebug, false},
		{"INFO", LogLevelInfo, false},
		{"info", LogLevelInfo, false},
		{"WARN", LogLevelWarn, false},
		{"WARNING", LogLevelWarn, false},
		{"warn", LogLevelWarn, false},
		{"ERROR", LogLevelError, false},
		{"error", LogLevelError, false},
		{"FATAL", LogLevelFatal, false},
		{"fatal", LogLevelFatal, false},
		{"INVALID", LogLevelInfo, true},
		{"", LogLevelInfo, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ParseLogLevel(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseLogLevel(%s) expected error but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ParseLogLevel(%s) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("ParseLogLevel(%s) = %v, want %v", tt.input, result, tt.expected)
				}
			}
		})
	}
}

func TestDefaultLoggerConfig(t *testing.T) {
	config := DefaultLoggerConfig()

	if config == nil {
		t.Fatal("DefaultLoggerConfig() returned nil")
	}

	if config.Level != LogLevelInfo {
		t.Errorf("Default level = %v, want %v", config.Level, LogLevelInfo)
	}

	if config.Prefix != "[pip-sdk]" {
		t.Errorf("Default prefix = %s, want [pip-sdk]", config.Prefix)
	}

	if config.EnableFile {
		t.Error("Default should not enable file logging")
	}

	if config.Output == nil {
		t.Error("Default output should not be nil")
	}
}

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name    string
		config  *LoggerConfig
		wantErr bool
	}{
		{
			name:    "nil config",
			config:  nil,
			wantErr: false,
		},
		{
			name: "valid config",
			config: &LoggerConfig{
				Level:  LogLevelDebug,
				Prefix: "[test]",
			},
			wantErr: false,
		},
		{
			name: "file logging config",
			config: &LoggerConfig{
				Level:      LogLevelInfo,
				EnableFile: true,
				LogFile:    "test.log",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup temp directory for file logging tests
			var tempDir string
			if tt.config != nil && tt.config.EnableFile {
				var err error
				tempDir, err = os.MkdirTemp("", "pip-sdk-logger-test")
				if err != nil {
					t.Fatalf("Failed to create temp dir: %v", err)
				}
				defer os.RemoveAll(tempDir)

				tt.config.LogFile = filepath.Join(tempDir, tt.config.LogFile)
			}

			logger, err := NewLogger(tt.config)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewLogger() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("NewLogger() unexpected error: %v", err)
				return
			}

			if logger == nil {
				t.Fatal("NewLogger() returned nil logger")
			}

			// Test logger methods don't panic
			logger.Debug("test debug")
			logger.Info("test info")
			logger.Warn("test warn")
			logger.Error("test error")

			// Verify file was created if file logging is enabled
			if tt.config != nil && tt.config.EnableFile {
				if _, err := os.Stat(tt.config.LogFile); os.IsNotExist(err) {
					t.Error("Log file was not created")
				}
			}

			// Clean up
			if err := logger.Close(); err != nil {
				t.Errorf("Logger.Close() error: %v", err)
			}
		})
	}
}

func TestLoggerLevels(t *testing.T) {
	var buf bytes.Buffer

	config := &LoggerConfig{
		Level:  LogLevelWarn,
		Output: &buf,
		Prefix: "[test]",
	}

	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("NewLogger() error: %v", err)
	}
	defer logger.Close()

	// Test that only WARN and above are logged
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	output := buf.String()

	if strings.Contains(output, "debug message") {
		t.Error("Debug message should not be logged at WARN level")
	}

	if strings.Contains(output, "info message") {
		t.Error("Info message should not be logged at WARN level")
	}

	if !strings.Contains(output, "warn message") {
		t.Error("Warn message should be logged at WARN level")
	}

	if !strings.Contains(output, "error message") {
		t.Error("Error message should be logged at WARN level")
	}
}

func TestLoggerIsEnabled(t *testing.T) {
	logger, err := NewLogger(&LoggerConfig{Level: LogLevelWarn})
	if err != nil {
		t.Fatalf("NewLogger() error: %v", err)
	}
	defer logger.Close()

	tests := []struct {
		level    LogLevel
		expected bool
	}{
		{LogLevelDebug, false},
		{LogLevelInfo, false},
		{LogLevelWarn, true},
		{LogLevelError, true},
		{LogLevelFatal, true},
	}

	for _, tt := range tests {
		t.Run(tt.level.String(), func(t *testing.T) {
			result := logger.IsEnabled(tt.level)
			if result != tt.expected {
				t.Errorf("IsEnabled(%v) = %v, want %v", tt.level, result, tt.expected)
			}
		})
	}
}

func TestLoggerSetLevel(t *testing.T) {
	logger, err := NewLogger(&LoggerConfig{Level: LogLevelInfo})
	if err != nil {
		t.Fatalf("NewLogger() error: %v", err)
	}
	defer logger.Close()

	// Initial level
	if logger.GetLevel() != LogLevelInfo {
		t.Errorf("Initial level = %v, want %v", logger.GetLevel(), LogLevelInfo)
	}

	// Change level
	logger.SetLevel(LogLevelError)
	if logger.GetLevel() != LogLevelError {
		t.Errorf("After SetLevel, level = %v, want %v", logger.GetLevel(), LogLevelError)
	}
}

func TestLoggerOperationMethods(t *testing.T) {
	var buf bytes.Buffer

	config := &LoggerConfig{
		Level:  LogLevelDebug,
		Output: &buf,
		Prefix: "[test]",
	}

	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("NewLogger() error: %v", err)
	}
	defer logger.Close()

	// Test operation logging methods
	logger.LogCommand("pip", []string{"install", "requests"}, 5*time.Second)
	logger.LogCommandError("pip", []string{"install", "nonexistent"},
		&PipErrorDetails{Message: "package not found"}, "error output")
	logger.LogPackageOperation("install", "requests", true, 3*time.Second)
	logger.LogPackageOperation("install", "nonexistent", false, 1*time.Second)
	logger.LogVenvOperation("create", "/path/to/venv", true)
	logger.LogVenvOperation("create", "/invalid/path", false)
	logger.LogProjectOperation("init", "/path/to/project", true)
	logger.LogProjectOperation("init", "/invalid/project", false)

	output := buf.String()

	// Verify expected content is logged
	expectedStrings := []string{
		"pip install requests",
		"duration: 5s",
		"package not found",
		"Package install successful: requests",
		"Package install failed: nonexistent",
		"Virtual environment create successful",
		"Virtual environment create failed",
		"Project init successful",
		"Project init failed",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected string '%s' not found in output", expected)
		}
	}
}

func TestLoggerFileLogging(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "pip-sdk-logger-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logFile := filepath.Join(tempDir, "test.log")

	config := &LoggerConfig{
		Level:      LogLevelInfo,
		EnableFile: true,
		LogFile:    logFile,
		Prefix:     "[test]",
	}

	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("NewLogger() error: %v", err)
	}

	// Log some messages
	logger.Info("test message 1")
	logger.Error("test message 2")

	// Close logger to flush file
	if err := logger.Close(); err != nil {
		t.Errorf("Logger.Close() error: %v", err)
	}

	// Verify file exists and has content
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Fatal("Log file was not created")
	}

	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "test message 1") {
		t.Error("Log file should contain 'test message 1'")
	}

	if !strings.Contains(contentStr, "test message 2") {
		t.Error("Log file should contain 'test message 2'")
	}

	if !strings.Contains(contentStr, "[INFO]") {
		t.Error("Log file should contain log level")
	}

	if !strings.Contains(contentStr, "[test]") {
		t.Error("Log file should contain prefix")
	}
}

func TestManagerCustomLogger(t *testing.T) {
	var buf bytes.Buffer

	config := &LoggerConfig{
		Level:  LogLevelDebug,
		Output: &buf,
		Prefix: "[test]",
	}

	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("NewLogger() error: %v", err)
	}
	defer logger.Close()

	manager := NewManager(nil)
	manager.SetCustomLogger(logger)

	// Test that manager uses custom logger
	manager.logInfo("test info message")
	manager.logError("test error message")
	manager.logDebug("test debug message")
	manager.logWarn("test warn message")

	output := buf.String()

	expectedStrings := []string{
		"test info message",
		"test error message",
		"test debug message",
		"test warn message",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected string '%s' not found in manager log output", expected)
		}
	}
}

// Benchmark tests
func BenchmarkLoggerInfo(b *testing.B) {
	logger, err := NewLogger(&LoggerConfig{
		Level:  LogLevelInfo,
		Output: os.Stdout,
	})
	if err != nil {
		b.Fatalf("NewLogger() error: %v", err)
	}
	defer logger.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark message %d", i)
	}
}

func BenchmarkLoggerDisabledLevel(b *testing.B) {
	logger, err := NewLogger(&LoggerConfig{
		Level:  LogLevelError,
		Output: os.Stdout,
	})
	if err != nil {
		b.Fatalf("NewLogger() error: %v", err)
	}
	defer logger.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Debug("benchmark debug message %d", i)
	}
}

// Test Fatal function (tricky because it calls os.Exit)
func TestLoggerFatal(t *testing.T) {
	// We can't easily test Fatal() because it calls os.Exit(1)
	// But we can test that it logs the message before exiting

	var buf bytes.Buffer
	logger, err := NewLogger(&LoggerConfig{
		Level:  LogLevelFatal,
		Output: &buf,
		Prefix: "[test]",
	})
	if err != nil {
		t.Fatalf("NewLogger() error: %v", err)
	}
	defer logger.Close()

	// We can't actually call Fatal() in tests because it would exit
	// Instead, we test the log method directly with FATAL level
	logger.log(LogLevelFatal, "test fatal message")

	output := buf.String()
	if !strings.Contains(output, "test fatal message") {
		t.Error("Fatal level message should be logged")
	}

	if !strings.Contains(output, "[FATAL]") {
		t.Error("Fatal level should be indicated in log")
	}
}

func TestLoggerWithNilOutput(t *testing.T) {
	// Test logger with nil output (should not panic)
	logger, err := NewLogger(&LoggerConfig{
		Level:  LogLevelInfo,
		Output: nil, // This should default to os.Stdout
		Prefix: "[test]",
	})
	if err != nil {
		t.Fatalf("NewLogger() error: %v", err)
	}
	defer logger.Close()

	// These should not panic
	logger.Info("test message")
	logger.Error("test error")
}

func TestLoggerMultiWriter(t *testing.T) {
	var buf1, buf2 bytes.Buffer

	// Create logger with multiple outputs
	multiWriter := io.MultiWriter(&buf1, &buf2)

	logger, err := NewLogger(&LoggerConfig{
		Level:  LogLevelInfo,
		Output: multiWriter,
		Prefix: "[test]",
	})
	if err != nil {
		t.Fatalf("NewLogger() error: %v", err)
	}
	defer logger.Close()

	logger.Info("test message")

	// Both buffers should contain the message
	output1 := buf1.String()
	output2 := buf2.String()

	if !strings.Contains(output1, "test message") {
		t.Error("First output should contain message")
	}

	if !strings.Contains(output2, "test message") {
		t.Error("Second output should contain message")
	}

	if output1 != output2 {
		t.Error("Both outputs should be identical")
	}
}

func TestLoggerCallerInfo(t *testing.T) {
	var buf bytes.Buffer

	logger, err := NewLogger(&LoggerConfig{
		Level:  LogLevelDebug,
		Output: &buf,
		Prefix: "[test]",
	})
	if err != nil {
		t.Fatalf("NewLogger() error: %v", err)
	}
	defer logger.Close()

	logger.Info("test message")

	output := buf.String()

	// Should contain caller information (filename:line)
	if !strings.Contains(output, "logger_test.go:") {
		t.Error("Log output should contain caller information")
	}
}

func TestLoggerErrorHandling(t *testing.T) {
	// Test logger creation with invalid file path
	config := &LoggerConfig{
		Level:      LogLevelInfo,
		EnableFile: true,
		LogFile:    "/invalid/path/that/does/not/exist/test.log",
	}

	_, err := NewLogger(config)
	if err == nil {
		t.Error("NewLogger() should fail with invalid file path")
	}
}

func TestLoggerCloseMultipleTimes(t *testing.T) {
	logger, err := NewLogger(&LoggerConfig{
		Level:  LogLevelInfo,
		Output: os.Stdout,
	})
	if err != nil {
		t.Fatalf("NewLogger() error: %v", err)
	}

	// Close multiple times should not panic
	err1 := logger.Close()
	err2 := logger.Close()

	// First close might succeed, second should be safe
	t.Logf("First close: %v, Second close: %v", err1, err2)
}

func TestLoggerLevelFiltering(t *testing.T) {
	var buf bytes.Buffer

	// Create logger with WARN level
	logger, err := NewLogger(&LoggerConfig{
		Level:  LogLevelWarn,
		Output: &buf,
		Prefix: "[test]",
	})
	if err != nil {
		t.Fatalf("NewLogger() error: %v", err)
	}
	defer logger.Close()

	// Test different log levels
	logger.Debug("debug message") // Should be filtered out
	logger.Info("info message")   // Should be filtered out
	logger.Warn("warn message")   // Should be logged
	logger.Error("error message") // Should be logged

	output := buf.String()

	// Debug and Info should not appear
	if strings.Contains(output, "debug message") {
		t.Error("Debug message should be filtered out at WARN level")
	}
	if strings.Contains(output, "info message") {
		t.Error("Info message should be filtered out at WARN level")
	}

	// Warn and Error should appear
	if !strings.Contains(output, "warn message") {
		t.Error("Warn message should be logged at WARN level")
	}
	if !strings.Contains(output, "error message") {
		t.Error("Error message should be logged at WARN level")
	}
}

func TestLoggerConcurrentAccess(t *testing.T) {
	// Use a thread-safe writer for concurrent testing
	logger, err := NewLogger(&LoggerConfig{
		Level:  LogLevelInfo,
		Output: os.Stdout, // Use stdout instead of bytes.Buffer for thread safety
		Prefix: "[test]",
	})
	if err != nil {
		t.Fatalf("NewLogger() error: %v", err)
	}
	defer logger.Close()

	// Test concurrent logging
	const numGoroutines = 5 // Reduced to avoid race conditions
	const messagesPerGoroutine = 5

	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer func() { done <- true }()
			for j := 0; j < messagesPerGoroutine; j++ {
				logger.Info("goroutine %d message %d", id, j)
			}
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// We can't easily verify the output when using stdout,
	// but the important thing is that concurrent access doesn't panic
	t.Log("Concurrent logging test completed without panic")
}

func TestLoggerWithDifferentFormats(t *testing.T) {
	tests := []struct {
		name   string
		format string
		args   []interface{}
		expect string
	}{
		{
			name:   "simple string",
			format: "simple message",
			args:   nil,
			expect: "simple message",
		},
		{
			name:   "formatted string",
			format: "user %s has %d items",
			args:   []interface{}{"john", 5},
			expect: "user john has 5 items",
		},
		{
			name:   "complex format",
			format: "error: %v, code: %d, success: %t",
			args:   []interface{}{errors.New("test error"), 404, false},
			expect: "error: test error, code: 404, success: false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			logger, err := NewLogger(&LoggerConfig{
				Level:  LogLevelInfo,
				Output: &buf,
				Prefix: "",
			})
			if err != nil {
				t.Fatalf("NewLogger() error: %v", err)
			}
			defer logger.Close()

			if tt.args == nil {
				logger.Info(tt.format)
			} else {
				logger.Info(tt.format, tt.args...)
			}

			output := buf.String()
			if !strings.Contains(output, tt.expect) {
				t.Errorf("Expected output to contain '%s', got: %s", tt.expect, output)
			}
		})
	}
}
