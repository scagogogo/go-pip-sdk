package pip

import (
	"bytes"
	"context"
	"errors"
	"log"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestNewManager(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		want   bool // whether manager should be created successfully
	}{
		{
			name:   "with nil config",
			config: nil,
			want:   true,
		},
		{
			name: "with custom config",
			config: &Config{
				Timeout:  60 * time.Second,
				LogLevel: "DEBUG",
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewManager(tt.config)
			if (manager != nil) != tt.want {
				t.Errorf("NewManager() = %v, want %v", manager != nil, tt.want)
			}

			if manager != nil {
				if manager.config == nil {
					t.Error("Manager config should not be nil")
				}
				if manager.ctx == nil {
					t.Error("Manager context should not be nil")
				}
			}
		})
	}
}

func TestNewManagerWithContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	manager := NewManagerWithContext(ctx, nil)
	if manager == nil {
		t.Fatal("NewManagerWithContext() returned nil")
	}

	if manager.ctx != ctx {
		t.Error("Manager context should match provided context")
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	if config == nil {
		t.Fatal("DefaultConfig() returned nil")
	}

	if config.Timeout != 30*time.Second {
		t.Errorf("Default timeout = %v, want %v", config.Timeout, 30*time.Second)
	}

	if config.Retries != 3 {
		t.Errorf("Default retries = %d, want %d", config.Retries, 3)
	}

	if config.LogLevel != "INFO" {
		t.Errorf("Default log level = %s, want %s", config.LogLevel, "INFO")
	}

	if config.Environment == nil {
		t.Error("Default environment should not be nil")
	}

	if config.ExtraOptions == nil {
		t.Error("Default extra options should not be nil")
	}
}

func TestManagerSetters(t *testing.T) {
	manager := NewManager(nil)

	// Test SetConfig
	newConfig := &Config{
		Timeout:  60 * time.Second,
		LogLevel: "DEBUG",
	}
	manager.SetConfig(newConfig)

	if manager.GetConfig() != newConfig {
		t.Error("SetConfig/GetConfig failed")
	}

	// Test SetContext
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	manager.SetContext(ctx)
	if manager.ctx != ctx {
		t.Error("SetContext failed")
	}
}

func TestValidatePackageSpec(t *testing.T) {
	manager := NewManager(nil)

	tests := []struct {
		name    string
		pkg     *PackageSpec
		wantErr bool
	}{
		{
			name:    "nil package",
			pkg:     nil,
			wantErr: true,
		},
		{
			name: "empty name",
			pkg: &PackageSpec{
				Name: "",
			},
			wantErr: true,
		},
		{
			name: "valid package",
			pkg: &PackageSpec{
				Name:    "requests",
				Version: ">=2.0.0",
			},
			wantErr: false,
		},
		{
			name: "package with extras",
			pkg: &PackageSpec{
				Name:   "requests",
				Extras: []string{"security", "socks"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.validatePackageSpec(tt.pkg)
			if (err != nil) != tt.wantErr {
				t.Errorf("validatePackageSpec() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreatePipError(t *testing.T) {
	manager := NewManager(nil)

	command := "pip install nonexistent-package"
	output := "ERROR: Could not find a version that satisfies the requirement"
	exitCode := 1

	pipErr := manager.createPipError(command, output, exitCode, nil)

	if pipErr == nil {
		t.Fatal("createPipError() returned nil")
	}

	if pipErr.Type != "command_failed" {
		t.Errorf("Error type = %s, want %s", pipErr.Type, "command_failed")
	}

	if pipErr.Command != command {
		t.Errorf("Error command = %s, want %s", pipErr.Command, command)
	}

	if pipErr.Output != output {
		t.Errorf("Error output = %s, want %s", pipErr.Output, output)
	}

	if pipErr.ExitCode != exitCode {
		t.Errorf("Error exit code = %d, want %d", pipErr.ExitCode, exitCode)
	}
}

// Integration tests (require actual pip installation)
func TestIntegrationIsInstalled(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	manager := NewManager(nil)

	// This test will check if pip is actually installed on the system
	installed, err := manager.IsInstalled()

	// We don't assert the result since pip might not be installed
	// but we check that the method doesn't panic
	t.Logf("Pip installed: %v, error: %v", installed, err)
}

func TestIntegrationGetVersion(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	manager := NewManager(nil)

	// Check if pip is installed first
	installed, err := manager.IsInstalled()
	if err != nil || !installed {
		t.Skip("Pip not installed, skipping version test")
	}

	version, err := manager.GetVersion()
	if err != nil {
		t.Errorf("GetVersion() error = %v", err)
	}

	if version == "" {
		t.Error("GetVersion() returned empty version")
	}

	t.Logf("Pip version: %s", version)
}

// Benchmark tests
func BenchmarkNewManager(b *testing.B) {
	for i := 0; i < b.N; i++ {
		manager := NewManager(nil)
		_ = manager
	}
}

func BenchmarkValidatePackageSpec(b *testing.B) {
	manager := NewManager(nil)
	pkg := &PackageSpec{
		Name:    "requests",
		Version: ">=2.0.0",
		Extras:  []string{"security"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = manager.validatePackageSpec(pkg)
	}
}

// Helper functions for tests
func setupTestManager(t *testing.T) *Manager {
	config := &Config{
		Timeout:     10 * time.Second,
		LogLevel:    "DEBUG",
		Environment: make(map[string]string),
	}

	return NewManager(config)
}

func skipIfNoPip(t *testing.T, manager *Manager) {
	installed, err := manager.IsInstalled()
	if err != nil || !installed {
		t.Skip("Pip not installed, skipping test")
	}
}

// Test cleanup
func TestMain(m *testing.M) {
	// Setup

	// Run tests
	code := m.Run()

	// Cleanup

	os.Exit(code)
}

func TestGetOSType(t *testing.T) {
	osType := GetOSType()

	// Test that we get a valid OS type
	validTypes := []OSType{OSWindows, OSMacOS, OSLinux, OSUnknown}
	found := false
	for _, validType := range validTypes {
		if osType == validType {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("GetOSType() returned invalid type: %s", osType)
	}

	// Test that the result is consistent with runtime.GOOS
	switch runtime.GOOS {
	case "windows":
		if osType != OSWindows {
			t.Errorf("Expected OSWindows on Windows, got %s", osType)
		}
	case "darwin":
		if osType != OSMacOS {
			t.Errorf("Expected OSMacOS on macOS, got %s", osType)
		}
	case "linux":
		if osType != OSLinux {
			t.Errorf("Expected OSLinux on Linux, got %s", osType)
		}
	default:
		if osType != OSUnknown {
			t.Errorf("Expected OSUnknown on %s, got %s", runtime.GOOS, osType)
		}
	}

	// Test OSType string representation
	tests := []struct {
		osType OSType
		want   string
	}{
		{OSWindows, "windows"},
		{OSMacOS, "darwin"}, // OSType returns the actual runtime.GOOS value
		{OSLinux, "linux"},
		{OSUnknown, "unknown"},
		{OSType("invalid"), "invalid"},
	}

	for _, tt := range tests {
		t.Run(string(tt.osType), func(t *testing.T) {
			got := string(tt.osType)
			if got != tt.want {
				t.Errorf("OSType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManagerWithCustomLogger(t *testing.T) {
	var buf bytes.Buffer

	// Create custom logger
	loggerConfig := &LoggerConfig{
		Level:  LogLevelDebug,
		Output: &buf,
		Prefix: "[test]",
	}

	logger, err := NewLogger(loggerConfig)
	if err != nil {
		t.Fatalf("NewLogger() error: %v", err)
	}
	defer logger.Close()

	// Create manager and set custom logger
	manager := NewManager(nil)
	manager.SetCustomLogger(logger)

	// Test logging through manager
	manager.logInfo("test info message")
	manager.logError("test error message")
	manager.logDebug("test debug message")
	manager.logWarn("test warn message")

	output := buf.String()

	// Verify messages were logged
	expectedMessages := []string{
		"test info message",
		"test error message",
		"test debug message",
		"test warn message",
	}

	for _, expected := range expectedMessages {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected message '%s' not found in output", expected)
		}
	}
}

func TestManagerErrorCreation(t *testing.T) {
	manager := NewManager(nil)

	command := "pip install nonexistent"
	output := "ERROR: Could not find a version"
	exitCode := 1
	cause := errors.New("command failed")

	pipErr := manager.createPipError(command, output, exitCode, cause)

	if pipErr == nil {
		t.Fatal("createPipError() returned nil")
	}

	// Verify error details
	if pipErr.Type != "command_failed" {
		t.Errorf("Error type = %s, want command_failed", pipErr.Type)
	}

	if pipErr.Command != command {
		t.Errorf("Error command = %s, want %s", pipErr.Command, command)
	}

	if pipErr.Output != output {
		t.Errorf("Error output = %s, want %s", pipErr.Output, output)
	}

	if pipErr.ExitCode != exitCode {
		t.Errorf("Error exit code = %d, want %d", pipErr.ExitCode, exitCode)
	}

	// Note: PipError doesn't have a Cause field, unlike PipErrorDetails

	// Test error message format
	errorMsg := pipErr.Error()
	if !strings.Contains(errorMsg, "command failed") {
		t.Error("Error message should contain error description")
	}
}

func TestManagerConfigurationEdgeCases(t *testing.T) {
	// Test with nil config
	manager := NewManager(nil)
	if manager.config == nil {
		t.Error("Manager should have default config when nil is provided")
	}

	// Test setting nil config - current implementation allows nil
	manager.SetConfig(nil)
	if manager.config != nil {
		t.Error("SetConfig(nil) should set config to nil (current implementation)")
	}

	// Test with empty config
	emptyConfig := &Config{}
	manager.SetConfig(emptyConfig)
	config := manager.GetConfig()
	if config != emptyConfig {
		t.Error("SetConfig/GetConfig should work with empty config")
	}
}

func TestManagerContextHandling(t *testing.T) {
	manager := NewManager(nil)

	// Test with timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	manager.SetContext(ctx)
	if manager.ctx != ctx {
		t.Error("SetContext should update manager context")
	}

	// Test with cancelled context
	cancelCtx, cancelFunc := context.WithCancel(context.Background())
	cancelFunc() // Cancel immediately

	manager.SetContext(cancelCtx)
	if manager.ctx != cancelCtx {
		t.Error("SetContext should work with cancelled context")
	}
}

func TestValidatePackageSpecEdgeCases(t *testing.T) {
	manager := NewManager(nil)

	tests := []struct {
		name    string
		pkg     *PackageSpec
		wantErr bool
	}{
		{
			name:    "nil package",
			pkg:     nil,
			wantErr: true,
		},
		{
			name: "empty name",
			pkg: &PackageSpec{
				Name: "",
			},
			wantErr: true,
		},
		{
			name: "whitespace only name",
			pkg: &PackageSpec{
				Name: "   ",
			},
			wantErr: false, // Current implementation doesn't trim whitespace
		},
		{
			name: "valid minimal package",
			pkg: &PackageSpec{
				Name: "a",
			},
			wantErr: false,
		},
		{
			name: "package with all options",
			pkg: &PackageSpec{
				Name:           "requests",
				Version:        ">=2.25.0",
				Extras:         []string{"security", "socks"},
				Upgrade:        true,
				ForceReinstall: true,
				Editable:       true,
				Index:          "https://pypi.org/simple/",
				Options: map[string]string{
					"timeout":      "30",
					"retries":      "3",
					"trusted-host": "pypi.org",
				},
			},
			wantErr: false,
		},
		{
			name: "package with empty extras",
			pkg: &PackageSpec{
				Name:   "requests",
				Extras: []string{},
			},
			wantErr: false,
		},
		{
			name: "package with nil options",
			pkg: &PackageSpec{
				Name:    "requests",
				Options: nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.validatePackageSpec(tt.pkg)

			if tt.wantErr {
				if err == nil {
					t.Errorf("validatePackageSpec() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("validatePackageSpec() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestSetLogger(t *testing.T) {
	manager := NewManager(nil)

	// Create a standard library logger
	var buf bytes.Buffer
	stdLogger := log.New(&buf, "[test] ", log.LstdFlags)

	// Test SetLogger
	manager.SetLogger(stdLogger)

	// The manager's logger field should be set
	if manager.logger != stdLogger {
		t.Error("SetLogger() should set the logger field")
	}

	// Test with nil logger
	manager.SetLogger(nil)
	if manager.logger != nil {
		t.Error("SetLogger(nil) should set logger to nil")
	}
}

func TestManagerLogWarn(t *testing.T) {
	var buf bytes.Buffer

	// Create custom logger
	loggerConfig := &LoggerConfig{
		Level:  LogLevelWarn,
		Output: &buf,
		Prefix: "[test]",
	}

	logger, err := NewLogger(loggerConfig)
	if err != nil {
		t.Fatalf("NewLogger() error: %v", err)
	}
	defer logger.Close()

	// Create manager and set custom logger
	manager := NewManager(nil)
	manager.SetCustomLogger(logger)

	// Test logWarn method
	manager.logWarn("test warning message")

	output := buf.String()
	if !strings.Contains(output, "test warning message") {
		t.Error("logWarn() should log warning messages")
	}

	if !strings.Contains(output, "[WARN]") {
		t.Error("logWarn() should include WARN level indicator")
	}
}
