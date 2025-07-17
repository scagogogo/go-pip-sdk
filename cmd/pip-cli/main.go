package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/scagogogo/go-pip-sdk/pkg/pip"
)

const (
	version = "1.0.0"
	usage   = `pip-cli - A command-line interface for the Go Pip SDK

Usage:
  pip-cli [global options] command [command options] [arguments...]

Commands:
  install     Install a package
  uninstall   Uninstall a package
  list        List installed packages
  show        Show package information
  freeze      Output installed packages in requirements format
  venv        Virtual environment operations
  project     Project operations
  version     Show version information
  help        Show help

Global Options:
  -timeout duration    Timeout for operations (default: 30s)
  -verbose            Enable verbose logging
  -python string      Path to Python executable
  -pip string         Path to pip executable

Examples:
  pip-cli install requests
  pip-cli install requests click flask
  pip-cli install "requests>=2.25.0"
  pip-cli venv create ./myenv
  pip-cli project init ./myproject
  pip-cli list
  pip-cli show requests

For more information about a command, use: pip-cli help <command>
`
)

var (
	timeoutFlag = flag.Duration("timeout", 30*time.Second, "Timeout for operations")
	verboseFlag = flag.Bool("verbose", false, "Enable verbose logging")
	pythonFlag  = flag.String("python", "", "Path to Python executable")
	pipFlag     = flag.String("pip", "", "Path to pip executable")
)

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	command := flag.Arg(0)
	args := flag.Args()[1:]

	// Create pip manager with configuration
	config := &pip.Config{
		Timeout:     *timeoutFlag,
		PythonPath:  *pythonFlag,
		PipPath:     *pipFlag,
		Environment: make(map[string]string),
	}

	if *verboseFlag {
		config.LogLevel = "DEBUG"
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeoutFlag)
	defer cancel()

	manager := pip.NewManagerWithContext(ctx, config)

	// Set up custom logger if verbose
	if *verboseFlag {
		loggerConfig := &pip.LoggerConfig{
			Level:  pip.LogLevelDebug,
			Output: os.Stdout,
			Prefix: "[pip-cli]",
		}
		logger, err := pip.NewLogger(loggerConfig)
		if err != nil {
			log.Fatalf("Failed to create logger: %v", err)
		}
		defer logger.Close()
		manager.SetCustomLogger(logger)
	}

	// Execute command
	switch command {
	case "install":
		handleInstall(manager, args)
	case "uninstall":
		handleUninstall(manager, args)
	case "list":
		handleList(manager, args)
	case "show":
		handleShow(manager, args)
	case "freeze":
		handleFreeze(manager, args)
	case "venv":
		handleVenv(manager, args)
	case "project":
		handleProject(manager, args)
	case "version":
		handleVersion(args)
	case "help":
		handleHelp(args)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		flag.Usage()
		os.Exit(1)
	}
}

// isVersionSpec checks if a string looks like a version specification
func isVersionSpec(s string) bool {
	// Check for version operators
	versionOperators := []string{">=", "==", "!=", "<=", "~=", "===", ">", "<"}
	for _, op := range versionOperators {
		if strings.Contains(s, op) {
			return true
		}
	}

	// Check if it starts with a digit and contains dots (like "2.25.0")
	if len(s) > 0 && s[0] >= '0' && s[0] <= '9' && strings.Contains(s, ".") {
		return true
	}

	// Check for version ranges with commas
	if strings.Contains(s, ",") {
		return true
	}

	return false
}

func handleInstall(manager pip.PipManager, args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: package name required\n")
		fmt.Fprintf(os.Stderr, "Usage: pip-cli install <package1> [package2] ...\n")
		fmt.Fprintf(os.Stderr, "       pip-cli install <package> <version>\n")
		os.Exit(1)
	}

	// Handle multiple packages or single package with version
	var packages []*pip.PackageSpec

	if len(args) == 2 && isVersionSpec(args[1]) {
		// Package + version format: pip-cli install requests ">=2.25.0"
		pkg := &pip.PackageSpec{
			Name:    args[0],
			Version: args[1],
		}
		packages = append(packages, pkg)
	} else {
		// Multiple packages: pip-cli install requests click flask
		for _, packageName := range args {
			if packageName == "" {
				continue
			}
			pkg := &pip.PackageSpec{
				Name: packageName,
			}
			packages = append(packages, pkg)
		}
	}

	if len(packages) == 1 {
		versionStr := ""
		if packages[0].Version != "" {
			versionStr = packages[0].Version
		}
		fmt.Printf("Installing %s%s...\n", packages[0].Name, versionStr)
	} else {
		packageNames := make([]string, len(packages))
		for i, pkg := range packages {
			packageNames[i] = pkg.Name
		}
		fmt.Printf("Installing %d packages: %s...\n", len(packages), strings.Join(packageNames, ", "))
	}

	start := time.Now()
	var errors []string
	successCount := 0

	for _, pkg := range packages {
		err := manager.InstallPackage(pkg)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", pkg.Name, err))
		} else {
			successCount++
		}
	}

	duration := time.Since(start)

	if len(errors) > 0 {
		if successCount > 0 {
			fmt.Printf("✓ %d packages installed successfully\n", successCount)
		}
		fmt.Fprintf(os.Stderr, "Installation failed for %d packages:\n", len(errors))
		for _, errMsg := range errors {
			fmt.Fprintf(os.Stderr, "  - %s\n", errMsg)
		}
		os.Exit(1)
	}

	if len(packages) == 1 {
		fmt.Printf("✓ Package %s installed successfully (took %v)\n", packages[0].Name, duration)
	} else {
		fmt.Printf("✓ All %d packages installed successfully (took %v)\n", len(packages), duration)
	}
}

func handleUninstall(manager pip.PipManager, args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: package name required\n")
		fmt.Fprintf(os.Stderr, "Usage: pip-cli uninstall <package>\n")
		os.Exit(1)
	}

	packageName := args[0]
	fmt.Printf("Uninstalling %s...\n", packageName)

	start := time.Now()
	err := manager.UninstallPackage(packageName)
	duration := time.Since(start)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Uninstallation failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Package %s uninstalled successfully (took %v)\n", packageName, duration)
}

func handleList(manager pip.PipManager, args []string) {
	fmt.Println("Listing installed packages...")

	packages, err := manager.ListPackages()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to list packages: %v\n", err)
		os.Exit(1)
	}

	if len(packages) == 0 {
		fmt.Println("No packages installed")
		return
	}

	fmt.Printf("Found %d installed packages:\n\n", len(packages))
	fmt.Printf("%-30s %-15s %s\n", "Package", "Version", "Location")
	fmt.Printf("%s\n", strings.Repeat("-", 80))

	for _, pkg := range packages {
		location := pkg.Location
		if location == "" {
			location = "N/A"
		}
		fmt.Printf("%-30s %-15s %s\n", pkg.Name, pkg.Version, location)
	}
}

func handleShow(manager pip.PipManager, args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: package name required\n")
		fmt.Fprintf(os.Stderr, "Usage: pip-cli show <package>\n")
		os.Exit(1)
	}

	packageName := args[0]
	fmt.Printf("Showing information for %s...\n\n", packageName)

	info, err := manager.ShowPackage(packageName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to show package info: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Name: %s\n", info.Name)
	fmt.Printf("Version: %s\n", info.Version)
	if info.Summary != "" {
		fmt.Printf("Summary: %s\n", info.Summary)
	}
	if info.HomePage != "" {
		fmt.Printf("Home-page: %s\n", info.HomePage)
	}
	if info.Author != "" {
		fmt.Printf("Author: %s\n", info.Author)
	}
	if info.AuthorEmail != "" {
		fmt.Printf("Author-email: %s\n", info.AuthorEmail)
	}
	if info.License != "" {
		fmt.Printf("License: %s\n", info.License)
	}
	if info.Location != "" {
		fmt.Printf("Location: %s\n", info.Location)
	}
	if len(info.Requires) > 0 {
		fmt.Printf("Requires: %s\n", strings.Join(info.Requires, ", "))
	}
	if len(info.RequiredBy) > 0 {
		fmt.Printf("Required-by: %s\n", strings.Join(info.RequiredBy, ", "))
	}
}

func handleFreeze(manager pip.PipManager, args []string) {
	fmt.Println("Freezing packages...")

	packages, err := manager.FreezePackages()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to freeze packages: %v\n", err)
		os.Exit(1)
	}

	for _, pkg := range packages {
		if pkg.Version != "" {
			fmt.Printf("%s==%s\n", pkg.Name, pkg.Version)
		} else {
			fmt.Printf("%s\n", pkg.Name)
		}
	}
}

func handleVenv(manager pip.PipManager, args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: venv subcommand required\n")
		fmt.Fprintf(os.Stderr, "Usage: pip-cli venv <create|activate|deactivate|remove|info> [path]\n")
		os.Exit(1)
	}

	subcommand := args[0]
	subargs := args[1:]

	switch subcommand {
	case "create":
		if len(subargs) == 0 {
			fmt.Fprintf(os.Stderr, "Error: path required for venv create\n")
			os.Exit(1)
		}
		path := subargs[0]
		fmt.Printf("Creating virtual environment at %s...\n", path)

		err := manager.CreateVenv(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create virtual environment: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Virtual environment created successfully\n")

	case "activate":
		if len(subargs) == 0 {
			fmt.Fprintf(os.Stderr, "Error: path required for venv activate\n")
			os.Exit(1)
		}
		path := subargs[0]
		fmt.Printf("Activating virtual environment at %s...\n", path)

		err := manager.ActivateVenv(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to activate virtual environment: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Virtual environment activated\n")

	case "deactivate":
		fmt.Println("Deactivating virtual environment...")

		err := manager.DeactivateVenv()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to deactivate virtual environment: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Virtual environment deactivated\n")

	case "remove":
		if len(subargs) == 0 {
			fmt.Fprintf(os.Stderr, "Error: path required for venv remove\n")
			os.Exit(1)
		}
		path := subargs[0]
		fmt.Printf("Removing virtual environment at %s...\n", path)

		err := manager.RemoveVenv(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to remove virtual environment: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Virtual environment removed\n")

	case "info":
		if len(subargs) == 0 {
			fmt.Fprintf(os.Stderr, "Error: path required for venv info\n")
			os.Exit(1)
		}
		path := subargs[0]

		// Check if GetVenvInfo method exists (it might not be in the interface)
		fmt.Printf("Virtual environment path: %s\n", path)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Printf("Status: Not found\n")
		} else {
			fmt.Printf("Status: Exists\n")
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown venv subcommand: %s\n", subcommand)
		os.Exit(1)
	}
}

func handleProject(manager pip.PipManager, args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: project subcommand required\n")
		fmt.Fprintf(os.Stderr, "Usage: pip-cli project <init> <path> [options]\n")
		os.Exit(1)
	}

	subcommand := args[0]
	subargs := args[1:]

	switch subcommand {
	case "init":
		if len(subargs) == 0 {
			fmt.Fprintf(os.Stderr, "Error: path required for project init\n")
			os.Exit(1)
		}
		path := subargs[0]

		projectName := filepath.Base(path)
		opts := &pip.ProjectOptions{
			Name:        projectName,
			Version:     "0.1.0",
			Description: fmt.Sprintf("A Python project: %s", projectName),
			Author:      "Developer",
			License:     "MIT",
			CreateVenv:  true,
			VenvPath:    filepath.Join(path, "venv"),
		}

		fmt.Printf("Initializing project at %s...\n", path)

		err := manager.InitProject(path, opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to initialize project: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Project initialized successfully\n")

	default:
		fmt.Fprintf(os.Stderr, "Unknown project subcommand: %s\n", subcommand)
		os.Exit(1)
	}
}

func handleVersion(args []string) {
	fmt.Printf("pip-cli version %s\n", version)
	fmt.Printf("Go Pip SDK - A Go library for managing Python pip operations\n")
}

func handleHelp(args []string) {
	if len(args) == 0 {
		flag.Usage()
		return
	}

	command := args[0]
	switch command {
	case "install":
		fmt.Println("Install Python packages")
		fmt.Println("Usage: pip-cli install <package1> [package2] ...")
		fmt.Println("       pip-cli install <package> <version>")
		fmt.Println("Examples:")
		fmt.Println("  pip-cli install requests")
		fmt.Println("  pip-cli install requests click flask")
		fmt.Println("  pip-cli install requests '>=2.25.0'")
	case "uninstall":
		fmt.Println("Uninstall a Python package")
		fmt.Println("Usage: pip-cli uninstall <package>")
	case "list":
		fmt.Println("List installed packages")
		fmt.Println("Usage: pip-cli list")
	case "show":
		fmt.Println("Show package information")
		fmt.Println("Usage: pip-cli show <package>")
	case "freeze":
		fmt.Println("Output installed packages in requirements format")
		fmt.Println("Usage: pip-cli freeze")
	case "venv":
		fmt.Println("Virtual environment operations")
		fmt.Println("Usage: pip-cli venv <create|activate|deactivate|remove|info> [path]")
	case "project":
		fmt.Println("Project operations")
		fmt.Println("Usage: pip-cli project <init> <path>")
	default:
		fmt.Printf("No help available for command: %s\n", command)
	}
}
