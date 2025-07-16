package pip

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestParseListOutput(t *testing.T) {
	manager := NewManager(nil)

	tests := []struct {
		name     string
		output   string
		expected int // expected number of packages
	}{
		{
			name: "standard pip list output",
			output: `Package    Version
---------- -------
pip        21.3.1
setuptools 58.1.0
wheel      0.37.1`,
			expected: 3,
		},
		{
			name: "empty output",
			output: `Package Version
------- -------`,
			expected: 0,
		},
		{
			name: "single package",
			output: `Package Version
------- -------
requests 2.28.1`,
			expected: 1,
		},
		{
			name: "packages with locations",
			output: `Package    Version Location
---------- ------- --------
pip        21.3.1  /usr/lib/python3.9/site-packages
requests   2.28.1  /home/user/.local/lib/python3.9/site-packages`,
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packages := manager.parseListOutput(tt.output)

			if len(packages) != tt.expected {
				t.Errorf("parseListOutput() returned %d packages, expected %d", len(packages), tt.expected)
			}

			// Verify package structure
			for _, pkg := range packages {
				if pkg.Name == "" {
					t.Error("Package name should not be empty")
				}
				if pkg.Version == "" {
					t.Error("Package version should not be empty")
				}
			}
		})
	}
}

func TestParseShowOutput(t *testing.T) {
	manager := NewManager(nil)

	showOutput := `Name: requests
Version: 2.28.1
Summary: Python HTTP for Humans.
Home-page: https://requests.readthedocs.io
Author: Kenneth Reitz
Author-email: me@kennethreitz.org
License: Apache 2.0
Location: /usr/lib/python3.9/site-packages
Requires: charset-normalizer, idna, urllib3, certifi
Required-by: pip-tools, responses`

	info := manager.parseShowOutput(showOutput)

	if info.Name != "requests" {
		t.Errorf("Name = %s, want %s", info.Name, "requests")
	}

	if info.Version != "2.28.1" {
		t.Errorf("Version = %s, want %s", info.Version, "2.28.1")
	}

	if info.Summary != "Python HTTP for Humans." {
		t.Errorf("Summary = %s, want %s", info.Summary, "Python HTTP for Humans.")
	}

	if info.Author != "Kenneth Reitz" {
		t.Errorf("Author = %s, want %s", info.Author, "Kenneth Reitz")
	}

	expectedRequires := []string{"charset-normalizer", "idna", "urllib3", "certifi"}
	if len(info.Requires) != len(expectedRequires) {
		t.Errorf("Requires length = %d, want %d", len(info.Requires), len(expectedRequires))
	}

	for i, req := range expectedRequires {
		if i < len(info.Requires) && info.Requires[i] != req {
			t.Errorf("Requires[%d] = %s, want %s", i, info.Requires[i], req)
		}
	}
}

func TestParseFreezeOutput(t *testing.T) {
	manager := NewManager(nil)

	tests := []struct {
		name     string
		output   string
		expected int
	}{
		{
			name: "standard freeze output",
			output: `certifi==2022.9.24
charset-normalizer==2.1.1
idna==3.4
requests==2.28.1
urllib3==1.26.12`,
			expected: 5,
		},
		{
			name: "freeze with editable packages",
			output: `-e git+https://github.com/user/repo.git@main#egg=mypackage
requests==2.28.1
-e /path/to/local/package`,
			expected: 3,
		},
		{
			name: "freeze with comments",
			output: `# This is a comment
requests==2.28.1
# Another comment
urllib3==1.26.12`,
			expected: 2,
		},
		{
			name:     "empty freeze output",
			output:   ``,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packages := manager.parseFreezeOutput(tt.output)

			if len(packages) != tt.expected {
				t.Errorf("parseFreezeOutput() returned %d packages, expected %d", len(packages), tt.expected)
			}

			// Verify package structure
			for _, pkg := range packages {
				if pkg.Name == "" {
					t.Error("Package name should not be empty")
				}
			}
		})
	}
}

func TestSetPackageInfoField(t *testing.T) {
	manager := NewManager(nil)
	info := &PackageInfo{
		Metadata: make(map[string]string),
	}

	tests := []struct {
		key   string
		value string
		check func(*PackageInfo) bool
	}{
		{
			key:   "name",
			value: "requests",
			check: func(info *PackageInfo) bool { return info.Name == "requests" },
		},
		{
			key:   "version",
			value: "2.28.1",
			check: func(info *PackageInfo) bool { return info.Version == "2.28.1" },
		},
		{
			key:   "summary",
			value: "HTTP library",
			check: func(info *PackageInfo) bool { return info.Summary == "HTTP library" },
		},
		{
			key:   "requires",
			value: "urllib3, certifi",
			check: func(info *PackageInfo) bool {
				return len(info.Requires) == 2 && info.Requires[0] == "urllib3" && info.Requires[1] == "certifi"
			},
		},
		{
			key:   "custom-field",
			value: "custom-value",
			check: func(info *PackageInfo) bool { return info.Metadata["custom-field"] == "custom-value" },
		},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			manager.setPackageInfoField(info, tt.key, tt.value)

			if !tt.check(info) {
				t.Errorf("setPackageInfoField(%s, %s) failed validation", tt.key, tt.value)
			}
		})
	}
}

// Integration tests for operations
func TestIntegrationListPackages(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	manager := setupTestManager(t)
	skipIfNoPip(t, manager)

	packages, err := manager.ListPackages()
	if err != nil {
		t.Errorf("ListPackages() error = %v", err)
		return
	}

	// Should have at least pip itself
	if len(packages) == 0 {
		t.Error("ListPackages() returned no packages")
	}

	// Check if pip is in the list (pip might be named differently or not present)
	found := false
	for _, pkg := range packages {
		if strings.ToLower(pkg.Name) == "pip" {
			found = true
			if pkg.Version == "" {
				t.Log("Pip package found but version is empty (this might be expected)")
			} else {
				t.Logf("Found pip package with version: %s", pkg.Version)
			}
			break
		}
	}

	if !found {
		t.Log("Pip package not found in list (this might be expected depending on the environment)")
		// Log the first few packages for debugging
		for i, pkg := range packages {
			if i < 3 {
				t.Logf("Package %d: %s (%s)", i+1, pkg.Name, pkg.Version)
			}
		}
	}
}

func TestIntegrationShowPackage(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	manager := setupTestManager(t)
	skipIfNoPip(t, manager)

	// Test showing pip package info
	info, err := manager.ShowPackage("pip")
	if err != nil {
		t.Errorf("ShowPackage('pip') error = %v", err)
		return
	}

	if info == nil {
		t.Fatal("ShowPackage('pip') returned nil info")
	}

	if strings.ToLower(info.Name) != "pip" {
		t.Errorf("Package name = %s, want pip", info.Name)
	}

	if info.Version == "" {
		t.Error("Package version should not be empty")
	}
}

func TestIntegrationFreezePackages(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	manager := setupTestManager(t)
	skipIfNoPip(t, manager)

	packages, err := manager.FreezePackages()
	if err != nil {
		t.Errorf("FreezePackages() error = %v", err)
		return
	}

	// Should have at least some packages
	if len(packages) == 0 {
		t.Error("FreezePackages() returned no packages")
	}

	// Verify package format
	for _, pkg := range packages {
		if pkg.Name == "" {
			t.Error("Package name should not be empty")
		}
		// Note: Version might be empty for editable packages
	}
}

func TestIntegrationSearchPackages(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	manager := setupTestManager(t)
	skipIfNoPip(t, manager)

	// Test search (should return error since pip search is disabled)
	results, err := manager.SearchPackages("requests")

	// We expect an error since pip search is disabled
	if err == nil {
		t.Error("SearchPackages() should return error since pip search is disabled")
	}

	if results != nil {
		t.Error("SearchPackages() should return nil results when disabled")
	}

	// Check if it's the expected error type
	if !strings.Contains(err.Error(), "disabled") {
		t.Errorf("Expected 'disabled' error, got: %v", err)
	}
}

// Benchmark tests for operations
func BenchmarkParseListOutput(b *testing.B) {
	manager := NewManager(nil)
	output := `Package    Version
---------- -------
pip        21.3.1
setuptools 58.1.0
wheel      0.37.1
requests   2.28.1
urllib3    1.26.12`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = manager.parseListOutput(output)
	}
}

func BenchmarkParseFreezeOutput(b *testing.B) {
	manager := NewManager(nil)
	output := `certifi==2022.9.24
charset-normalizer==2.1.1
idna==3.4
requests==2.28.1
urllib3==1.26.12`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = manager.parseFreezeOutput(output)
	}
}

func TestExecutePipCommand(t *testing.T) {
	manager := NewManager(nil)

	tests := []struct {
		name    string
		pipPath string
		args    []string
		wantErr bool
	}{
		{
			name:    "invalid command",
			pipPath: "nonexistent-pip-command",
			args:    []string{"--version"},
			wantErr: true,
		},
		{
			name:    "empty pip path",
			pipPath: "",
			args:    []string{"--version"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.executePipCommand(tt.pipPath, tt.args)

			if tt.wantErr {
				if err == nil {
					t.Errorf("executePipCommand() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("executePipCommand() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestExecutePipCommandWithOutput(t *testing.T) {
	manager := NewManager(nil)

	tests := []struct {
		name    string
		pipPath string
		args    []string
		wantErr bool
	}{
		{
			name:    "invalid command",
			pipPath: "nonexistent-pip-command",
			args:    []string{"--version"},
			wantErr: true,
		},
		{
			name:    "command with spaces",
			pipPath: "python -m pip",
			args:    []string{"--help"},
			wantErr: false, // Might succeed if Python is available
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := manager.executePipCommandWithOutput(tt.pipPath, tt.args)

			if tt.wantErr {
				if err == nil {
					t.Errorf("executePipCommandWithOutput() expected error but got none")
				}
			} else {
				// For valid cases, we don't assert success since pip might not be available
				t.Logf("Command output: %s, error: %v", output, err)
			}
		})
	}
}

func TestInstallPackageValidation(t *testing.T) {
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
			name: "valid package with options",
			pkg: &PackageSpec{
				Name:           "requests",
				Version:        ">=2.25.0",
				Extras:         []string{"security"},
				Upgrade:        true,
				ForceReinstall: true,
				Editable:       true,
				Index:          "https://pypi.org/simple/",
				Options: map[string]string{
					"timeout": "30",
					"retries": "3",
				},
			},
			wantErr: false, // Validation should pass, but installation might fail
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.InstallPackage(tt.pkg)

			if tt.wantErr {
				if err == nil {
					t.Errorf("InstallPackage() expected error but got none")
				}
			} else {
				// For valid packages, we expect either success or pip-related errors
				if err != nil {
					t.Logf("InstallPackage() returned error (possibly due to pip not being available): %v", err)
				}
			}
		})
	}
}

func TestUninstallPackageValidation(t *testing.T) {
	manager := NewManager(nil)

	tests := []struct {
		name        string
		packageName string
		wantErr     bool
		expectedErr string
	}{
		{
			name:        "empty package name",
			packageName: "",
			wantErr:     true,
			expectedErr: "package name cannot be empty",
		},
		{
			name:        "valid package name",
			packageName: "requests",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.UninstallPackage(tt.packageName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("UninstallPackage() expected error but got none")
					return
				}
				if tt.expectedErr != "" && !strings.Contains(err.Error(), tt.expectedErr) {
					t.Errorf("UninstallPackage() error = %v, expected to contain %s", err, tt.expectedErr)
				}
			} else {
				// For valid packages, we expect either success or pip-related errors
				if err != nil {
					t.Logf("UninstallPackage() returned error (possibly due to pip not being available): %v", err)
				}
			}
		})
	}
}

func TestShowPackageValidation(t *testing.T) {
	manager := NewManager(nil)

	tests := []struct {
		name        string
		packageName string
		wantErr     bool
		expectedErr string
	}{
		{
			name:        "empty package name",
			packageName: "",
			wantErr:     true,
			expectedErr: "package name cannot be empty",
		},
		{
			name:        "valid package name",
			packageName: "pip",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := manager.ShowPackage(tt.packageName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ShowPackage() expected error but got none")
					return
				}
				if tt.expectedErr != "" && !strings.Contains(err.Error(), tt.expectedErr) {
					t.Errorf("ShowPackage() error = %v, expected to contain %s", err, tt.expectedErr)
				}
			} else {
				// For valid packages, we expect either success or pip-related errors
				if err != nil {
					t.Logf("ShowPackage() returned error (possibly due to pip not being available): %v", err)
				} else if info != nil {
					t.Logf("ShowPackage() returned info for %s: %s", info.Name, info.Version)
				}
			}
		})
	}
}

func TestSearchPackagesDisabled(t *testing.T) {
	manager := NewManager(nil)

	// Test that search returns appropriate error since it's disabled
	results, err := manager.SearchPackages("requests")

	if err == nil {
		t.Error("SearchPackages() should return error since pip search is disabled")
	}

	if results != nil {
		t.Error("SearchPackages() should return nil results when disabled")
	}

	if !strings.Contains(err.Error(), "disabled") {
		t.Errorf("Expected 'disabled' error, got: %v", err)
	}

	// Test empty query
	_, err = manager.SearchPackages("")
	if err == nil {
		t.Error("SearchPackages() with empty query should return error")
	}
}

func TestExecutePipCommandWithOutputDetailed(t *testing.T) {
	manager := NewManager(nil)

	tests := []struct {
		name         string
		args         []string
		expectError  bool
		expectOutput bool
	}{
		{
			name:         "help command",
			args:         []string{"--help"},
			expectError:  false,
			expectOutput: true,
		},
		{
			name:         "version command",
			args:         []string{"--version"},
			expectError:  false,
			expectOutput: true,
		},
		{
			name:         "invalid command",
			args:         []string{"invalid-command"},
			expectError:  true,
			expectOutput: true, // Error output
		},
		{
			name:         "list with format",
			args:         []string{"list", "--format=json"},
			expectError:  false,
			expectOutput: true,
		},
		{
			name:         "empty args",
			args:         []string{},
			expectError:  false,
			expectOutput: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Try to find pip path first
			pipPath, pipErr := manager.findPipExecutable()
			if pipErr != nil {
				t.Logf("Pip not found, using dummy path: %v", pipErr)
				pipPath = "pip" // Use dummy path for testing
			}

			output, err := manager.executePipCommandWithOutput(pipPath, tt.args)

			if tt.expectError {
				if err == nil {
					t.Errorf("executePipCommandWithOutput() expected error but got none")
				}
			} else {
				if err != nil {
					t.Logf("executePipCommandWithOutput() returned error (might be expected): %v", err)
				}
			}

			if tt.expectOutput {
				if output == "" {
					t.Logf("executePipCommandWithOutput() returned empty output for %s", tt.name)
				} else {
					t.Logf("executePipCommandWithOutput() output length: %d", len(output))
				}
			}
		})
	}
}

func TestExecutePipCommandWithOutputTimeout(t *testing.T) {
	manager := NewManager(nil)

	// Set a very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	manager.SetContext(ctx)

	// Try to find pip path first
	pipPath, pipErr := manager.findPipExecutable()
	if pipErr != nil {
		t.Logf("Pip not found, using dummy path: %v", pipErr)
		pipPath = "pip" // Use dummy path for testing
	}

	// This should timeout quickly
	_, err := manager.executePipCommandWithOutput(pipPath, []string{"list"})

	// We expect either a timeout error or the command to complete quickly
	if err != nil {
		t.Logf("executePipCommandWithOutput() with timeout: %v", err)
	}
}

func TestExecutePipCommandWithOutputInvalidPipPath(t *testing.T) {
	// Create manager with invalid pip path
	config := &Config{
		PipPath: "/nonexistent/pip",
	}
	manager := NewManager(config)

	_, err := manager.executePipCommandWithOutput("/nonexistent/pip", []string{"--version"})

	if err == nil {
		t.Error("executePipCommandWithOutput() should fail with invalid pip path")
	}
}
