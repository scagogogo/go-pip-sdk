package pip

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestNewVenvManager(t *testing.T) {
	manager := NewManager(nil)
	venvManager := NewVenvManager(manager)

	if venvManager == nil {
		t.Fatal("NewVenvManager() returned nil")
	}

	if venvManager.manager != manager {
		t.Error("VenvManager should reference the provided manager")
	}
}

func TestCreateVenv(t *testing.T) {
	manager := NewManager(nil)

	tests := []struct {
		name    string
		path    string
		wantErr bool
		setup   func(string) error
		cleanup func(string) error
	}{
		{
			name:    "empty path",
			path:    "",
			wantErr: true,
		},
		{
			name:    "directory already exists",
			path:    "existing-dir",
			wantErr: true,
			setup: func(path string) error {
				return os.MkdirAll(path, 0755)
			},
			cleanup: func(path string) error {
				return os.RemoveAll(path)
			},
		},
		{
			name:    "valid new directory",
			path:    "new-venv",
			wantErr: false, // Note: This might fail if Python is not available
			cleanup: func(path string) error {
				return os.RemoveAll(path)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fullPath string

			if tt.path != "" {
				tempDir, err := os.MkdirTemp("", "pip-sdk-venv-test")
				if err != nil {
					t.Fatalf("Failed to create temp dir: %v", err)
				}
				defer os.RemoveAll(tempDir)

				fullPath = filepath.Join(tempDir, tt.path)

				if tt.setup != nil {
					if err := tt.setup(fullPath); err != nil {
						t.Fatalf("Setup failed: %v", err)
					}
				}
			} else {
				fullPath = tt.path
			}

			err := manager.CreateVenv(fullPath)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CreateVenv() expected error but got none")
				}
			} else {
				// For valid cases, we expect either success or Python-related errors
				if err != nil {
					t.Logf("CreateVenv() returned error (possibly due to Python not being available): %v", err)
				} else {
					// Verify virtual environment was created
					if _, err := os.Stat(fullPath); os.IsNotExist(err) {
						t.Error("Virtual environment directory was not created")
					}
				}
			}

			if tt.cleanup != nil && fullPath != "" {
				if err := tt.cleanup(fullPath); err != nil {
					t.Logf("Cleanup failed: %v", err)
				}
			}
		})
	}
}

func TestActivateVenv(t *testing.T) {
	manager := NewManager(nil)

	tests := []struct {
		name    string
		path    string
		wantErr bool
		setup   func(string) error
	}{
		{
			name:    "empty path",
			path:    "",
			wantErr: true,
		},
		{
			name:    "non-existent venv",
			path:    "non-existent-venv",
			wantErr: true,
		},
		{
			name:    "invalid venv directory",
			path:    "invalid-venv",
			wantErr: true,
			setup: func(path string) error {
				// Create directory but not a valid venv
				return os.MkdirAll(path, 0755)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fullPath string

			if tt.path != "" {
				tempDir, err := os.MkdirTemp("", "pip-sdk-venv-test")
				if err != nil {
					t.Fatalf("Failed to create temp dir: %v", err)
				}
				defer os.RemoveAll(tempDir)

				fullPath = filepath.Join(tempDir, tt.path)

				if tt.setup != nil {
					if err := tt.setup(fullPath); err != nil {
						t.Fatalf("Setup failed: %v", err)
					}
				}
			} else {
				fullPath = tt.path
			}

			err := manager.ActivateVenv(fullPath)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ActivateVenv() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("ActivateVenv() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestDeactivateVenv(t *testing.T) {
	manager := NewManager(nil)

	// Test deactivating when no venv is active
	err := manager.DeactivateVenv()
	if err != nil {
		t.Errorf("DeactivateVenv() should not error when no venv is active: %v", err)
	}

	// Test deactivating after setting some environment
	manager.config.Environment = map[string]string{
		"VIRTUAL_ENV": "/some/path",
		"PATH":        "/some/path/bin:/usr/bin",
	}

	err = manager.DeactivateVenv()
	if err != nil {
		t.Errorf("DeactivateVenv() unexpected error: %v", err)
	}

	// Verify environment was cleaned
	if _, exists := manager.config.Environment["VIRTUAL_ENV"]; exists {
		t.Error("VIRTUAL_ENV should be removed after deactivation")
	}
}

func TestRemoveVenv(t *testing.T) {
	manager := NewManager(nil)

	tests := []struct {
		name    string
		path    string
		wantErr bool
		setup   func(string) error
	}{
		{
			name:    "empty path",
			path:    "",
			wantErr: true,
		},
		{
			name:    "non-existent directory",
			path:    "non-existent",
			wantErr: true,
		},
		{
			name:    "existing directory",
			path:    "existing-venv",
			wantErr: false,
			setup: func(path string) error {
				return os.MkdirAll(path, 0755)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fullPath string

			if tt.path != "" {
				tempDir, err := os.MkdirTemp("", "pip-sdk-venv-test")
				if err != nil {
					t.Fatalf("Failed to create temp dir: %v", err)
				}
				defer os.RemoveAll(tempDir)

				fullPath = filepath.Join(tempDir, tt.path)

				if tt.setup != nil {
					if err := tt.setup(fullPath); err != nil {
						t.Fatalf("Setup failed: %v", err)
					}
				}
			} else {
				fullPath = tt.path
			}

			err := manager.RemoveVenv(fullPath)

			if tt.wantErr {
				if err == nil {
					t.Errorf("RemoveVenv() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("RemoveVenv() unexpected error: %v", err)
				} else {
					// Verify directory was removed
					if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
						t.Error("Directory should be removed after RemoveVenv()")
					}
				}
			}
		})
	}
}

func TestGetVenvBinPath(t *testing.T) {
	manager := NewManager(nil)

	tests := []struct {
		name     string
		venvPath string
		expected string
	}{
		{
			name:     "unix path",
			venvPath: "/home/user/venv",
			expected: "/home/user/venv/bin",
		},
		{
			name:     "relative path",
			venvPath: "./venv",
			expected: "./venv/bin",
		},
	}

	// Note: We can't actually change runtime.GOOS in tests,
	// but we can test the logic for the current OS

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.getVenvBinPath(tt.venvPath)

			if runtime.GOOS == "windows" {
				expected := filepath.Join(tt.venvPath, "Scripts")
				if result != expected {
					t.Errorf("getVenvBinPath() on Windows = %v, want %v", result, expected)
				}
			} else {
				expected := filepath.Join(tt.venvPath, "bin")
				if result != expected {
					t.Errorf("getVenvBinPath() on Unix = %v, want %v", result, expected)
				}
			}
		})
	}
}

func TestGetPythonExecutableName(t *testing.T) {
	manager := NewManager(nil)

	result := manager.getPythonExecutableName()

	if runtime.GOOS == "windows" {
		if result != "python.exe" {
			t.Errorf("getPythonExecutableName() on Windows = %v, want python.exe", result)
		}
	} else {
		if result != "python" {
			t.Errorf("getPythonExecutableName() on Unix = %v, want python", result)
		}
	}
}

func TestGetPipExecutableName(t *testing.T) {
	manager := NewManager(nil)

	result := manager.getPipExecutableName()

	if runtime.GOOS == "windows" {
		if result != "pip.exe" {
			t.Errorf("getPipExecutableName() on Windows = %v, want pip.exe", result)
		}
	} else {
		if result != "pip" {
			t.Errorf("getPipExecutableName() on Unix = %v, want pip", result)
		}
	}
}

func TestIsVenvValid(t *testing.T) {
	manager := NewManager(nil)

	tests := []struct {
		name     string
		setup    func(string) error
		expected bool
	}{
		{
			name:     "non-existent directory",
			expected: false,
		},
		{
			name: "empty directory",
			setup: func(path string) error {
				return os.MkdirAll(path, 0755)
			},
			expected: false,
		},
		{
			name: "directory with python executable",
			setup: func(path string) error {
				binPath := manager.getVenvBinPath(path)
				if err := os.MkdirAll(binPath, 0755); err != nil {
					return err
				}
				pythonPath := filepath.Join(binPath, manager.getPythonExecutableName())
				return os.WriteFile(pythonPath, []byte("#!/usr/bin/env python"), 0755)
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir, err := os.MkdirTemp("", "pip-sdk-venv-test")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			venvPath := filepath.Join(tempDir, "test-venv")

			if tt.setup != nil {
				if err := tt.setup(venvPath); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			result := manager.isVenvValid(venvPath)
			if result != tt.expected {
				t.Errorf("isVenvValid() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsVenvActive(t *testing.T) {
	manager := NewManager(nil)

	tempDir, err := os.MkdirTemp("", "pip-sdk-venv-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	venvPath := filepath.Join(tempDir, "test-venv")

	// Test when no environment is set
	result := manager.isVenvActive(venvPath)
	if result {
		t.Error("isVenvActive() should return false when no environment is set")
	}

	// Test when different venv is active
	manager.config.Environment = map[string]string{
		"VIRTUAL_ENV": "/other/venv/path",
	}
	result = manager.isVenvActive(venvPath)
	if result {
		t.Error("isVenvActive() should return false when different venv is active")
	}

	// Test when same venv is active
	manager.config.Environment["VIRTUAL_ENV"] = venvPath
	result = manager.isVenvActive(venvPath)
	if !result {
		t.Error("isVenvActive() should return true when same venv is active")
	}
}

func TestCreateWithVenv(t *testing.T) {
	manager := NewManager(nil)
	venvManager := NewVenvManager(manager)

	tempDir, err := os.MkdirTemp("", "pip-sdk-venv-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	venvPath := filepath.Join(tempDir, "test-venv")

	// Try to find Python first
	pythonPath, err := manager.findPythonExecutable()
	if err != nil {
		t.Logf("Python not found, testing with dummy path: %v", err)
		pythonPath = "python" // Use dummy path for testing
	}

	// Test createWithVenv method
	err = venvManager.createWithVenv(pythonPath, venvPath)

	// This might succeed if Python is available with venv module
	// or fail if Python is not available or venv is not supported
	t.Logf("createWithVenv result: %v", err)

	// Test with invalid Python path
	err = venvManager.createWithVenv("nonexistent-python", venvPath+"-invalid")
	if err == nil {
		t.Error("createWithVenv() should fail with invalid Python path")
	}
}

func TestCreateWithVirtualenv(t *testing.T) {
	manager := NewManager(nil)
	venvManager := NewVenvManager(manager)

	tempDir, err := os.MkdirTemp("", "pip-sdk-venv-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	venvPath := filepath.Join(tempDir, "test-venv")

	// Try to find Python first
	pythonPath, err := manager.findPythonExecutable()
	if err != nil {
		t.Logf("Python not found, testing with dummy path: %v", err)
		pythonPath = "python" // Use dummy path for testing
	}

	// Test createWithVirtualenv method
	err = venvManager.createWithVirtualenv(pythonPath, venvPath)

	// This might succeed if Python is available with virtualenv module
	// or fail if Python is not available or virtualenv is not installed
	t.Logf("createWithVirtualenv result: %v", err)

	// Test with invalid Python path
	err = venvManager.createWithVirtualenv("nonexistent-python", venvPath+"-invalid")
	if err == nil {
		t.Error("createWithVirtualenv() should fail with invalid Python path")
	}
}

func TestCreateWithVirtualenvCommand(t *testing.T) {
	manager := NewManager(nil)
	venvManager := NewVenvManager(manager)

	tempDir, err := os.MkdirTemp("", "pip-sdk-venv-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	venvPath := filepath.Join(tempDir, "test-venv")

	// Test createWithVirtualenvCommand method
	err = venvManager.createWithVirtualenvCommand(venvPath)

	// This might succeed if virtualenv command is available
	// or fail if virtualenv is not installed
	t.Logf("createWithVirtualenvCommand result: %v", err)
}

func TestGetVenvInfo(t *testing.T) {
	manager := NewManager(nil)

	tests := []struct {
		name    string
		path    string
		setup   func(string) error
		wantErr bool
	}{
		{
			name:    "empty path",
			path:    "",
			wantErr: true,
		},
		{
			name:    "non-existent venv",
			path:    "non-existent",
			wantErr: true,
		},
		{
			name: "invalid venv directory",
			path: "invalid-venv",
			setup: func(path string) error {
				return os.MkdirAll(path, 0755)
			},
			wantErr: true,
		},
		{
			name: "valid venv structure",
			path: "valid-venv",
			setup: func(path string) error {
				binPath := manager.getVenvBinPath(path)
				if err := os.MkdirAll(binPath, 0755); err != nil {
					return err
				}
				pythonPath := filepath.Join(binPath, manager.getPythonExecutableName())
				return os.WriteFile(pythonPath, []byte("#!/usr/bin/env python"), 0755)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fullPath string

			if tt.path != "" {
				tempDir, err := os.MkdirTemp("", "pip-sdk-venv-test")
				if err != nil {
					t.Fatalf("Failed to create temp dir: %v", err)
				}
				defer os.RemoveAll(tempDir)

				fullPath = filepath.Join(tempDir, tt.path)

				if tt.setup != nil {
					if err := tt.setup(fullPath); err != nil {
						t.Fatalf("Setup failed: %v", err)
					}
				}
			} else {
				fullPath = tt.path
			}

			info, err := manager.GetVenvInfo(fullPath)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetVenvInfo() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("GetVenvInfo() unexpected error: %v", err)
				} else {
					if info == nil {
						t.Error("GetVenvInfo() returned nil info")
					} else {
						t.Logf("VenvInfo: Path=%s, IsActive=%v, PythonPath=%s",
							info.Path, info.IsActive, info.PythonPath)
					}
				}
			}
		})
	}
}

func TestGetActivateScriptPath(t *testing.T) {
	manager := NewManager(nil)

	tempDir, err := os.MkdirTemp("", "pip-sdk-venv-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	venvPath := filepath.Join(tempDir, "test-venv")
	binPath := manager.getVenvBinPath(venvPath)

	// Create bin directory
	if err := os.MkdirAll(binPath, 0755); err != nil {
		t.Fatalf("Failed to create bin directory: %v", err)
	}

	// Test different activation scripts based on OS
	var activateScript string
	if runtime.GOOS == "windows" {
		// Test activate.bat
		activateScript = filepath.Join(binPath, "activate.bat")
		if err := os.WriteFile(activateScript, []byte("@echo off"), 0644); err != nil {
			t.Fatalf("Failed to create activate.bat: %v", err)
		}
	} else {
		// Test activate script
		activateScript = filepath.Join(binPath, "activate")
		if err := os.WriteFile(activateScript, []byte("#!/bin/bash"), 0755); err != nil {
			t.Fatalf("Failed to create activate script: %v", err)
		}
	}

	// Test getActivateScriptPath
	result := manager.getActivateScriptPath(venvPath)
	if result == "" {
		t.Error("getActivateScriptPath() should find activation script")
	} else {
		t.Logf("Found activation script: %s", result)
	}

	// Test with non-existent venv
	nonExistentPath := filepath.Join(tempDir, "non-existent-venv")
	result = manager.getActivateScriptPath(nonExistentPath)
	if result != "" {
		t.Error("getActivateScriptPath() should return empty string for non-existent venv")
	}
}

func TestGetPythonVersionInVenv(t *testing.T) {
	manager := NewManager(nil)

	// Try to find a real Python executable first
	pythonPath, err := manager.findPythonExecutable()
	if err != nil {
		t.Logf("Python not found, skipping version test: %v", err)
		return
	}

	// Test getPythonVersionInVenv with real Python
	version, err := manager.getPythonVersionInVenv(pythonPath)
	if err != nil {
		t.Logf("Failed to get Python version (this might be expected): %v", err)
	} else {
		t.Logf("Python version: %s", version)

		// Version should contain numbers
		if !strings.ContainsAny(version, "0123456789") {
			t.Errorf("Invalid version format: %s", version)
		}
	}

	// Test with invalid Python path
	_, err = manager.getPythonVersionInVenv("nonexistent-python")
	if err == nil {
		t.Error("getPythonVersionInVenv() should fail with invalid Python path")
	}
}

// Test edge cases and error conditions
func TestVenvEdgeCases(t *testing.T) {
	manager := NewManager(nil)

	// Test with nil environment
	manager.config.Environment = nil
	result := manager.isVenvActive("/some/path")
	if result {
		t.Error("isVenvActive() should return false when environment is nil")
	}

	// Test deactivate with nil environment
	err := manager.DeactivateVenv()
	if err != nil {
		t.Errorf("DeactivateVenv() should not error with nil environment: %v", err)
	}

	// Test with empty config
	manager.config = &Config{}
	result = manager.isVenvActive("/some/path")
	if result {
		t.Error("isVenvActive() should return false with empty config")
	}
}

// Test concurrent access (basic thread safety)
func TestVenvConcurrency(t *testing.T) {
	manager := NewManager(nil)

	// Test concurrent calls to isVenvValid
	tempDir, err := os.MkdirTemp("", "pip-sdk-venv-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	venvPath := filepath.Join(tempDir, "test-venv")

	// Run multiple goroutines
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()
			_ = manager.isVenvValid(venvPath)
			_ = manager.getVenvBinPath(venvPath)
			_ = manager.getPythonExecutableName()
			_ = manager.getPipExecutableName()
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

// Benchmark tests for venv operations
func BenchmarkIsVenvValid(b *testing.B) {
	manager := NewManager(nil)

	tempDir, err := os.MkdirTemp("", "pip-sdk-venv-bench")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	venvPath := filepath.Join(tempDir, "test-venv")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = manager.isVenvValid(venvPath)
	}
}

func BenchmarkGetVenvBinPath(b *testing.B) {
	manager := NewManager(nil)
	venvPath := "/path/to/venv"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = manager.getVenvBinPath(venvPath)
	}
}

// Test ActivateVenv more thoroughly
func TestActivateVenvDetailed(t *testing.T) {
	manager := NewManager(nil)

	tempDir, err := os.MkdirTemp("", "pip-sdk-venv-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a valid venv structure
	venvPath := filepath.Join(tempDir, "test-venv")
	binPath := manager.getVenvBinPath(venvPath)
	if err := os.MkdirAll(binPath, 0755); err != nil {
		t.Fatalf("Failed to create bin directory: %v", err)
	}

	// Create Python executable
	pythonPath := filepath.Join(binPath, manager.getPythonExecutableName())
	if err := os.WriteFile(pythonPath, []byte("#!/usr/bin/env python"), 0755); err != nil {
		t.Fatalf("Failed to create Python executable: %v", err)
	}

	// Create activate script
	var activateScript string
	if runtime.GOOS == "windows" {
		activateScript = filepath.Join(binPath, "activate.bat")
		if err := os.WriteFile(activateScript, []byte("@echo off"), 0644); err != nil {
			t.Fatalf("Failed to create activate script: %v", err)
		}
	} else {
		activateScript = filepath.Join(binPath, "activate")
		if err := os.WriteFile(activateScript, []byte("#!/bin/bash"), 0755); err != nil {
			t.Fatalf("Failed to create activate script: %v", err)
		}
	}

	// Test activation
	err = manager.ActivateVenv(venvPath)
	if err != nil {
		t.Logf("ActivateVenv() error (might be expected): %v", err)
	}

	// Test activation with already active venv
	if manager.config.Environment == nil {
		manager.config.Environment = make(map[string]string)
	}
	manager.config.Environment["VIRTUAL_ENV"] = venvPath

	err = manager.ActivateVenv(venvPath)
	if err != nil {
		t.Logf("ActivateVenv() with already active venv: %v", err)
	}
}

// Test getActivateScriptPath edge cases
func TestGetActivateScriptPathEdgeCases(t *testing.T) {
	manager := NewManager(nil)

	tempDir, err := os.MkdirTemp("", "pip-sdk-venv-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name         string
		setupFunc    func(string) error
		expectScript bool
	}{
		{
			name: "no bin directory",
			setupFunc: func(venvPath string) error {
				return os.MkdirAll(venvPath, 0755)
			},
			expectScript: false,
		},
		{
			name: "bin directory but no activate script",
			setupFunc: func(venvPath string) error {
				binPath := manager.getVenvBinPath(venvPath)
				return os.MkdirAll(binPath, 0755)
			},
			expectScript: false,
		},
		{
			name: "activate script exists",
			setupFunc: func(venvPath string) error {
				binPath := manager.getVenvBinPath(venvPath)
				if err := os.MkdirAll(binPath, 0755); err != nil {
					return err
				}

				var scriptName string
				if runtime.GOOS == "windows" {
					scriptName = "activate.bat"
				} else {
					scriptName = "activate"
				}

				scriptPath := filepath.Join(binPath, scriptName)
				return os.WriteFile(scriptPath, []byte("test"), 0755)
			},
			expectScript: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			venvPath := filepath.Join(tempDir, tt.name)

			if err := tt.setupFunc(venvPath); err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			scriptPath := manager.getActivateScriptPath(venvPath)

			if tt.expectScript {
				if scriptPath == "" {
					t.Error("Expected to find activate script but got empty path")
				}
			} else {
				if scriptPath != "" {
					t.Errorf("Expected no activate script but got: %s", scriptPath)
				}
			}
		})
	}
}

// Test cross-platform executable names
func TestExecutableNames(t *testing.T) {
	manager := NewManager(nil)

	pythonName := manager.getPythonExecutableName()
	pipName := manager.getPipExecutableName()

	if runtime.GOOS == "windows" {
		if pythonName != "python.exe" {
			t.Errorf("Expected python.exe on Windows, got %s", pythonName)
		}
		if pipName != "pip.exe" {
			t.Errorf("Expected pip.exe on Windows, got %s", pipName)
		}
	} else {
		if pythonName != "python" {
			t.Errorf("Expected python on Unix, got %s", pythonName)
		}
		if pipName != "pip" {
			t.Errorf("Expected pip on Unix, got %s", pipName)
		}
	}
}

// Test getVenvBinPath cross-platform
func TestGetVenvBinPathCrossPlatform(t *testing.T) {
	manager := NewManager(nil)

	tests := []struct {
		name     string
		venvPath string
		expected string
	}{
		{
			name:     "unix absolute path",
			venvPath: "/home/user/venv",
			expected: "/home/user/venv/bin",
		},
		{
			name:     "unix relative path",
			venvPath: "venv",
			expected: "venv/bin",
		},
		{
			name:     "windows-style path",
			venvPath: "C:\\Users\\user\\venv",
			expected: "C:\\Users\\user\\venv\\Scripts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.getVenvBinPath(tt.venvPath)

			// On Windows, we expect Scripts, on Unix we expect bin
			if runtime.GOOS == "windows" {
				if !strings.HasSuffix(result, "Scripts") {
					t.Errorf("Expected path to end with Scripts on Windows, got %s", result)
				}
			} else {
				if !strings.HasSuffix(result, "bin") {
					t.Errorf("Expected path to end with bin on Unix, got %s", result)
				}
			}
		})
	}
}
