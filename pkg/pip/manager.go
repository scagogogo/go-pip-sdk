package pip

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

// Manager implements the PipManager interface
type Manager struct {
	config       *Config
	logger       *log.Logger
	customLogger *Logger
	ctx          context.Context
}

// NewManager creates a new pip manager instance
func NewManager(config *Config) *Manager {
	if config == nil {
		config = DefaultConfig()
	}

	logger := log.New(os.Stdout, "[pip-sdk] ", log.LstdFlags)

	return &Manager{
		config: config,
		logger: logger,
		ctx:    context.Background(),
	}
}

// NewManagerWithContext creates a new pip manager instance with context
func NewManagerWithContext(ctx context.Context, config *Config) *Manager {
	manager := NewManager(config)
	manager.ctx = ctx
	return manager
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		Timeout:      30 * time.Second,
		Retries:      3,
		LogLevel:     "INFO",
		CacheDir:     "",
		Environment:  make(map[string]string),
		ExtraOptions: make(map[string]string),
	}
}

// GetOSType detects the current operating system
func GetOSType() OSType {
	switch runtime.GOOS {
	case "windows":
		return OSWindows
	case "darwin":
		return OSMacOS
	case "linux":
		return OSLinux
	default:
		return OSUnknown
	}
}

// SetConfig updates the manager configuration
func (m *Manager) SetConfig(config *Config) {
	m.config = config
}

// GetConfig returns the current configuration
func (m *Manager) GetConfig() *Config {
	return m.config
}

// SetLogger sets a custom logger
func (m *Manager) SetLogger(logger *log.Logger) {
	m.logger = logger
}

// SetContext sets the context for operations
func (m *Manager) SetContext(ctx context.Context) {
	m.ctx = ctx
}

// validatePackageSpec validates a package specification
func (m *Manager) validatePackageSpec(pkg *PackageSpec) error {
	if pkg == nil {
		return ErrInvalidPackageSpec
	}

	if pkg.Name == "" {
		return &PipError{
			Type:    "invalid_package_spec",
			Message: "package name cannot be empty",
		}
	}

	return nil
}

// createPipError creates a pip error from command execution
func (m *Manager) createPipError(command string, output string, exitCode int, err error) *PipError {
	pipErr := &PipError{
		Type:     "command_failed",
		Command:  command,
		Output:   output,
		ExitCode: exitCode,
	}

	if err != nil {
		pipErr.Message = fmt.Sprintf("pip command failed: %v", err)
	} else {
		pipErr.Message = fmt.Sprintf("pip command failed with exit code %d", exitCode)
	}

	return pipErr
}
