package pip

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// InstallPackage installs a Python package
func (m *Manager) InstallPackage(pkg *PackageSpec) error {
	if err := m.validatePackageSpec(pkg); err != nil {
		return err
	}

	m.logInfo("Installing package: %s", pkg.Name)

	pipPath, err := m.findPipExecutable()
	if err != nil {
		return ErrPipNotInstalled
	}

	// Build command arguments
	args := []string{"install"}

	// Add package specification
	packageSpec := pkg.Name
	if pkg.Version != "" {
		packageSpec += pkg.Version
	}
	if len(pkg.Extras) > 0 {
		packageSpec += "[" + strings.Join(pkg.Extras, ",") + "]"
	}
	args = append(args, packageSpec)

	// Add options
	if pkg.Upgrade {
		args = append(args, "--upgrade")
	}
	if pkg.ForceReinstall {
		args = append(args, "--force-reinstall")
	}
	if pkg.Editable {
		args = append(args, "--editable")
	}
	if pkg.Index != "" {
		args = append(args, "--index-url", pkg.Index)
	}

	// Add extra options
	for key, value := range pkg.Options {
		if value == "" {
			args = append(args, "--"+key)
		} else {
			args = append(args, "--"+key, value)
		}
	}

	// Execute command
	return m.executePipCommand(pipPath, args)
}

// UninstallPackage uninstalls a Python package
func (m *Manager) UninstallPackage(name string) error {
	if name == "" {
		return &PipError{
			Type:    "invalid_package_spec",
			Message: "package name cannot be empty",
		}
	}

	m.logInfo("Uninstalling package: %s", name)

	pipPath, err := m.findPipExecutable()
	if err != nil {
		return ErrPipNotInstalled
	}

	args := []string{"uninstall", "-y", name}
	return m.executePipCommand(pipPath, args)
}

// ListPackages lists all installed packages
func (m *Manager) ListPackages() ([]*Package, error) {
	m.logDebug("Listing installed packages")

	pipPath, err := m.findPipExecutable()
	if err != nil {
		return nil, ErrPipNotInstalled
	}

	args := []string{"list", "--format=json"}
	output, err := m.executePipCommandWithOutput(pipPath, args)
	if err != nil {
		return nil, err
	}

	// Parse JSON output
	var packages []*Package
	if err := json.Unmarshal([]byte(output), &packages); err != nil {
		// Fallback to parsing text format
		return m.parseListOutput(output), nil
	}

	return packages, nil
}

// ShowPackage shows detailed information about a package
func (m *Manager) ShowPackage(name string) (*PackageInfo, error) {
	if name == "" {
		return nil, &PipError{
			Type:    "invalid_package_spec",
			Message: "package name cannot be empty",
		}
	}

	m.logDebug("Showing package info: %s", name)

	pipPath, err := m.findPipExecutable()
	if err != nil {
		return nil, ErrPipNotInstalled
	}

	args := []string{"show", name}
	output, err := m.executePipCommandWithOutput(pipPath, args)
	if err != nil {
		return nil, err
	}

	return m.parseShowOutput(output), nil
}

// SearchPackages searches for packages (Note: pip search was disabled, using alternative approach)
func (m *Manager) SearchPackages(query string) ([]*SearchResult, error) {
	if query == "" {
		return nil, &PipError{
			Type:    "invalid_package_spec",
			Message: "search query cannot be empty",
		}
	}

	m.logDebug("Searching packages: %s", query)

	// Since pip search is disabled, we'll return an error with a helpful message
	return nil, &PipError{
		Type:    "feature_disabled",
		Message: "pip search has been disabled. Please use https://pypi.org to search for packages",
	}
}

// FreezePackages returns a list of installed packages with versions
func (m *Manager) FreezePackages() ([]*Package, error) {
	m.logDebug("Freezing packages")

	pipPath, err := m.findPipExecutable()
	if err != nil {
		return nil, ErrPipNotInstalled
	}

	args := []string{"freeze"}
	output, err := m.executePipCommandWithOutput(pipPath, args)
	if err != nil {
		return nil, err
	}

	return m.parseFreezeOutput(output), nil
}

// executePipCommand executes a pip command and returns error if any
func (m *Manager) executePipCommand(pipPath string, args []string) error {
	_, err := m.executePipCommandWithOutput(pipPath, args)
	return err
}

// executePipCommandWithOutput executes a pip command and returns output
func (m *Manager) executePipCommandWithOutput(pipPath string, args []string) (string, error) {
	var cmd *exec.Cmd

	if strings.Contains(pipPath, " ") {
		// Handle commands like "python -m pip"
		parts := strings.Fields(pipPath)
		cmdArgs := append(parts[1:], args...)
		cmd = exec.CommandContext(m.ctx, parts[0], cmdArgs...)
	} else {
		cmd = exec.CommandContext(m.ctx, pipPath, args...)
	}

	// Set environment variables
	if len(m.config.Environment) > 0 {
		env := cmd.Env
		for key, value := range m.config.Environment {
			env = append(env, fmt.Sprintf("%s=%s", key, value))
		}
		cmd.Env = env
	}

	m.logDebug("Executing command: %s", cmd.String())

	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return outputStr, m.createPipError(cmd.String(), outputStr, exitError.ExitCode(), err)
		}
		return outputStr, m.createPipError(cmd.String(), outputStr, -1, err)
	}

	m.logDebug("Command output: %s", outputStr)
	return outputStr, nil
}

// parseListOutput parses pip list output in text format
func (m *Manager) parseListOutput(output string) []*Package {
	var packages []*Package
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Package") || strings.HasPrefix(line, "---") {
			continue
		}

		// Parse format: "package-name version"
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			pkg := &Package{
				Name:    parts[0],
				Version: parts[1],
			}

			// Check for additional info (location, editable)
			if len(parts) > 2 {
				for i := 2; i < len(parts); i++ {
					if strings.Contains(parts[i], "/") {
						pkg.Location = parts[i]
					} else if parts[i] == "editable" {
						pkg.Editable = true
					}
				}
			}

			packages = append(packages, pkg)
		}
	}

	return packages
}

// parseShowOutput parses pip show output
func (m *Manager) parseShowOutput(output string) *PackageInfo {
	info := &PackageInfo{
		Metadata: make(map[string]string),
	}

	lines := strings.Split(output, "\n")
	var currentKey string
	var currentValue strings.Builder

	for _, line := range lines {
		if strings.Contains(line, ":") && !strings.HasPrefix(line, " ") {
			// Save previous key-value pair
			if currentKey != "" {
				m.setPackageInfoField(info, currentKey, strings.TrimSpace(currentValue.String()))
			}

			// Parse new key-value pair
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				currentKey = strings.TrimSpace(parts[0])
				currentValue.Reset()
				currentValue.WriteString(strings.TrimSpace(parts[1]))
			}
		} else if currentKey != "" && strings.HasPrefix(line, " ") {
			// Continuation of previous value
			currentValue.WriteString("\n")
			currentValue.WriteString(strings.TrimSpace(line))
		}
	}

	// Save last key-value pair
	if currentKey != "" {
		m.setPackageInfoField(info, currentKey, strings.TrimSpace(currentValue.String()))
	}

	return info
}

// setPackageInfoField sets a field in PackageInfo based on key
func (m *Manager) setPackageInfoField(info *PackageInfo, key, value string) {
	switch strings.ToLower(key) {
	case "name":
		info.Name = value
	case "version":
		info.Version = value
	case "summary":
		info.Summary = value
	case "home-page":
		info.HomePage = value
	case "author":
		info.Author = value
	case "author-email":
		info.AuthorEmail = value
	case "license":
		info.License = value
	case "location":
		info.Location = value
	case "requires":
		if value != "" {
			info.Requires = strings.Split(value, ", ")
		}
	case "required-by":
		if value != "" {
			info.RequiredBy = strings.Split(value, ", ")
		}
	case "files":
		if value != "" {
			info.Files = strings.Split(value, "\n")
		}
	default:
		info.Metadata[key] = value
	}
}

// parseFreezeOutput parses pip freeze output
func (m *Manager) parseFreezeOutput(output string) []*Package {
	var packages []*Package
	lines := strings.Split(output, "\n")

	// Regex to parse package==version format
	re := regexp.MustCompile(`^([^=<>!]+)([=<>!].+)?$`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Handle editable packages
		if strings.HasPrefix(line, "-e ") {
			// Extract package name from editable line
			editableLine := strings.TrimPrefix(line, "-e ")
			editableLine = strings.TrimSpace(editableLine)

			// Try to extract package name from various formats
			var packageName string
			if strings.Contains(editableLine, "#egg=") {
				// Format: -e git+https://github.com/user/repo.git@main#egg=mypackage
				parts := strings.Split(editableLine, "#egg=")
				if len(parts) > 1 {
					packageName = parts[1]
				}
			} else if strings.Contains(editableLine, "/") {
				// Format: -e /path/to/local/package
				parts := strings.Split(editableLine, "/")
				packageName = parts[len(parts)-1]
			} else {
				packageName = editableLine
			}

			if packageName != "" {
				pkg := &Package{
					Name:     packageName,
					Editable: true,
				}
				packages = append(packages, pkg)
			}
			continue
		}

		// Skip other lines starting with "-" (but not "-e")
		if strings.HasPrefix(line, "-") {
			continue
		}

		matches := re.FindStringSubmatch(line)
		if len(matches) >= 2 {
			pkg := &Package{
				Name: matches[1],
			}

			// Parse version
			if len(matches) > 2 && matches[2] != "" {
				versionPart := matches[2]
				if strings.HasPrefix(versionPart, "==") {
					pkg.Version = strings.TrimPrefix(versionPart, "==")
				} else {
					pkg.Version = versionPart
				}
			}

			packages = append(packages, pkg)
		}
	}

	return packages
}
