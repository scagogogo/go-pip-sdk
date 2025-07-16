# Security Policy

## Supported Versions

We actively support the following versions of the Go Pip SDK:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take the security of the Go Pip SDK seriously. If you discover a security vulnerability, please follow these steps:

### 1. Do Not Create a Public Issue

Please **do not** create a public GitHub issue for security vulnerabilities. This could put users at risk.

### 2. Report Privately

Send an email to **security@scagogogo.com** with the following information:

- A clear description of the vulnerability
- Steps to reproduce the issue
- Potential impact of the vulnerability
- Any suggested fixes or mitigations
- Your contact information for follow-up

### 3. Response Timeline

We will acknowledge receipt of your vulnerability report within **48 hours** and provide a detailed response within **7 days** indicating:

- Confirmation of the vulnerability
- Our assessment of the impact
- Timeline for a fix
- Any immediate mitigation steps

### 4. Coordinated Disclosure

We follow a coordinated disclosure process:

1. We will work with you to understand and validate the vulnerability
2. We will develop and test a fix
3. We will prepare a security advisory
4. We will release the fix and publish the advisory
5. We will credit you for the discovery (if desired)

## Security Best Practices

When using the Go Pip SDK, please follow these security best practices:

### 1. Keep Dependencies Updated

- Regularly update the SDK to the latest version
- Monitor security advisories for dependencies
- Use tools like `go mod tidy` and `go list -m -u all` to check for updates

### 2. Validate Input

- Always validate package names and versions before installation
- Sanitize user input when constructing package specifications
- Use version constraints to prevent unexpected package versions

### 3. Use Trusted Package Indexes

- Configure trusted package indexes in your pip configuration
- Avoid installing packages from untrusted sources
- Verify package integrity when possible

### 4. Virtual Environment Isolation

- Use virtual environments to isolate package installations
- Avoid installing packages globally when possible
- Clean up virtual environments when no longer needed

### 5. Logging and Monitoring

- Enable appropriate logging levels for your environment
- Monitor pip operations for suspicious activity
- Review logs regularly for security events

### 6. Network Security

- Use HTTPS for package downloads
- Configure appropriate proxy settings if required
- Validate SSL certificates for package indexes

## Common Security Considerations

### Package Installation

```go
// Good: Validate package names
func installPackage(name string) error {
    if !isValidPackageName(name) {
        return errors.New("invalid package name")
    }
    
    pkg := &pip.PackageSpec{
        Name: name,
        // Use version constraints
        Version: ">=1.0.0,<2.0.0",
    }
    
    return manager.InstallPackage(pkg)
}

// Bad: Direct user input without validation
func installPackageUnsafe(userInput string) error {
    pkg := &pip.PackageSpec{Name: userInput}
    return manager.InstallPackage(pkg)
}
```

### Virtual Environment Management

```go
// Good: Validate paths
func createVenv(path string) error {
    cleanPath := filepath.Clean(path)
    if !isValidPath(cleanPath) {
        return errors.New("invalid path")
    }
    
    return manager.CreateVenv(cleanPath)
}

// Bad: Direct path usage without validation
func createVenvUnsafe(userPath string) error {
    return manager.CreateVenv(userPath)
}
```

### Configuration Security

```go
// Good: Secure configuration
config := &pip.Config{
    DefaultIndex: "https://pypi.org/simple/",
    TrustedHosts: []string{"pypi.org", "pypi.python.org"},
    Timeout:      30 * time.Second,
    Environment: map[string]string{
        "PIP_DISABLE_PIP_VERSION_CHECK": "1",
    },
}

// Bad: Insecure configuration
config := &pip.Config{
    DefaultIndex: "http://untrusted-index.com/simple/",
    // No trusted hosts specified
    // No timeout specified
}
```

## Vulnerability Disclosure History

We will maintain a record of disclosed vulnerabilities here:

- No vulnerabilities have been disclosed to date.

## Security Tools

We use the following tools to maintain security:

- **Dependabot**: Automated dependency updates
- **CodeQL**: Static analysis for security vulnerabilities
- **gosec**: Go security checker
- **nancy**: Vulnerability scanner for Go dependencies

## Contact

For security-related questions or concerns:

- Email: security@scagogogo.com
- For general questions: [GitHub Discussions](https://github.com/scagogogo/go-pip-sdk/discussions)

## Acknowledgments

We appreciate the security research community and will acknowledge researchers who responsibly disclose vulnerabilities to us.

---

Thank you for helping keep the Go Pip SDK and its users safe!
