package pip

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// VenvManager handles virtual environment operations
type VenvManager struct {
	manager *Manager
}

// NewVenvManager creates a new virtual environment manager
func NewVenvManager(manager *Manager) *VenvManager {
	return &VenvManager{
		manager: manager,
	}
}

// CreateVenv creates a virtual environment
func (m *Manager) CreateVenv(path string) error {
	if path == "" {
		return &PipError{
			Type:    "invalid_path",
			Message: "virtual environment path cannot be empty",
		}
	}

	m.logInfo("Creating virtual environment: %s", path)

	// Check if directory already exists
	if _, err := os.Stat(path); err == nil {
		return &PipError{
			Type:    "venv_already_exists",
			Message: fmt.Sprintf("directory already exists: %s", path),
		}
	}

	// Find Python executable
	pythonPath, err := m.findPythonExecutable()
	if err != nil {
		return err
	}

	// Try different methods to create virtual environment
	venvManager := NewVenvManager(m)

	// Method 1: python -m venv
	if err := venvManager.createWithVenv(pythonPath, path); err == nil {
		m.logInfo("Successfully created virtual environment using venv module")
		return nil
	}

	// Method 2: python -m virtualenv
	if err := venvManager.createWithVirtualenv(pythonPath, path); err == nil {
		m.logInfo("Successfully created virtual environment using virtualenv")
		return nil
	}

	// Method 3: virtualenv command
	if err := venvManager.createWithVirtualenvCommand(path); err == nil {
		m.logInfo("Successfully created virtual environment using virtualenv command")
		return nil
	}

	return &PipError{
		Type:    "venv_creation_failed",
		Message: "failed to create virtual environment using any available method",
	}
}

// createWithVenv creates virtual environment using python -m venv
func (vm *VenvManager) createWithVenv(pythonPath, path string) error {
	vm.manager.logDebug("Trying to create venv using python -m venv")

	cmd := exec.CommandContext(vm.manager.ctx, pythonPath, "-m", "venv", path)
	output, err := cmd.CombinedOutput()

	if err != nil {
		vm.manager.logDebug("venv creation failed: %v, output: %s", err, string(output))
		return err
	}

	return nil
}

// createWithVirtualenv creates virtual environment using python -m virtualenv
func (vm *VenvManager) createWithVirtualenv(pythonPath, path string) error {
	vm.manager.logDebug("Trying to create venv using python -m virtualenv")

	cmd := exec.CommandContext(vm.manager.ctx, pythonPath, "-m", "virtualenv", path)
	output, err := cmd.CombinedOutput()

	if err != nil {
		vm.manager.logDebug("virtualenv creation failed: %v, output: %s", err, string(output))
		return err
	}

	return nil
}

// createWithVirtualenvCommand creates virtual environment using virtualenv command
func (vm *VenvManager) createWithVirtualenvCommand(path string) error {
	vm.manager.logDebug("Trying to create venv using virtualenv command")

	// Check if virtualenv command exists
	if _, err := exec.LookPath("virtualenv"); err != nil {
		return err
	}

	cmd := exec.CommandContext(vm.manager.ctx, "virtualenv", path)
	output, err := cmd.CombinedOutput()

	if err != nil {
		vm.manager.logDebug("virtualenv command failed: %v, output: %s", err, string(output))
		return err
	}

	return nil
}

// ActivateVenv activates a virtual environment
func (m *Manager) ActivateVenv(path string) error {
	if path == "" {
		return &PipError{
			Type:    "invalid_path",
			Message: "virtual environment path cannot be empty",
		}
	}

	m.logInfo("Activating virtual environment: %s", path)

	// Check if virtual environment exists
	if !m.isVenvValid(path) {
		return &PipError{
			Type:    "venv_not_found",
			Message: fmt.Sprintf("virtual environment not found or invalid: %s", path),
		}
	}

	// Get activation script path
	activateScript := m.getActivateScriptPath(path)
	if activateScript == "" {
		return &PipError{
			Type:    "venv_activation_failed",
			Message: "activation script not found",
		}
	}

	// Update environment variables
	venvBinPath := m.getVenvBinPath(path)
	pythonPath := filepath.Join(venvBinPath, m.getPythonExecutableName())
	pipPath := filepath.Join(venvBinPath, m.getPipExecutableName())

	// Update manager configuration
	m.config.PythonPath = pythonPath
	m.config.PipPath = pipPath

	// Set environment variables
	if m.config.Environment == nil {
		m.config.Environment = make(map[string]string)
	}

	m.config.Environment["VIRTUAL_ENV"] = path
	m.config.Environment["PATH"] = venvBinPath + string(os.PathListSeparator) + os.Getenv("PATH")

	m.logInfo("Virtual environment activated: %s", path)
	return nil
}

// DeactivateVenv deactivates the current virtual environment
func (m *Manager) DeactivateVenv() error {
	m.logInfo("Deactivating virtual environment")

	// Reset paths
	m.config.PythonPath = ""
	m.config.PipPath = ""

	// Remove virtual environment variables
	if m.config.Environment != nil {
		delete(m.config.Environment, "VIRTUAL_ENV")
		delete(m.config.Environment, "PATH")
	}

	m.logInfo("Virtual environment deactivated")
	return nil
}

// RemoveVenv removes a virtual environment
func (m *Manager) RemoveVenv(path string) error {
	if path == "" {
		return &PipError{
			Type:    "invalid_path",
			Message: "virtual environment path cannot be empty",
		}
	}

	m.logInfo("Removing virtual environment: %s", path)

	// Check if directory exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &PipError{
			Type:    "venv_not_found",
			Message: fmt.Sprintf("virtual environment not found: %s", path),
		}
	}

	// Remove directory
	if err := os.RemoveAll(path); err != nil {
		return &PipError{
			Type:    "venv_removal_failed",
			Message: fmt.Sprintf("failed to remove virtual environment: %v", err),
		}
	}

	m.logInfo("Virtual environment removed: %s", path)
	return nil
}

// GetVenvInfo returns information about a virtual environment
func (m *Manager) GetVenvInfo(path string) (*VenvInfo, error) {
	if path == "" {
		return nil, &PipError{
			Type:    "invalid_path",
			Message: "virtual environment path cannot be empty",
		}
	}

	if !m.isVenvValid(path) {
		return nil, &PipError{
			Type:    "venv_not_found",
			Message: fmt.Sprintf("virtual environment not found or invalid: %s", path),
		}
	}

	info := &VenvInfo{
		Path:     path,
		IsActive: m.isVenvActive(path),
	}

	// Get Python path
	venvBinPath := m.getVenvBinPath(path)
	info.PythonPath = filepath.Join(venvBinPath, m.getPythonExecutableName())

	// Get creation time (if available)
	if stat, err := os.Stat(path); err == nil {
		info.CreatedAt = stat.ModTime()
	}

	// Get Python version
	if version, err := m.getPythonVersionInVenv(info.PythonPath); err == nil {
		info.PythonVersion = version
	}

	return info, nil
}

// isVenvValid checks if a virtual environment is valid
func (m *Manager) isVenvValid(path string) bool {
	// Check if directory exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	// Check for required files/directories
	venvBinPath := m.getVenvBinPath(path)
	pythonPath := filepath.Join(venvBinPath, m.getPythonExecutableName())

	if _, err := os.Stat(pythonPath); os.IsNotExist(err) {
		return false
	}

	return true
}

// isVenvActive checks if a virtual environment is currently active
func (m *Manager) isVenvActive(path string) bool {
	if m.config.Environment == nil {
		return false
	}

	virtualEnv, exists := m.config.Environment["VIRTUAL_ENV"]
	if !exists {
		return false
	}

	// Compare absolute paths
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	absVirtualEnv, err := filepath.Abs(virtualEnv)
	if err != nil {
		return false
	}

	return absPath == absVirtualEnv
}

// getVenvBinPath returns the bin/Scripts directory path for a virtual environment
func (m *Manager) getVenvBinPath(venvPath string) string {
	if runtime.GOOS == "windows" {
		return filepath.Join(venvPath, "Scripts")
	}
	return filepath.Join(venvPath, "bin")
}

// getActivateScriptPath returns the activation script path for a virtual environment
func (m *Manager) getActivateScriptPath(venvPath string) string {
	binPath := m.getVenvBinPath(venvPath)

	if runtime.GOOS == "windows" {
		activateScript := filepath.Join(binPath, "activate.bat")
		if _, err := os.Stat(activateScript); err == nil {
			return activateScript
		}

		activateScript = filepath.Join(binPath, "Activate.ps1")
		if _, err := os.Stat(activateScript); err == nil {
			return activateScript
		}
	} else {
		activateScript := filepath.Join(binPath, "activate")
		if _, err := os.Stat(activateScript); err == nil {
			return activateScript
		}
	}

	return ""
}

// getPythonExecutableName returns the Python executable name for the current OS
func (m *Manager) getPythonExecutableName() string {
	if runtime.GOOS == "windows" {
		return "python.exe"
	}
	return "python"
}

// getPipExecutableName returns the pip executable name for the current OS
func (m *Manager) getPipExecutableName() string {
	if runtime.GOOS == "windows" {
		return "pip.exe"
	}
	return "pip"
}

// getPythonVersionInVenv gets the Python version in a virtual environment
func (m *Manager) getPythonVersionInVenv(pythonPath string) (string, error) {
	cmd := exec.CommandContext(m.ctx, pythonPath, "--version")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Parse version from output (e.g., "Python 3.9.7")
	versionStr := strings.TrimSpace(string(output))
	parts := strings.Fields(versionStr)
	if len(parts) >= 2 {
		return parts[1], nil
	}

	return versionStr, nil
}
