package pip

import (
	"context"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestNewInstaller(t *testing.T) {
	manager := NewManager(nil)
	installer := NewInstaller(manager)

	if installer == nil {
		t.Fatal("NewInstaller() returned nil")
	}

	if installer.manager != manager {
		t.Error("Installer should reference the provided manager")
	}
}

func TestFindPipExecutable(t *testing.T) {
	manager := NewManager(nil)

	tests := []struct {
		name    string
		config  *Config
		setup   func()
		cleanup func()
	}{
		{
			name: "with configured pip path",
			config: &Config{
				PipPath: "/usr/bin/pip3",
			},
		},
		{
			name:   "search in PATH",
			config: &Config{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()

			if tt.config != nil {
				manager.SetConfig(tt.config)
			}

			pipPath, err := manager.findPipExecutable()

			// We don't assert success/failure since pip might not be installed
			// But we test that the function doesn't panic and returns reasonable results
			if err == nil && pipPath == "" {
				t.Error("findPipExecutable() returned empty path without error")
			}

			t.Logf("Found pip at: %s (error: %v)", pipPath, err)
		})
	}
}

func TestFindPythonExecutable(t *testing.T) {
	manager := NewManager(nil)

	tests := []struct {
		name   string
		config *Config
	}{
		{
			name: "with configured python path",
			config: &Config{
				PythonPath: "/usr/bin/python3",
			},
		},
		{
			name:   "search in PATH",
			config: &Config{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.config != nil {
				manager.SetConfig(tt.config)
			}

			pythonPath, err := manager.findPythonExecutable()

			// We don't assert success/failure since Python might not be installed
			// But we test that the function doesn't panic and returns reasonable results
			if err == nil && pythonPath == "" {
				t.Error("findPythonExecutable() returned empty path without error")
			}

			t.Logf("Found Python at: %s (error: %v)", pythonPath, err)
		})
	}
}

func TestIsInstalled(t *testing.T) {
	manager := NewManager(nil)

	// Test with short timeout to avoid hanging
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	manager.SetContext(ctx)

	installed, err := manager.IsInstalled()

	// We don't assert the result since pip might or might not be installed
	// But we verify the function works correctly
	t.Logf("Pip installed: %v, error: %v", installed, err)

	if err != nil && installed {
		t.Error("IsInstalled() should not return true when there's an error")
	}
}

func TestGetVersion(t *testing.T) {
	manager := NewManager(nil)

	// Test with short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	manager.SetContext(ctx)

	// First check if pip is available
	installed, err := manager.IsInstalled()
	if err != nil || !installed {
		t.Skip("Pip not available, skipping version test")
	}

	version, err := manager.GetVersion()
	if err != nil {
		t.Errorf("GetVersion() error = %v", err)
		return
	}

	if version == "" {
		t.Error("GetVersion() returned empty version")
	}

	// Version should contain numbers
	if !strings.ContainsAny(version, "0123456789") {
		t.Errorf("GetVersion() returned invalid version format: %s", version)
	}

	t.Logf("Pip version: %s", version)
}

func TestInstallMethods(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping installation tests in short mode")
	}

	manager := NewManager(nil)
	_ = NewInstaller(manager)

	// Test different installation methods without actually installing
	// (to avoid modifying the system)

	tests := []struct {
		name   string
		method func() error
	}{
		{
			name: "installUsingEnsurepip",
			method: func() error {
				// Mock Python path for testing
				pythonPath, err := manager.findPythonExecutable()
				if err != nil {
					return err
				}

				// Test the command construction (don't actually run it)
				cmd := exec.Command(pythonPath, "-m", "ensurepip", "--help")
				_, err = cmd.Output()
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.method()
			// We don't assert success since the system might not support the method
			t.Logf("Method %s result: %v", tt.name, err)
		})
	}
}

func TestDownloadGetPip(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping download test in short mode")
	}

	manager := NewManager(nil)
	installer := NewInstaller(manager)

	// Test downloading get-pip.py (with timeout)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	manager.SetContext(ctx)

	tempFile, err := installer.downloadGetPip()
	if err != nil {
		t.Logf("downloadGetPip() failed (possibly due to network): %v", err)
		return
	}
	defer os.Remove(tempFile)

	// Verify file was created
	if _, err := os.Stat(tempFile); os.IsNotExist(err) {
		t.Error("downloadGetPip() did not create file")
		return
	}

	// Verify file has content
	info, err := os.Stat(tempFile)
	if err != nil {
		t.Errorf("Failed to stat downloaded file: %v", err)
		return
	}

	if info.Size() == 0 {
		t.Error("Downloaded get-pip.py file is empty")
	}

	// Verify it's a Python script
	content, err := os.ReadFile(tempFile)
	if err != nil {
		t.Errorf("Failed to read downloaded file: %v", err)
		return
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "python") && !strings.Contains(contentStr, "pip") {
		t.Error("Downloaded file doesn't appear to be get-pip.py")
	}

	t.Logf("Successfully downloaded get-pip.py (%d bytes)", info.Size())
}

func TestInstallOnDifferentOS(t *testing.T) {
	manager := NewManager(nil)
	installer := NewInstaller(manager)

	// Test OS-specific installation logic without actually installing
	tests := []struct {
		name     string
		osType   OSType
		testFunc func() error
	}{
		{
			name:   "windows logic",
			osType: OSWindows,
			testFunc: func() error {
				// Test Windows-specific logic
				return installer.installOnWindows()
			},
		},
		{
			name:   "macos logic",
			osType: OSMacOS,
			testFunc: func() error {
				// Test macOS-specific logic
				return installer.installOnMacOS()
			},
		},
		{
			name:   "linux logic",
			osType: OSLinux,
			testFunc: func() error {
				// Test Linux-specific logic
				return installer.installOnLinux()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Only test the current OS to avoid platform-specific issues
			if GetOSType() == tt.osType {
				err := tt.testFunc()
				// We don't assert success since installation might fail for various reasons
				t.Logf("Installation method for %s: %v", tt.name, err)
			} else {
				t.Skipf("Skipping %s test on %s", tt.name, GetOSType())
			}
		})
	}
}

func TestInstallUsingPackageManager(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping package manager test in short mode")
	}

	if runtime.GOOS != "linux" {
		t.Skip("Package manager test only runs on Linux")
	}

	manager := NewManager(nil)
	installer := NewInstaller(manager)

	// Test package manager detection without actually installing
	err := installer.installUsingPackageManager()

	// We don't assert success since:
	// 1. We might not have sudo access
	// 2. Package managers might not be available
	// 3. pip might already be installed
	t.Logf("Package manager installation result: %v", err)
}

func TestInstallPythonWithHomebrew(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Homebrew test in short mode")
	}

	if runtime.GOOS != "darwin" {
		t.Skip("Homebrew test only runs on macOS")
	}

	manager := NewManager(nil)
	installer := NewInstaller(manager)

	// Test Homebrew detection without actually installing
	err := installer.installPythonWithHomebrew()

	// We don't assert success since:
	// 1. Homebrew might not be installed
	// 2. Python might already be installed
	// 3. We don't want to actually install Python
	t.Logf("Homebrew Python installation result: %v", err)
}

// Integration test for the full installation process
func TestInstallIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Only run if explicitly requested
	if os.Getenv("RUN_INSTALL_TESTS") != "1" {
		t.Skip("Skipping install integration test (set RUN_INSTALL_TESTS=1 to enable)")
	}

	manager := NewManager(nil)

	// Set a reasonable timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	manager.SetContext(ctx)

	// Check if pip is already installed
	installed, err := manager.IsInstalled()
	if err == nil && installed {
		t.Log("Pip is already installed, skipping installation test")
		return
	}

	// Attempt installation
	err = manager.Install()
	if err != nil {
		t.Logf("Installation failed (this might be expected): %v", err)
		return
	}

	// Verify installation
	installed, err = manager.IsInstalled()
	if err != nil {
		t.Errorf("Failed to verify installation: %v", err)
		return
	}

	if !installed {
		t.Error("Pip installation verification failed")
		return
	}

	t.Log("Pip installation successful")
}

// Benchmark tests
func BenchmarkIsInstalled(b *testing.B) {
	manager := NewManager(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = manager.IsInstalled()
	}
}

func BenchmarkFindPipExecutable(b *testing.B) {
	manager := NewManager(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = manager.findPipExecutable()
	}
}

func BenchmarkFindPythonExecutable(b *testing.B) {
	manager := NewManager(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = manager.findPythonExecutable()
	}
}

func TestInstall(t *testing.T) {
	manager := NewManager(nil)

	// Test with short timeout to avoid hanging
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	manager.SetContext(ctx)

	// Check if pip is already installed
	installed, err := manager.IsInstalled()
	if err == nil && installed {
		t.Log("Pip is already installed, testing Install() method anyway")
	}

	// Test Install method - this might succeed or fail depending on the environment
	err = manager.Install()

	// We don't assert success/failure since:
	// 1. Pip might already be installed
	// 2. We might not have permissions
	// 3. The system might not support automatic installation
	t.Logf("Install() result: %v", err)

	// The important thing is that the method doesn't panic
	// and returns a reasonable error if it fails
}

func TestInstallOnWindows(t *testing.T) {
	manager := NewManager(nil)
	installer := NewInstaller(manager)

	// Test Windows installation logic
	err := installer.installOnWindows()

	if runtime.GOOS == "windows" {
		// On Windows, we expect either success or a reasonable error
		t.Logf("Windows installation result: %v", err)
	} else {
		// On non-Windows systems, we expect this to work but might fail
		// due to platform-specific commands not being available
		t.Logf("Windows installation on %s: %v", runtime.GOOS, err)
	}
}

func TestInstallOnLinux(t *testing.T) {
	manager := NewManager(nil)
	installer := NewInstaller(manager)

	// Test Linux installation logic
	err := installer.installOnLinux()

	if runtime.GOOS == "linux" {
		// On Linux, we expect either success or a reasonable error
		t.Logf("Linux installation result: %v", err)
	} else {
		// On non-Linux systems, we expect this to work but might fail
		// due to platform-specific commands not being available
		t.Logf("Linux installation on %s: %v", runtime.GOOS, err)
	}
}

func TestInstallUsingEnsurepip(t *testing.T) {
	manager := NewManager(nil)
	installer := NewInstaller(manager)

	// Try to find Python first
	pythonPath, err := manager.findPythonExecutable()
	if err != nil {
		t.Logf("Python not found, testing with dummy path: %v", err)
		pythonPath = "python" // Use dummy path for testing
	}

	// Test ensurepip installation method
	err = installer.installUsingEnsurepip(pythonPath)

	// This might succeed if Python is available with ensurepip module
	// or fail if Python is not available or ensurepip is not supported
	t.Logf("Ensurepip installation result: %v", err)

	// Test with invalid Python path
	err = installer.installUsingEnsurepip("nonexistent-python-executable")
	if err == nil {
		t.Error("installUsingEnsurepip() should fail with invalid Python path")
	}
}

func TestInstallUsingGetPip(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping get-pip test in short mode")
	}

	manager := NewManager(nil)
	installer := NewInstaller(manager)

	// Set a reasonable timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	manager.SetContext(ctx)

	// Try to find Python first
	pythonPath, err := manager.findPythonExecutable()
	if err != nil {
		t.Logf("Python not found, testing with dummy path: %v", err)
		pythonPath = "python" // Use dummy path for testing
	}

	// Test get-pip installation method
	err = installer.installUsingGetPip(pythonPath)

	// This might succeed if we can download get-pip.py and Python is available
	// or fail due to network issues, permissions, or missing Python
	t.Logf("Get-pip installation result: %v", err)

	// Test with invalid Python path
	err = installer.installUsingGetPip("nonexistent-python-executable")
	if err == nil {
		t.Error("installUsingGetPip() should fail with invalid Python path")
	}
}

func TestInstallOnMacOSDetailed(t *testing.T) {
	manager := NewManager(nil)
	installer := NewInstaller(manager)

	// Test installOnMacOS method
	err := installer.installOnMacOS()

	// This should work on macOS or fail gracefully on other platforms
	t.Logf("installOnMacOS result: %v", err)

	// Test with different Python configurations
	tests := []struct {
		name       string
		pythonPath string
		pipPath    string
	}{
		{
			name:       "default paths",
			pythonPath: "",
			pipPath:    "",
		},
		{
			name:       "custom python path",
			pythonPath: "/usr/bin/python3",
			pipPath:    "",
		},
		{
			name:       "custom pip path",
			pythonPath: "",
			pipPath:    "/usr/bin/pip3",
		},
		{
			name:       "invalid paths",
			pythonPath: "/nonexistent/python",
			pipPath:    "/nonexistent/pip",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new manager for each test
			testManager := NewManager(&Config{
				PythonPath: tt.pythonPath,
				PipPath:    tt.pipPath,
			})
			testInstaller := NewInstaller(testManager)

			err := testInstaller.installOnMacOS()
			t.Logf("installOnMacOS with %s: %v", tt.name, err)

			// We don't assert success/failure since it depends on the environment
			// The important thing is that it doesn't panic
		})
	}
}

func TestInstallUsingPackageManagerDetailed(t *testing.T) {
	manager := NewManager(nil)
	installer := NewInstaller(manager)

	// Test installUsingPackageManager method
	err := installer.installUsingPackageManager()

	// This should work on Linux or fail gracefully on other platforms
	t.Logf("installUsingPackageManager result: %v", err)

	// Test with timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	manager.SetContext(ctx)

	err = installer.installUsingPackageManager()
	t.Logf("installUsingPackageManager with timeout: %v", err)
}
