package pip

import (
	"time"
)

// OSType represents the operating system type
type OSType string

const (
	OSWindows OSType = "windows"
	OSMacOS   OSType = "darwin"
	OSLinux   OSType = "linux"
	OSUnknown OSType = "unknown"
)

// PipManager represents the main pip manager interface
type PipManager interface {
	// System operations
	IsInstalled() (bool, error)
	Install() error
	GetVersion() (string, error)

	// Package operations
	InstallPackage(pkg *PackageSpec) error
	UninstallPackage(name string) error
	ListPackages() ([]*Package, error)
	ShowPackage(name string) (*PackageInfo, error)
	SearchPackages(query string) ([]*SearchResult, error)
	FreezePackages() ([]*Package, error)

	// Virtual environment operations
	CreateVenv(path string) error
	ActivateVenv(path string) error
	DeactivateVenv() error
	RemoveVenv(path string) error

	// Project operations
	InitProject(path string, opts *ProjectOptions) error
	InstallRequirements(path string) error
	GenerateRequirements(path string) error
}

// PackageSpec represents a package specification for installation
type PackageSpec struct {
	Name           string            `json:"name"`
	Version        string            `json:"version,omitempty"`  // e.g., ">=1.0.0", "==2.1.0"
	Extras         []string          `json:"extras,omitempty"`   // e.g., ["dev", "test"]
	Index          string            `json:"index,omitempty"`    // custom index URL
	Options        map[string]string `json:"options,omitempty"`  // additional pip options
	Editable       bool              `json:"editable,omitempty"` // install in editable mode
	Upgrade        bool              `json:"upgrade,omitempty"`  // upgrade if already installed
	ForceReinstall bool              `json:"force_reinstall,omitempty"`
}

// Package represents an installed package
type Package struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	Location  string `json:"location,omitempty"`
	Editable  bool   `json:"editable,omitempty"`
	Installer string `json:"installer,omitempty"`
}

// PackageInfo represents detailed package information
type PackageInfo struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Summary     string            `json:"summary,omitempty"`
	HomePage    string            `json:"home_page,omitempty"`
	Author      string            `json:"author,omitempty"`
	AuthorEmail string            `json:"author_email,omitempty"`
	License     string            `json:"license,omitempty"`
	Location    string            `json:"location,omitempty"`
	Requires    []string          `json:"requires,omitempty"`
	RequiredBy  []string          `json:"required_by,omitempty"`
	Files       []string          `json:"files,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// SearchResult represents a package search result
type SearchResult struct {
	Name    string  `json:"name"`
	Version string  `json:"version"`
	Summary string  `json:"summary,omitempty"`
	Score   float64 `json:"score,omitempty"`
}

// ProjectOptions represents options for project initialization
type ProjectOptions struct {
	Name            string            `json:"name"`
	Version         string            `json:"version,omitempty"`
	Description     string            `json:"description,omitempty"`
	Author          string            `json:"author,omitempty"`
	AuthorEmail     string            `json:"author_email,omitempty"`
	License         string            `json:"license,omitempty"`
	PythonVersion   string            `json:"python_version,omitempty"`
	Dependencies    []string          `json:"dependencies,omitempty"`
	DevDependencies []string          `json:"dev_dependencies,omitempty"`
	CreateVenv      bool              `json:"create_venv,omitempty"`
	VenvPath        string            `json:"venv_path,omitempty"`
	Template        string            `json:"template,omitempty"`
	ExtraFiles      map[string]string `json:"extra_files,omitempty"`
}

// VenvInfo represents virtual environment information
type VenvInfo struct {
	Path          string    `json:"path"`
	PythonPath    string    `json:"python_path"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	PythonVersion string    `json:"python_version,omitempty"`
}

// InstallResult represents the result of a package installation
type InstallResult struct {
	Package     *PackageSpec  `json:"package"`
	Success     bool          `json:"success"`
	Message     string        `json:"message,omitempty"`
	Error       error         `json:"error,omitempty"`
	Duration    time.Duration `json:"duration"`
	OutputLines []string      `json:"output_lines,omitempty"`
}

// UninstallResult represents the result of a package uninstallation
type UninstallResult struct {
	PackageName string        `json:"package_name"`
	Success     bool          `json:"success"`
	Message     string        `json:"message,omitempty"`
	Error       error         `json:"error,omitempty"`
	Duration    time.Duration `json:"duration"`
	OutputLines []string      `json:"output_lines,omitempty"`
}

// Config represents pip manager configuration
type Config struct {
	PythonPath   string            `json:"python_path,omitempty"`
	PipPath      string            `json:"pip_path,omitempty"`
	DefaultIndex string            `json:"default_index,omitempty"`
	TrustedHosts []string          `json:"trusted_hosts,omitempty"`
	Timeout      time.Duration     `json:"timeout,omitempty"`
	Retries      int               `json:"retries,omitempty"`
	LogLevel     string            `json:"log_level,omitempty"`
	CacheDir     string            `json:"cache_dir,omitempty"`
	ExtraOptions map[string]string `json:"extra_options,omitempty"`
	Environment  map[string]string `json:"environment,omitempty"`
}

// Error types
type PipError struct {
	Type     string `json:"type"`
	Message  string `json:"message"`
	Command  string `json:"command,omitempty"`
	Output   string `json:"output,omitempty"`
	ExitCode int    `json:"exit_code,omitempty"`
}

func (e *PipError) Error() string {
	return e.Message
}

// Note: Common error instances are defined in errors.go
