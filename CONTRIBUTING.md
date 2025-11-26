# Contributing to Go-SecureVibes

Thank you for your interest in contributing to Go-SecureVibes! This document provides guidelines and instructions for contributing.

---

## Table of Contents

1. [Code of Conduct](#code-of-conduct)
2. [Getting Started](#getting-started)
3. [Development Setup](#development-setup)
4. [Making Changes](#making-changes)
5. [Testing](#testing)
6. [Submitting Changes](#submitting-changes)
7. [Coding Standards](#coding-standards)
8. [Areas for Contribution](#areas-for-contribution)

---

## Code of Conduct

This project adheres to a code of conduct. By participating, you are expected to uphold this code. Please be respectful and constructive in all interactions.

---

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Git
- Make (optional but recommended)
- Docker (for container testing)

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork:

```bash
git clone https://github.com/YOUR_USERNAME/Go-SecureVibes.git
cd Go-SecureVibes
```

3. Add upstream remote:

```bash
git remote add upstream https://github.com/rizkylab/Go-SecureVibes.git
```

---

## Development Setup

### Install Dependencies

```bash
# Download Go modules
go mod download

# Verify dependencies
go mod verify
```

### Setup Git Hooks

```bash
# Install pre-commit and pre-push hooks
make setup-hooks

# Or manually
git config core.hooksPath .githooks
chmod +x .githooks/*
```

### Build the Project

```bash
# Using Make
make build

# Or directly with Go
go build -o gosecvibes cmd/gosecvibes/main.go
```

### Run Tests

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# View coverage in browser
make coverage-html
```

---

## Making Changes

### Create a Branch

```bash
# Update your fork
git fetch upstream
git checkout main
git merge upstream/main

# Create a feature branch
git checkout -b feature/your-feature-name
```

### Branch Naming Convention

- `feature/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation updates
- `refactor/` - Code refactoring
- `test/` - Test additions or updates
- `ci/` - CI/CD changes

Examples:
- `feature/add-python-support`
- `fix/scanner-crash-on-empty-file`
- `docs/update-installation-guide`

### Commit Messages

Follow conventional commits format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation
- `style`: Formatting, missing semicolons, etc.
- `refactor`: Code restructuring
- `test`: Adding tests
- `chore`: Maintenance tasks

Examples:

```
feat(scanner): add Python language support

- Implemented Python AST parser
- Added Python-specific vulnerability patterns
- Updated documentation

Closes #123
```

```
fix(threatmodel): correct STRIDE category mapping

The STRIDE mapping was incorrectly categorizing some threats.
This fix ensures proper categorization based on OWASP guidelines.

Fixes #456
```

---

## Testing

### Unit Tests

```bash
# Run all tests
go test ./...

# Run specific package
go test ./internal/scanner

# Run with race detection
go test -race ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
```

### Integration Tests

```bash
# Run integration tests
make run-example

# Or manually
./gosecvibes scan examples --format json --output test_report.json
```

### Test Coverage Requirements

- Aim for >70% coverage for new code
- All public functions should have tests
- Critical paths must have tests

### Writing Tests

```go
func TestScannerBasic(t *testing.T) {
    // Arrange
    scanner := NewScanner(Config{
        Target: "testdata/sample",
    })
    
    // Act
    result, err := scanner.Scan()
    
    // Assert
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    
    if len(result.Findings) == 0 {
        t.Error("expected findings, got none")
    }
}
```

---

## Submitting Changes

### Before Submitting

1. **Run all checks**:
   ```bash
   make ci
   ```

2. **Update documentation** if needed

3. **Add tests** for new functionality

4. **Run security scan**:
   ```bash
   make security
   ```

### Create Pull Request

1. Push your branch:
   ```bash
   git push origin feature/your-feature-name
   ```

2. Go to GitHub and create a Pull Request

3. Fill out the PR template:
   - Description of changes
   - Related issues
   - Testing performed
   - Screenshots (if UI changes)

### PR Review Process

1. **Automated Checks**: CI/CD pipelines will run
2. **Code Review**: Maintainers will review your code
3. **Feedback**: Address any requested changes
4. **Approval**: Once approved, your PR will be merged

---

## Coding Standards

### Go Style Guide

Follow the [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md) and [Effective Go](https://golang.org/doc/effective_go).

### Code Formatting

```bash
# Format code
make fmt

# Or
gofmt -s -w .
```

### Linting

```bash
# Run linters
make lint

# Or
golangci-lint run
```

### Documentation

- Add godoc comments for all exported functions
- Update README.md for user-facing changes
- Add examples for new features

Example:

```go
// Scan performs a comprehensive security scan on the target directory.
// It returns a ScanResult containing all findings and metadata.
//
// Example:
//   scanner := NewScanner(Config{Target: "/path/to/project"})
//   result, err := scanner.Scan()
//   if err != nil {
//       log.Fatal(err)
//   }
func (s *Scanner) Scan() (*ScanResult, error) {
    // Implementation
}
```

### Error Handling

```go
// Good: Wrap errors with context
if err != nil {
    return fmt.Errorf("failed to parse file %s: %w", filename, err)
}

// Bad: Lose error context
if err != nil {
    return err
}
```

### Logging

Use the internal logger, not `fmt.Println`:

```go
// Good
logger.Info("Starting scan", "target", target)
logger.Error("Scan failed", "error", err)

// Bad
fmt.Println("Starting scan:", target)
```

---

## Areas for Contribution

### High Priority

1. **Language Support**
   - Python parser and detectors
   - JavaScript/TypeScript support
   - Java support
   - PHP support

2. **Vulnerability Detection**
   - New vulnerability patterns
   - CWE coverage expansion
   - OWASP Top 10 coverage

3. **DAST Improvements**
   - More test scenarios
   - Authentication handling
   - API testing support

### Medium Priority

4. **Reporting**
   - HTML report generation
   - PDF export
   - SARIF format support
   - Custom report templates

5. **Performance**
   - Parallel scanning optimization
   - Memory usage reduction
   - Caching improvements

6. **CI/CD Integration**
   - Azure DevOps support
   - Bitbucket Pipelines
   - Travis CI examples

### Nice to Have

7. **UI/UX**
   - Web dashboard
   - VSCode extension
   - CLI improvements

8. **Documentation**
   - Video tutorials
   - Blog posts
   - Use case examples

9. **Testing**
   - Benchmark suite
   - Fuzzing tests
   - More integration tests

---

## Getting Help

- **Questions**: Open a GitHub Discussion
- **Bugs**: Open a GitHub Issue
- **Security**: Email security@example.com (do not open public issues)
- **Chat**: Join our Discord/Slack (if available)

---

## Recognition

Contributors will be:
- Listed in CONTRIBUTORS.md
- Mentioned in release notes
- Credited in documentation

---

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to Go-SecureVibes! 🔒
