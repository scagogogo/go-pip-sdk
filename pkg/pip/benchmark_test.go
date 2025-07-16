package pip

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// BenchmarkManagerCreation benchmarks the creation of pip managers
func BenchmarkManagerCreation(b *testing.B) {
	config := &Config{
		Timeout: 30 * time.Second,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewManager(config)
	}
}

// BenchmarkManagerCreationWithContext benchmarks manager creation with context
func BenchmarkManagerCreationWithContext(b *testing.B) {
	config := &Config{
		Timeout: 30 * time.Second,
	}
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewManagerWithContext(ctx, config)
	}
}

// BenchmarkIsInstalled benchmarks the pip installation check
func BenchmarkIsInstalled(b *testing.B) {
	manager := NewManager(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := manager.IsInstalled()
		if err != nil {
			b.Fatalf("IsInstalled failed: %v", err)
		}
	}
}

// BenchmarkGetVersion benchmarks getting pip version
func BenchmarkGetVersion(b *testing.B) {
	manager := NewManager(nil)

	// Skip if pip is not installed
	installed, err := manager.IsInstalled()
	if err != nil || !installed {
		b.Skip("pip not installed")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := manager.GetVersion()
		if err != nil {
			b.Fatalf("GetVersion failed: %v", err)
		}
	}
}

// BenchmarkListPackages benchmarks listing installed packages
func BenchmarkListPackages(b *testing.B) {
	manager := NewManager(nil)

	// Skip if pip is not installed
	installed, err := manager.IsInstalled()
	if err != nil || !installed {
		b.Skip("pip not installed")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := manager.ListPackages()
		if err != nil {
			b.Fatalf("ListPackages failed: %v", err)
		}
	}
}

// BenchmarkPackageSpecCreation benchmarks creating package specifications
func BenchmarkPackageSpecCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = &PackageSpec{
			Name:    "requests",
			Version: ">=2.25.0",
			Extras:  []string{"security", "socks"},
			Options: map[string]string{
				"no-cache-dir": "",
				"timeout":      "60",
			},
		}
	}
}

// BenchmarkVirtualEnvironmentOperations benchmarks virtual environment operations
func BenchmarkVirtualEnvironmentOperations(b *testing.B) {
	if testing.Short() {
		b.Skip("skipping benchmark in short mode")
	}

	manager := NewManager(nil)

	// Skip if pip is not installed
	installed, err := manager.IsInstalled()
	if err != nil || !installed {
		b.Skip("pip not installed")
	}

	b.Run("CreateVenv", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tempDir, err := os.MkdirTemp("", "bench-venv-*")
			if err != nil {
				b.Fatalf("Failed to create temp dir: %v", err)
			}

			venvPath := filepath.Join(tempDir, "venv")

			b.StartTimer()
			err = manager.CreateVenv(venvPath)
			b.StopTimer()

			if err != nil {
				b.Fatalf("CreateVenv failed: %v", err)
			}

			// Clean up
			os.RemoveAll(tempDir)
		}
	})

	b.Run("VenvInfo", func(b *testing.B) {
		// Create a test virtual environment
		tempDir, err := os.MkdirTemp("", "bench-venv-info-*")
		if err != nil {
			b.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		venvPath := filepath.Join(tempDir, "venv")
		err = manager.CreateVenv(venvPath)
		if err != nil {
			b.Fatalf("Failed to create test venv: %v", err)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := manager.GetVenvInfo(venvPath)
			if err != nil {
				b.Fatalf("GetVenvInfo failed: %v", err)
			}
		}
	})
}

// BenchmarkErrorHandling benchmarks error creation and handling
func BenchmarkErrorHandling(b *testing.B) {
	b.Run("NewPipError", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = NewPipError(ErrorTypePackageNotFound, "package not found")
		}
	})

	b.Run("ErrorWithSuggestion", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			err := NewPipError(ErrorTypePackageNotFound, "package not found")
			_ = err.WithSuggestion("Check package name")
		}
	})

	b.Run("IsErrorType", func(b *testing.B) {
		err := NewPipError(ErrorTypePackageNotFound, "package not found")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = IsErrorType(err, ErrorTypePackageNotFound)
		}
	})
}

// BenchmarkConfigurationParsing benchmarks configuration parsing
func BenchmarkConfigurationParsing(b *testing.B) {
	b.Run("DefaultConfig", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			config := &Config{
				Timeout:      30 * time.Second,
				Retries:      3,
				DefaultIndex: "https://pypi.org/simple/",
				TrustedHosts: []string{"pypi.org", "pypi.python.org"},
				LogLevel:     "INFO",
				ExtraOptions: map[string]string{
					"no-cache-dir": "",
				},
				Environment: map[string]string{
					"PIP_DISABLE_PIP_VERSION_CHECK": "1",
				},
			}
			_ = config
		}
	})
}

// BenchmarkLoggerOperations benchmarks logger operations
func BenchmarkLoggerOperations(b *testing.B) {
	config := &LoggerConfig{
		Level:  LogLevelInfo,
		Output: os.Stdout,
		Prefix: "[bench]",
	}

	logger, err := NewLogger(config)
	if err != nil {
		b.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	b.Run("LogInfo", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			logger.Info("Benchmark log message %d", i)
		}
	})

	b.Run("LogDebug", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			logger.Debug("Benchmark debug message %d", i)
		}
	})

	b.Run("IsEnabled", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = logger.IsEnabled(LogLevelDebug)
		}
	})
}

// BenchmarkConcurrentOperations benchmarks concurrent operations
func BenchmarkConcurrentOperations(b *testing.B) {
	manager := NewManager(nil)

	// Skip if pip is not installed
	installed, err := manager.IsInstalled()
	if err != nil || !installed {
		b.Skip("pip not installed")
	}

	b.Run("ConcurrentIsInstalled", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, err := manager.IsInstalled()
				if err != nil {
					b.Errorf("IsInstalled failed: %v", err)
				}
			}
		})
	})

	b.Run("ConcurrentGetVersion", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, err := manager.GetVersion()
				if err != nil {
					b.Errorf("GetVersion failed: %v", err)
				}
			}
		})
	})
}
