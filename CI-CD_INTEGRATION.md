# CI/CD Integration Guide

This guide provides comprehensive instructions for integrating **Go-SecureVibes** into various CI/CD platforms.

---

## Table of Contents

1. [Overview](#overview)
2. [CI Mode](#ci-mode)
3. [GitHub Actions](#github-actions)
4. [GitLab CI/CD](#gitlab-cicd)
5. [Jenkins](#jenkins)
6. [CircleCI](#circleci)
7. [Docker Integration](#docker-integration)
8. [Local Development](#local-development)
9. [Best Practices](#best-practices)

---

## Overview

Go-SecureVibes provides first-class CI/CD integration with:

- **CI Mode**: Specialized mode for automated pipelines
- **Exit Codes**: Meaningful exit codes for pipeline logic
- **Multiple Formats**: JSON, Markdown, or both
- **Configurable Thresholds**: Fail builds on specific severity levels
- **Artifact Generation**: Reports suitable for archiving and review

### Exit Codes

| Code | Meaning | Typical Action |
|------|---------|----------------|
| `0` | No issues found | Continue pipeline |
| `1` | Low/Medium issues found | Continue with warning |
| `2` | High/Critical issues found | Fail build (configurable) |
| `3` | Scanner internal error | Fail build |

---

## CI Mode

Enable CI mode with the `--ci-mode` flag:

```bash
./gosecvibes scan . --ci-mode --format json --output report.json --fail-on high
```

CI mode features:
- ✅ No interactive output (colors, spinners)
- ✅ Machine-readable JSON output
- ✅ Deterministic exit codes
- ✅ Verbose logging to stderr
- ✅ Progress indicators disabled

---

## GitHub Actions

### Quick Start

Create `.github/workflows/security-scan.yml`:

```yaml
name: Security Scan

on: [push, pull_request]

jobs:
  security:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      
      - name: Build Scanner
        run: go build -o gosecvibes cmd/gosecvibes/main.go
      
      - name: Run Security Scan
        run: ./gosecvibes scan . --ci-mode --format both --fail-on high
      
      - name: Upload Reports
        if: always()
        uses: actions/upload-artifact@v4
        with:
          name: security-reports
          path: |
            security_report.json
            security_report.md
```

### Advanced Features

The provided workflow (`.github/workflows/security-scan.yml`) includes:

- ✅ **PR Comments**: Automatic security summary on pull requests
- ✅ **Artifact Upload**: Reports saved for 30 days
- ✅ **Security Events**: Integration with GitHub Security tab
- ✅ **Scheduled Scans**: Daily security audits
- ✅ **Matrix Builds**: Multi-OS and Go version testing

### Using Pre-built Workflows

We provide three workflows:

1. **`security-scan.yml`** - Comprehensive security scanning
2. **`build-test.yml`** - Build, test, and lint
3. **`release.yml`** - Multi-platform releases

Simply copy these to your `.github/workflows/` directory.

---

## GitLab CI/CD

### Quick Start

Create `.gitlab-ci.yml`:

```yaml
stages:
  - security

security:scan:
  stage: security
  image: golang:1.21
  script:
    - go build -o gosecvibes cmd/gosecvibes/main.go
    - ./gosecvibes scan . --ci-mode --format json --output report.json
  artifacts:
    reports:
      sast: report.json
    paths:
      - report.json
    expire_in: 30 days
```

### Advanced Configuration

The provided `.gitlab-ci.yml` includes:

- ✅ **Multi-stage Pipeline**: Build → Test → Security → Release → Deploy
- ✅ **Caching**: Go module and build caching
- ✅ **Coverage Reports**: Integrated test coverage
- ✅ **Secret Scanning**: GitLeaks integration
- ✅ **Dependency Scanning**: govulncheck integration
- ✅ **Docker Build**: Automated container builds
- ✅ **Manual Deployments**: Staging and production gates

### GitLab Security Dashboard

The scanner integrates with GitLab's Security Dashboard:

```yaml
artifacts:
  reports:
    sast: security_report.json
```

---

## Jenkins

### Declarative Pipeline

Create `Jenkinsfile`:

```groovy
pipeline {
    agent any
    
    environment {
        SCAN_TARGET = '.'
        SEVERITY_THRESHOLD = 'high'
    }
    
    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }
        
        stage('Security Scan') {
            steps {
                script {
                    sh 'go build -o gosecvibes cmd/gosecvibes/main.go'
                    
                    def exitCode = sh(
                        script: "./gosecvibes scan ${SCAN_TARGET} --ci-mode --format json --output security_report.json --fail-on ${SEVERITY_THRESHOLD}",
                        returnStatus: true
                    )
                    
                    if (exitCode == 3) {
                        error "Scanner internal error"
                    } else if (exitCode == 2) {
                        currentBuild.result = 'UNSTABLE'
                        echo "High/Critical vulnerabilities found!"
                    }
                }
            }
        }
        
        stage('Publish Reports') {
            steps {
                archiveArtifacts artifacts: 'security_report.json', fingerprint: true
            }
        }
    }
}
```

### Jenkins Shared Library

For reusable pipeline code, see `internal/ci/jenkins/securevibes.groovy`.

Usage:

```groovy
@Library('my-shared-lib') _

pipeline {
    agent any
    stages {
        stage('Security') {
            steps {
                secureScan(
                    target: '.',
                    threshold: 'high',
                    failBuild: true
                )
            }
        }
    }
}
```

---

## CircleCI

### Quick Start

Create `.circleci/config.yml`:

```yaml
version: 2.1

jobs:
  security-scan:
    docker:
      - image: cimg/go:1.21
    steps:
      - checkout
      - run:
          name: Build Scanner
          command: go build -o gosecvibes cmd/gosecvibes/main.go
      - run:
          name: Run Security Scan
          command: ./gosecvibes scan . --ci-mode --format json --output report.json
      - store_artifacts:
          path: report.json

workflows:
  scan:
    jobs:
      - security-scan
```

### Advanced Configuration

The provided `.circleci/config.yml` includes:

- ✅ **Orbs**: Reusable CircleCI packages
- ✅ **Executors**: Custom execution environments
- ✅ **Caching**: Go module caching
- ✅ **Workflows**: Separate build/test/release workflows
- ✅ **Artifacts**: Report storage and retrieval

---

## Docker Integration

### Using Docker

Run the scanner in a container:

```bash
# Build the image
docker build -t gosecvibes:latest .

# Scan a project
docker run --rm -v $(pwd):/workspace gosecvibes:latest scan /workspace
```

### Docker Compose

Use the provided `docker-compose.yml`:

```bash
# Scan current directory
docker-compose up

# Scan specific directory
SCAN_TARGET=/path/to/project docker-compose up

# Custom output format
FORMAT=json docker-compose up
```

### CI/CD with Docker

```yaml
# GitHub Actions example
- name: Run Security Scan
  run: |
    docker run --rm \
      -v ${{ github.workspace }}:/workspace \
      gosecvibes:latest \
      scan /workspace --ci-mode --format json
```

---

## Local Development

### Pre-commit Hooks

Install Git hooks for automatic checks:

```bash
# Setup hooks
make setup-hooks

# Or manually
git config core.hooksPath .githooks
chmod +x .githooks/*
```

Pre-commit hook checks:
- ✅ Code formatting
- ✅ Go vet
- ✅ Unit tests
- ✅ Common issues (TODOs, hardcoded secrets)
- ✅ Build verification

Pre-push hook checks:
- ✅ Full test suite with coverage
- ✅ Self security scan
- ✅ Dependency vulnerability check

### Makefile Commands

```bash
# Build the scanner
make build

# Run tests
make test

# Run tests with coverage
make test-coverage

# Format code
make fmt

# Run linters
make lint

# Run security scan on itself
make security

# Run all CI checks locally
make ci

# Build for all platforms
make build-all

# Build Docker image
make docker-build

# Setup git hooks
make setup-hooks
```

---

## Best Practices

### 1. **Fail Fast**

Configure appropriate severity thresholds:

```bash
# Fail on high or critical
--fail-on high

# Fail on critical only
--fail-on critical
```

### 2. **Archive Reports**

Always save reports as artifacts:

```yaml
# GitHub Actions
- uses: actions/upload-artifact@v4
  with:
    name: security-reports
    path: security_report.*

# GitLab CI
artifacts:
  paths:
    - security_report.*
  expire_in: 30 days
```

### 3. **Scheduled Scans**

Run regular security audits:

```yaml
# GitHub Actions - daily at 2 AM
on:
  schedule:
    - cron: '0 2 * * *'

# GitLab CI
security:scan:
  only:
    - schedules
```

### 4. **PR Integration**

Add security checks to pull requests:

```yaml
on:
  pull_request:
    branches: [main, develop]
```

### 5. **Caching**

Cache dependencies for faster builds:

```yaml
# GitHub Actions
- uses: actions/cache@v4
  with:
    path: ~/go/pkg/mod
    key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
```

### 6. **Parallel Execution**

Run security scans in parallel with other jobs:

```yaml
jobs:
  test:
    # ...
  security:
    # ...
  lint:
    # ...
```

### 7. **Security Dashboard Integration**

Integrate with platform security features:

- **GitHub**: Upload SARIF to Security tab
- **GitLab**: Use SAST report format
- **Jenkins**: Use security plugins

### 8. **Notifications**

Configure notifications for security issues:

```yaml
# Slack, email, or other notification systems
- name: Notify on failure
  if: failure()
  # Your notification step
```

---

## Troubleshooting

### Scanner Fails to Build

```bash
# Ensure Go version is correct
go version  # Should be 1.21+

# Clean and rebuild
make clean
make build
```

### Exit Code 3 (Internal Error)

```bash
# Run with verbose logging
./gosecvibes scan . --ci-mode --verbose

# Check scanner logs
```

### Reports Not Generated

```bash
# Ensure output directory exists
mkdir -p reports

# Specify full output path
./gosecvibes scan . --output reports/security_report.json
```

### Permission Denied

```bash
# Make binary executable
chmod +x gosecvibes

# For Docker
chmod +x .githooks/*
```

---

## Support

For issues or questions:

1. Check the [main README](README.md)
2. Review [example workflows](.github/workflows/)
3. Open an issue on GitHub
4. Consult platform-specific CI/CD documentation

---

## Additional Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [GitLab CI/CD Documentation](https://docs.gitlab.com/ee/ci/)
- [Jenkins Pipeline Documentation](https://www.jenkins.io/doc/book/pipeline/)
- [CircleCI Documentation](https://circleci.com/docs/)
- [Docker Documentation](https://docs.docker.com/)

---

**Last Updated**: 2025-11-26  
**Version**: 1.0.0
