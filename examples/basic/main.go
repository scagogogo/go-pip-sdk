package main

import (
	"fmt"
	"log"

	"github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func main() {
	fmt.Println("=== Go Pip SDK Basic Example ===")

	// Create a new pip manager with default configuration
	manager := pip.NewManager(nil)

	// Check if pip is installed
	fmt.Println("\n1. Checking if pip is installed...")
	installed, err := manager.IsInstalled()
	if err != nil {
		log.Printf("Error checking pip installation: %v", err)
	}

	if installed {
		fmt.Println("✓ Pip is installed")

		// Get pip version
		version, err := manager.GetVersion()
		if err != nil {
			log.Printf("Error getting pip version: %v", err)
		} else {
			fmt.Printf("✓ Pip version: %s", version)
		}
	} else {
		fmt.Println("✗ Pip is not installed")
		fmt.Println("Installing pip...")

		if err := manager.Install(); err != nil {
			log.Fatalf("Failed to install pip: %v", err)
		}

		fmt.Println("✓ Pip installed successfully")
	}

	// List currently installed packages
	fmt.Println("\n2. Listing installed packages...")
	packages, err := manager.ListPackages()
	if err != nil {
		log.Printf("Error listing packages: %v", err)
	} else {
		fmt.Printf("Found %d installed packages:\n", len(packages))
		for i, pkg := range packages {
			if i < 5 { // Show only first 5 packages
				fmt.Printf("  - %s (%s)\n", pkg.Name, pkg.Version)
			}
		}
		if len(packages) > 5 {
			fmt.Printf("  ... and %d more packages\n", len(packages)-5)
		}
	}

	// Show detailed information about a specific package (pip itself)
	fmt.Println("\n3. Showing package information...")
	info, err := manager.ShowPackage("pip")
	if err != nil {
		log.Printf("Error showing package info: %v", err)
	} else {
		fmt.Printf("Package: %s\n", info.Name)
		fmt.Printf("Version: %s\n", info.Version)
		fmt.Printf("Summary: %s\n", info.Summary)
		fmt.Printf("Location: %s\n", info.Location)
		if len(info.Requires) > 0 {
			fmt.Printf("Requires: %v\n", info.Requires)
		}
	}

	// Install a package (requests is a popular HTTP library)
	fmt.Println("\n4. Installing a package...")
	pkg := &pip.PackageSpec{
		Name:    "requests",
		Version: ">=2.25.0",
	}

	fmt.Printf("Installing %s%s...\n", pkg.Name, pkg.Version)
	err = manager.InstallPackage(pkg)
	if err != nil {
		// Check if it's already installed
		if pip.IsErrorType(err, pip.ErrorTypeCommandFailed) {
			fmt.Printf("Package might already be installed: %v\n", err)
		} else {
			log.Printf("Error installing package: %v", err)
		}
	} else {
		fmt.Println("✓ Package installed successfully")
	}

	// Verify the installation by showing package info
	fmt.Println("\n5. Verifying installation...")
	requestsInfo, err := manager.ShowPackage("requests")
	if err != nil {
		log.Printf("Error verifying installation: %v", err)
	} else {
		fmt.Printf("✓ Requests %s is installed at %s\n", requestsInfo.Version, requestsInfo.Location)
	}

	// Freeze packages (equivalent to pip freeze)
	fmt.Println("\n6. Freezing packages...")
	frozenPackages, err := manager.FreezePackages()
	if err != nil {
		log.Printf("Error freezing packages: %v", err)
	} else {
		fmt.Printf("Frozen %d packages (showing first 5):\n", len(frozenPackages))
		for i, pkg := range frozenPackages {
			if i < 5 {
				if pkg.Version != "" {
					fmt.Printf("  %s==%s\n", pkg.Name, pkg.Version)
				} else {
					fmt.Printf("  %s\n", pkg.Name)
				}
			}
		}
		if len(frozenPackages) > 5 {
			fmt.Printf("  ... and %d more packages\n", len(frozenPackages)-5)
		}
	}

	// Demonstrate error handling
	fmt.Println("\n7. Demonstrating error handling...")
	nonExistentPkg := &pip.PackageSpec{
		Name: "this-package-definitely-does-not-exist-12345",
	}

	err = manager.InstallPackage(nonExistentPkg)
	if err != nil {
		fmt.Printf("Expected error occurred: %v\n", err)

		// Check error type and provide appropriate response
		switch pip.GetErrorType(err) {
		case pip.ErrorTypePackageNotFound:
			fmt.Println("→ This is a package not found error")
		case pip.ErrorTypeCommandFailed:
			fmt.Println("→ This is a command execution error")
		case pip.ErrorTypeNetworkError:
			fmt.Println("→ This is a network error")
		default:
			fmt.Println("→ This is an unknown error type")
		}
	}

	fmt.Println("\n=== Example completed ===")
}
