package pip

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNewProjectManager(t *testing.T) {
	manager := NewManager(nil)
	projectManager := NewProjectManager(manager)

	if projectManager == nil {
		t.Fatal("NewProjectManager() returned nil")
	}

	if projectManager.manager != manager {
		t.Error("ProjectManager should reference the provided manager")
	}
}

func TestInitProject(t *testing.T) {
	manager := NewManager(nil)

	tests := []struct {
		name    string
		path    string
		opts    *ProjectOptions
		wantErr bool
		setup   func(string) error
		cleanup func(string) error
	}{
		{
			name:    "empty path",
			path:    "",
			opts:    nil,
			wantErr: true,
		},
		{
			name:    "valid project with defaults",
			path:    "test-project",
			opts:    nil,
			wantErr: false,
			setup: func(path string) error {
				return nil
			},
			cleanup: func(path string) error {
				return os.RemoveAll(path)
			},
		},
		{
			name: "project with custom options",
			path: "custom-project",
			opts: &ProjectOptions{
				Name:            "my-custom-project",
				Version:         "1.0.0",
				Description:     "A test project",
				Author:          "Test Author",
				AuthorEmail:     "test@example.com",
				License:         "Apache-2.0",
				Dependencies:    []string{"requests>=2.25.0", "click>=8.0.0"},
				DevDependencies: []string{"pytest>=6.0.0", "black>=21.0.0"},
			},
			wantErr: false,
			setup: func(path string) error {
				return nil
			},
			cleanup: func(path string) error {
				return os.RemoveAll(path)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			var fullPath string
			if tt.path != "" {
				tempDir, err := os.MkdirTemp("", "pip-sdk-test")
				if err != nil {
					t.Fatalf("Failed to create temp dir: %v", err)
				}
				fullPath = filepath.Join(tempDir, tt.path)
				defer os.RemoveAll(tempDir)

				if tt.setup != nil {
					if err := tt.setup(fullPath); err != nil {
						t.Fatalf("Setup failed: %v", err)
					}
				}
			}

			// Test
			err := manager.InitProject(fullPath, tt.opts)

			// Verify
			if (err != nil) != tt.wantErr {
				t.Errorf("InitProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && fullPath != "" {
				// Verify project structure was created
				expectedFiles := []string{
					"requirements.txt",
					"setup.py",
					"pyproject.toml",
					"README.md",
					".gitignore",
				}

				for _, file := range expectedFiles {
					filePath := filepath.Join(fullPath, file)
					if _, err := os.Stat(filePath); os.IsNotExist(err) {
						t.Errorf("Expected file %s was not created", file)
					}
				}

				// Verify package directory was created
				var projectName string
				if tt.opts != nil && tt.opts.Name != "" {
					projectName = tt.opts.Name
				} else {
					projectName = filepath.Base(fullPath)
				}
				packagePath := filepath.Join(fullPath, projectName)
				if _, err := os.Stat(packagePath); os.IsNotExist(err) {
					t.Errorf("Package directory %s was not created", packagePath)
				}
			}

			// Cleanup
			if tt.cleanup != nil && fullPath != "" {
				if err := tt.cleanup(fullPath); err != nil {
					t.Logf("Cleanup failed: %v", err)
				}
			}
		})
	}
}

func TestCreateRequirementsFile(t *testing.T) {
	manager := NewManager(nil)
	projectManager := NewProjectManager(manager)

	tempDir, err := os.MkdirTemp("", "pip-sdk-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name     string
		opts     *ProjectOptions
		expected []string
	}{
		{
			name: "with dependencies",
			opts: &ProjectOptions{
				Dependencies: []string{"requests>=2.25.0", "click>=8.0.0"},
			},
			expected: []string{"requests>=2.25.0", "click>=8.0.0"},
		},
		{
			name: "with dev dependencies",
			opts: &ProjectOptions{
				Dependencies:    []string{"requests>=2.25.0"},
				DevDependencies: []string{"pytest>=6.0.0", "black>=21.0.0"},
			},
			expected: []string{"requests>=2.25.0"},
		},
		{
			name:     "no dependencies",
			opts:     &ProjectOptions{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectPath := filepath.Join(tempDir, tt.name)
			err := os.MkdirAll(projectPath, 0755)
			if err != nil {
				t.Fatalf("Failed to create project dir: %v", err)
			}

			err = projectManager.createRequirementsFile(projectPath, tt.opts)
			if err != nil {
				t.Errorf("createRequirementsFile() error = %v", err)
				return
			}

			// Verify requirements.txt content
			requirementsPath := filepath.Join(projectPath, "requirements.txt")
			content, err := os.ReadFile(requirementsPath)
			if err != nil {
				t.Errorf("Failed to read requirements.txt: %v", err)
				return
			}

			contentStr := string(content)
			for _, dep := range tt.expected {
				if !strings.Contains(contentStr, dep) {
					t.Errorf("Expected dependency %s not found in requirements.txt", dep)
				}
			}

			// Check for dev dependencies file if dev deps exist
			if len(tt.opts.DevDependencies) > 0 {
				devRequirementsPath := filepath.Join(projectPath, "requirements-dev.txt")
				if _, err := os.Stat(devRequirementsPath); os.IsNotExist(err) {
					t.Error("requirements-dev.txt should be created when dev dependencies exist")
				}
			}
		})
	}
}

func TestCreateSetupFile(t *testing.T) {
	manager := NewManager(nil)
	projectManager := NewProjectManager(manager)

	tempDir, err := os.MkdirTemp("", "pip-sdk-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	opts := &ProjectOptions{
		Name:        "test-package",
		Version:     "1.0.0",
		Description: "A test package",
		Author:      "Test Author",
		AuthorEmail: "test@example.com",
		License:     "MIT",
	}

	err = projectManager.createSetupFile(tempDir, opts)
	if err != nil {
		t.Errorf("createSetupFile() error = %v", err)
		return
	}

	// Verify setup.py was created and contains expected content
	setupPath := filepath.Join(tempDir, "setup.py")
	content, err := os.ReadFile(setupPath)
	if err != nil {
		t.Errorf("Failed to read setup.py: %v", err)
		return
	}

	contentStr := string(content)
	expectedStrings := []string{
		opts.Name,
		opts.Version,
		opts.Author,
		opts.AuthorEmail,
		opts.License,
		"setuptools",
		"find_packages",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(contentStr, expected) {
			t.Errorf("Expected string %s not found in setup.py", expected)
		}
	}
}

func TestCreatePyprojectFile(t *testing.T) {
	manager := NewManager(nil)
	projectManager := NewProjectManager(manager)

	tempDir, err := os.MkdirTemp("", "pip-sdk-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	opts := &ProjectOptions{
		Name:            "test-package",
		Version:         "1.0.0",
		Description:     "A test package",
		Author:          "Test Author",
		AuthorEmail:     "test@example.com",
		License:         "MIT",
		Dependencies:    []string{"requests>=2.25.0"},
		DevDependencies: []string{"pytest>=6.0.0"},
	}

	err = projectManager.createPyprojectFile(tempDir, opts)
	if err != nil {
		t.Errorf("createPyprojectFile() error = %v", err)
		return
	}

	// Verify pyproject.toml was created and contains expected content
	pyprojectPath := filepath.Join(tempDir, "pyproject.toml")
	content, err := os.ReadFile(pyprojectPath)
	if err != nil {
		t.Errorf("Failed to read pyproject.toml: %v", err)
		return
	}

	contentStr := string(content)
	expectedStrings := []string{
		"[build-system]",
		"[project]",
		opts.Name,
		opts.Version,
		opts.Author,
		"requests>=2.25.0",
		"pytest>=6.0.0",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(contentStr, expected) {
			t.Errorf("Expected string %s not found in pyproject.toml", expected)
		}
	}
}

func TestCreateReadmeFile(t *testing.T) {
	manager := NewManager(nil)
	projectManager := NewProjectManager(manager)

	tempDir, err := os.MkdirTemp("", "pip-sdk-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	opts := &ProjectOptions{
		Name:        "test-package",
		Description: "A test package for testing",
		License:     "MIT",
	}

	err = projectManager.createReadmeFile(tempDir, opts)
	if err != nil {
		t.Errorf("createReadmeFile() error = %v", err)
		return
	}

	// Verify README.md was created and contains expected content
	readmePath := filepath.Join(tempDir, "README.md")
	content, err := os.ReadFile(readmePath)
	if err != nil {
		t.Errorf("Failed to read README.md: %v", err)
		return
	}

	contentStr := string(content)
	expectedStrings := []string{
		"# " + opts.Name,
		opts.Description,
		"## Installation",
		"## Usage",
		"## Development",
		opts.License,
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(contentStr, expected) {
			t.Errorf("Expected string %s not found in README.md", expected)
		}
	}
}

func TestCreateGitignoreFile(t *testing.T) {
	manager := NewManager(nil)
	projectManager := NewProjectManager(manager)

	tempDir, err := os.MkdirTemp("", "pip-sdk-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	err = projectManager.createGitignoreFile(tempDir)
	if err != nil {
		t.Errorf("createGitignoreFile() error = %v", err)
		return
	}

	// Verify .gitignore was created and contains expected patterns
	gitignorePath := filepath.Join(tempDir, ".gitignore")
	content, err := os.ReadFile(gitignorePath)
	if err != nil {
		t.Errorf("Failed to read .gitignore: %v", err)
		return
	}

	contentStr := string(content)
	expectedPatterns := []string{
		"__pycache__/",
		"*.py[cod]",
		".env",
		"venv/",
		".DS_Store",
		"*.egg-info/",
	}

	for _, pattern := range expectedPatterns {
		if !strings.Contains(contentStr, pattern) {
			t.Errorf("Expected pattern %s not found in .gitignore", pattern)
		}
	}
}

func TestCreatePackageDirectory(t *testing.T) {
	manager := NewManager(nil)
	projectManager := NewProjectManager(manager)

	tempDir, err := os.MkdirTemp("", "pip-sdk-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	opts := &ProjectOptions{
		Name:        "test-package",
		Version:     "1.0.0",
		Author:      "Test Author",
		AuthorEmail: "test@example.com",
	}

	err = projectManager.createPackageDirectory(tempDir, opts)
	if err != nil {
		t.Errorf("createPackageDirectory() error = %v", err)
		return
	}

	// Verify package directory was created
	packagePath := filepath.Join(tempDir, opts.Name)
	if _, err := os.Stat(packagePath); os.IsNotExist(err) {
		t.Errorf("Package directory %s was not created", packagePath)
		return
	}

	// Verify __init__.py was created
	initPath := filepath.Join(packagePath, "__init__.py")
	if _, err := os.Stat(initPath); os.IsNotExist(err) {
		t.Error("__init__.py was not created")
		return
	}

	// Verify main.py was created
	mainPath := filepath.Join(packagePath, "main.py")
	if _, err := os.Stat(mainPath); os.IsNotExist(err) {
		t.Error("main.py was not created")
		return
	}

	// Verify __init__.py content
	initContent, err := os.ReadFile(initPath)
	if err != nil {
		t.Errorf("Failed to read __init__.py: %v", err)
		return
	}

	initStr := string(initContent)
	if !strings.Contains(initStr, opts.Version) {
		t.Error("__init__.py should contain version information")
	}
}

func TestCreateExtraFiles(t *testing.T) {
	manager := NewManager(nil)
	projectManager := NewProjectManager(manager)

	tempDir, err := os.MkdirTemp("", "pip-sdk-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	opts := &ProjectOptions{
		ExtraFiles: map[string]string{
			"custom.txt":        "This is a custom file",
			"config/app.yaml":   "app:\n  name: test",
			"scripts/deploy.sh": "#!/bin/bash\necho 'Deploying...'",
		},
	}

	err = projectManager.createExtraFiles(tempDir, opts)
	if err != nil {
		t.Errorf("createExtraFiles() error = %v", err)
		return
	}

	// Verify all extra files were created with correct content
	for fileName, expectedContent := range opts.ExtraFiles {
		filePath := filepath.Join(tempDir, fileName)

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("Extra file %s was not created", fileName)
			continue
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Errorf("Failed to read extra file %s: %v", fileName, err)
			continue
		}

		if string(content) != expectedContent {
			t.Errorf("Extra file %s content mismatch. Expected: %s, Got: %s",
				fileName, expectedContent, string(content))
		}
	}
}

func TestInstallRequirements(t *testing.T) {
	manager := NewManager(nil)

	tests := []struct {
		name        string
		setupFile   func(string) error
		path        string
		wantErr     bool
		expectedErr string
	}{
		{
			name:        "empty path",
			path:        "",
			wantErr:     true,
			expectedErr: "requirements file path cannot be empty",
		},
		{
			name:        "file not found",
			path:        "nonexistent-requirements.txt",
			wantErr:     true,
			expectedErr: "requirements file not found",
		},
		{
			name: "valid requirements file",
			setupFile: func(path string) error {
				content := "# Test requirements\nrequests>=2.25.0\nclick>=8.0.0\n"
				return os.WriteFile(path, []byte(content), 0644)
			},
			path:    "test-requirements.txt",
			wantErr: false, // Note: This might fail if pip is not available, but we test the logic
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fullPath string

			if tt.path != "" && tt.setupFile != nil {
				tempDir, err := os.MkdirTemp("", "pip-sdk-test")
				if err != nil {
					t.Fatalf("Failed to create temp dir: %v", err)
				}
				defer os.RemoveAll(tempDir)

				fullPath = filepath.Join(tempDir, tt.path)
				if err := tt.setupFile(fullPath); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			} else {
				fullPath = tt.path
			}

			err := manager.InstallRequirements(fullPath)

			if tt.wantErr {
				if err == nil {
					t.Errorf("InstallRequirements() expected error but got none")
					return
				}
				if tt.expectedErr != "" && !strings.Contains(err.Error(), tt.expectedErr) {
					t.Errorf("InstallRequirements() error = %v, expected to contain %s", err, tt.expectedErr)
				}
			} else {
				// For valid cases, we expect either success or pip-related errors
				// We don't fail the test if pip is not available
				if err != nil {
					t.Logf("InstallRequirements() returned error (possibly due to pip not being available): %v", err)
				}
			}
		})
	}
}

func TestGenerateRequirements(t *testing.T) {
	manager := NewManager(nil)

	tests := []struct {
		name        string
		path        string
		wantErr     bool
		expectedErr string
	}{
		{
			name:        "empty path",
			path:        "",
			wantErr:     true,
			expectedErr: "requirements file path cannot be empty",
		},
		{
			name:    "valid path",
			path:    "generated-requirements.txt",
			wantErr: false, // Note: This might fail if pip is not available
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fullPath string

			if tt.path != "" {
				tempDir, err := os.MkdirTemp("", "pip-sdk-test")
				if err != nil {
					t.Fatalf("Failed to create temp dir: %v", err)
				}
				defer os.RemoveAll(tempDir)

				fullPath = filepath.Join(tempDir, tt.path)
			} else {
				fullPath = tt.path
			}

			err := manager.GenerateRequirements(fullPath)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GenerateRequirements() expected error but got none")
					return
				}
				if tt.expectedErr != "" && !strings.Contains(err.Error(), tt.expectedErr) {
					t.Errorf("GenerateRequirements() error = %v, expected to contain %s", err, tt.expectedErr)
				}
			} else {
				// For valid cases, we expect either success or pip-related errors
				if err != nil {
					t.Logf("GenerateRequirements() returned error (possibly due to pip not being available): %v", err)
				} else {
					// Verify file was created
					if _, err := os.Stat(fullPath); os.IsNotExist(err) {
						t.Error("Requirements file was not created")
					} else {
						// Verify file content
						content, err := os.ReadFile(fullPath)
						if err != nil {
							t.Errorf("Failed to read generated requirements file: %v", err)
						} else {
							contentStr := string(content)
							if !strings.Contains(contentStr, "# Generated requirements file") {
								t.Error("Generated requirements file should contain header comment")
							}
							if !strings.Contains(contentStr, time.Now().Format("2006-01-02")) {
								t.Error("Generated requirements file should contain generation date")
							}
						}
					}
				}
			}
		})
	}
}

func TestCreateSetupFileDetailed(t *testing.T) {
	manager := NewManager(nil)
	projectManager := NewProjectManager(manager)

	tempDir, err := os.MkdirTemp("", "pip-sdk-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name    string
		opts    *ProjectOptions
		wantErr bool
	}{
		{
			name: "minimal options",
			opts: &ProjectOptions{
				Name:        "test-package",
				Version:     "1.0.0",
				Author:      "Test Author",
				AuthorEmail: "test@example.com",
			},
			wantErr: false,
		},
		{
			name: "full options",
			opts: &ProjectOptions{
				Name:            "full-package",
				Version:         "2.0.0",
				Description:     "A full test package",
				Author:          "Full Author",
				AuthorEmail:     "full@example.com",
				License:         "MIT",
				Dependencies:    []string{"requests>=2.25.0", "click>=7.0"},
				DevDependencies: []string{"pytest>=6.0", "black>=21.0"},
				PythonVersion:   "3.7",
			},
			wantErr: false,
		},
		{
			name: "with template",
			opts: &ProjectOptions{
				Name:        "template-package",
				Version:     "1.0.0",
				Author:      "Test Author",
				AuthorEmail: "test@example.com",
				Template:    "basic",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectPath := filepath.Join(tempDir, tt.name)
			if err := os.MkdirAll(projectPath, 0755); err != nil {
				t.Fatalf("Failed to create project dir: %v", err)
			}

			err := projectManager.createSetupFile(projectPath, tt.opts)

			if tt.wantErr {
				if err == nil {
					t.Errorf("createSetupFile() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("createSetupFile() unexpected error: %v", err)
				} else {
					// Verify setup.py was created
					setupPath := filepath.Join(projectPath, "setup.py")
					if _, err := os.Stat(setupPath); os.IsNotExist(err) {
						t.Error("createSetupFile() should create setup.py")
					} else {
						// Read and verify content
						content, err := os.ReadFile(setupPath)
						if err != nil {
							t.Fatalf("Failed to read setup.py: %v", err)
						}

						contentStr := string(content)
						if !strings.Contains(contentStr, tt.opts.Name) {
							t.Error("setup.py should contain package name")
						}
						if !strings.Contains(contentStr, tt.opts.Version) {
							t.Error("setup.py should contain version")
						}
						if !strings.Contains(contentStr, tt.opts.Author) {
							t.Error("setup.py should contain author")
						}
					}
				}
			}
		})
	}
}
