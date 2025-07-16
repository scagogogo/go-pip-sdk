package pip

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// Installer handles pip installation across different operating systems
type Installer struct {
	manager *Manager
}

// NewInstaller creates a new pip installer
func NewInstaller(manager *Manager) *Installer {
	return &Installer{
		manager: manager,
	}
}

// IsInstalled checks if pip is installed and accessible
func (m *Manager) IsInstalled() (bool, error) {
	m.logDebug("Checking if pip is installed")

	// Try to find pip executable
	pipPath, err := m.findPipExecutable()
	if err != nil {
		return false, nil
	}

	// Try to run pip --version
	cmd := exec.CommandContext(m.ctx, pipPath, "--version")
	output, err := cmd.Output()
	if err != nil {
		m.logDebug("pip --version failed: %v", err)
		return false, nil
	}

	m.logDebug("pip version output: %s", string(output))
	return true, nil
}

// findPipExecutable finds the pip executable path
func (m *Manager) findPipExecutable() (string, error) {
	// If pip path is configured, use it
	if m.config.PipPath != "" {
		if _, err := os.Stat(m.config.PipPath); err == nil {
			return m.config.PipPath, nil
		}
	}

	// Try common pip command names
	pipCommands := []string{"pip", "pip3", "python -m pip", "python3 -m pip"}

	for _, cmd := range pipCommands {
		parts := strings.Fields(cmd)
		if len(parts) == 0 {
			continue
		}

		// Check if command exists
		if len(parts) == 1 {
			if path, err := exec.LookPath(parts[0]); err == nil {
				return path, nil
			}
		} else {
			// For commands like "python -m pip"
			if _, err := exec.LookPath(parts[0]); err == nil {
				return cmd, nil
			}
		}
	}

	return "", fmt.Errorf("pip executable not found")
}

// findPythonExecutable finds the Python executable path
func (m *Manager) findPythonExecutable() (string, error) {
	// If python path is configured, use it
	if m.config.PythonPath != "" {
		if _, err := os.Stat(m.config.PythonPath); err == nil {
			return m.config.PythonPath, nil
		}
	}

	// Try common python command names
	pythonCommands := []string{"python", "python3", "py"}

	for _, cmd := range pythonCommands {
		if path, err := exec.LookPath(cmd); err == nil {
			return path, nil
		}
	}

	return "", ErrPythonNotFound
}

// Install installs pip on the system
func (m *Manager) Install() error {
	m.logInfo("Installing pip for %s", runtime.GOOS)

	installer := NewInstaller(m)

	switch GetOSType() {
	case OSWindows:
		return installer.installOnWindows()
	case OSMacOS:
		return installer.installOnMacOS()
	case OSLinux:
		return installer.installOnLinux()
	default:
		return ErrUnsupportedOS
	}
}

// installOnWindows installs pip on Windows
func (i *Installer) installOnWindows() error {
	i.manager.logInfo("Installing pip on Windows")

	// First, try to find Python
	pythonPath, err := i.manager.findPythonExecutable()
	if err != nil {
		return fmt.Errorf("Python not found. Please install Python first: %w", err)
	}

	// Check if pip is already available through python -m pip
	cmd := exec.CommandContext(i.manager.ctx, pythonPath, "-m", "pip", "--version")
	if err := cmd.Run(); err == nil {
		i.manager.logInfo("pip is already available through python -m pip")
		return nil
	}

	// Download and install get-pip.py
	return i.installUsingGetPip(pythonPath)
}

// installOnMacOS installs pip on macOS
func (i *Installer) installOnMacOS() error {
	i.manager.logInfo("Installing pip on macOS")

	// First, try to find Python
	pythonPath, err := i.manager.findPythonExecutable()
	if err != nil {
		// Try to install Python using Homebrew if available
		if err := i.installPythonWithHomebrew(); err != nil {
			return fmt.Errorf("Python not found and Homebrew installation failed. Please install Python first: %w", err)
		}

		// Try to find Python again
		pythonPath, err = i.manager.findPythonExecutable()
		if err != nil {
			return fmt.Errorf("Python still not found after Homebrew installation: %w", err)
		}
	}

	// Check if pip is already available
	cmd := exec.CommandContext(i.manager.ctx, pythonPath, "-m", "pip", "--version")
	if err := cmd.Run(); err == nil {
		i.manager.logInfo("pip is already available through python -m pip")
		return nil
	}

	// Try to install pip using ensurepip
	if err := i.installUsingEnsurepip(pythonPath); err == nil {
		return nil
	}

	// Fallback to get-pip.py
	return i.installUsingGetPip(pythonPath)
}

// installOnLinux installs pip on Linux
func (i *Installer) installOnLinux() error {
	i.manager.logInfo("Installing pip on Linux")

	// Try package manager first
	if err := i.installUsingPackageManager(); err == nil {
		return nil
	}

	// Try to find Python
	pythonPath, err := i.manager.findPythonExecutable()
	if err != nil {
		return fmt.Errorf("Python not found. Please install Python first: %w", err)
	}

	// Check if pip is already available
	cmd := exec.CommandContext(i.manager.ctx, pythonPath, "-m", "pip", "--version")
	if err := cmd.Run(); err == nil {
		i.manager.logInfo("pip is already available through python -m pip")
		return nil
	}

	// Try to install pip using ensurepip
	if err := i.installUsingEnsurepip(pythonPath); err == nil {
		return nil
	}

	// Fallback to get-pip.py
	return i.installUsingGetPip(pythonPath)
}

// installUsingPackageManager installs pip using system package manager on Linux
func (i *Installer) installUsingPackageManager() error {
	i.manager.logInfo("Trying to install pip using package manager")

	// Try different package managers
	packageManagers := []struct {
		cmd  string
		args []string
	}{
		{"apt-get", []string{"update", "&&", "apt-get", "install", "-y", "python3-pip"}},
		{"yum", []string{"install", "-y", "python3-pip"}},
		{"dnf", []string{"install", "-y", "python3-pip"}},
		{"pacman", []string{"-S", "--noconfirm", "python-pip"}},
		{"zypper", []string{"install", "-y", "python3-pip"}},
	}

	for _, pm := range packageManagers {
		if _, err := exec.LookPath(pm.cmd); err == nil {
			i.manager.logInfo("Found package manager: %s", pm.cmd)

			cmd := exec.CommandContext(i.manager.ctx, "sudo", append([]string{pm.cmd}, pm.args...)...)
			if err := cmd.Run(); err == nil {
				i.manager.logInfo("Successfully installed pip using %s", pm.cmd)
				return nil
			}
		}
	}

	return fmt.Errorf("no suitable package manager found")
}

// installPythonWithHomebrew installs Python using Homebrew on macOS
func (i *Installer) installPythonWithHomebrew() error {
	i.manager.logInfo("Trying to install Python using Homebrew")

	// Check if Homebrew is available
	if _, err := exec.LookPath("brew"); err != nil {
		return fmt.Errorf("Homebrew not found")
	}

	// Install Python
	cmd := exec.CommandContext(i.manager.ctx, "brew", "install", "python")
	return cmd.Run()
}

// installUsingEnsurepip installs pip using Python's ensurepip module
func (i *Installer) installUsingEnsurepip(pythonPath string) error {
	i.manager.logInfo("Trying to install pip using ensurepip")

	cmd := exec.CommandContext(i.manager.ctx, pythonPath, "-m", "ensurepip", "--upgrade")
	output, err := cmd.CombinedOutput()

	if err != nil {
		i.manager.logError("ensurepip failed: %v, output: %s", err, string(output))
		return err
	}

	i.manager.logInfo("Successfully installed pip using ensurepip")
	return nil
}

// installUsingGetPip downloads and runs get-pip.py
func (i *Installer) installUsingGetPip(pythonPath string) error {
	i.manager.logInfo("Installing pip using get-pip.py")

	// Download get-pip.py
	getPipPath, err := i.downloadGetPip()
	if err != nil {
		return fmt.Errorf("failed to download get-pip.py: %w", err)
	}
	defer os.Remove(getPipPath)

	// Run get-pip.py
	cmd := exec.CommandContext(i.manager.ctx, pythonPath, getPipPath)
	output, err := cmd.CombinedOutput()

	if err != nil {
		i.manager.logError("get-pip.py failed: %v, output: %s", err, string(output))
		return fmt.Errorf("get-pip.py execution failed: %w", err)
	}

	i.manager.logInfo("Successfully installed pip using get-pip.py")
	return nil
}

// downloadGetPip downloads get-pip.py script
func (i *Installer) downloadGetPip() (string, error) {
	i.manager.logInfo("Downloading get-pip.py")

	// Create temporary file
	tmpFile, err := os.CreateTemp("", "get-pip-*.py")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	// Download get-pip.py
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get("https://bootstrap.pypa.io/get-pip.py")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download get-pip.py: HTTP %d", resp.StatusCode)
	}

	// Copy content to temporary file
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

// GetVersion returns the pip version
func (m *Manager) GetVersion() (string, error) {
	m.logDebug("Getting pip version")

	pipPath, err := m.findPipExecutable()
	if err != nil {
		return "", ErrPipNotInstalled
	}

	var cmd *exec.Cmd
	if strings.Contains(pipPath, " ") {
		// Handle commands like "python -m pip"
		parts := strings.Fields(pipPath)
		args := append(parts[1:], "--version")
		cmd = exec.CommandContext(m.ctx, parts[0], args...)
	} else {
		cmd = exec.CommandContext(m.ctx, pipPath, "--version")
	}

	output, err := cmd.Output()
	if err != nil {
		return "", m.createPipError(cmd.String(), string(output), -1, err)
	}

	// Parse version from output (e.g., "pip 21.3.1 from ...")
	versionStr := string(output)
	parts := strings.Fields(versionStr)
	if len(parts) >= 2 {
		return parts[1], nil
	}

	return strings.TrimSpace(versionStr), nil
}
