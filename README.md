# 🔒 Go-SecureVibes

**A comprehensive security scanner built in Golang with multi-agent architecture for automated threat modeling, code review, and vulnerability assessment.**

---

## 🎯 Overview

Go-SecureVibes is a powerful security analysis tool that combines static code analysis, threat modeling, and optional dynamic testing to provide comprehensive security assessments of your applications. Built with a modular multi-agent architecture, it can analyze codebases in multiple languages and generate detailed security reports.

### Key Features

- 🏗️ **Architecture Assessment** - Automatically maps your application's structure, dependencies, and data flows
- 🎯 **STRIDE Threat Modeling** - Identifies potential security threats using industry-standard STRIDE methodology
- 🔍 **Static Code Analysis** - Detects vulnerabilities, insecure patterns, and code smells
- 🚀 **Dynamic Testing (Optional)** - Runtime security testing for web applications
- 📊 **Comprehensive Reporting** - Generates detailed reports in Markdown and JSON formats
- 🔧 **CI/CD Integration** - Exit codes for automated security gates
- 🌐 **Multi-Language Support** - Currently supports Go, with extensible architecture for other languages

---

## 🚀 Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/gosecvibes.git
cd gosecvibes

# Build the binary
go build -o gosecvibes cmd/gosecvibes/main.go

# Or install globally
go install ./cmd/gosecvibes
```

### Basic Usage

```bash
# Scan a project directory
./gosecvibes scan /path/to/your/project

# Scan with specific options
./gosecvibes scan /path/to/project --severity high --format json --output report.json

# Skip DAST (dynamic testing)
./gosecvibes scan /path/to/project --skip-dast

# Verbose mode for debugging
./gosecvibes scan /path/to/project --verbose
```

---

## 📋 Command Line Options

```
Flags:
  -h, --help              Show help information
  -o, --output string     Output file path (default: "SECURITY_AUDIT.md")
  -f, --format string     Output format: markdown, json, both (default: "markdown")
  -s, --severity string   Minimum severity to report: low, medium, high, critical (default: "low")
  --skip-dast            Skip dynamic application security testing
  --skip-threats         Skip threat modeling phase
  --skip-review          Skip code review phase
  -v, --verbose          Enable verbose logging
  --fail-on string       Exit with error on severity: high, critical (default: "critical")
  --exclude strings      Directories to exclude (default: vendor,node_modules,bin,dist)
```

---

## 🏗️ Architecture

Go-SecureVibes uses a **multi-agent pipeline architecture**:

```
┌─────────────────────────────────────────────────────────────┐
│                     Scanner Orchestrator                     │
└──────────────────────┬──────────────────────────────────────┘
                       │
        ┌──────────────┼──────────────┐
        │              │              │
        ▼              ▼              ▼
┌───────────────┐ ┌──────────┐ ┌─────────────┐
│  Assessment   │ │  Threat  │ │    Code     │
│     Agent     │ │ Modeling │ │   Review    │
│               │ │   Agent  │ │    Agent    │
└───────┬───────┘ └────┬─────┘ └──────┬──────┘
        │              │              │
        └──────────────┼──────────────┘
                       │
                       ▼
              ┌────────────────┐
              │  DAST Agent    │
              │   (Optional)   │
              └────────┬───────┘
                       │
                       ▼
              ┌────────────────┐
              │    Report      │
              │   Generator    │
              └────────────────┘
```

### Agent Responsibilities

1. **Assessment Agent** (`pkg/assessment/`)
   - Directory tree walking
   - Code structure parsing
   - Dependency detection
   - API endpoint identification
   - Data flow analysis
   - Configuration extraction

2. **Threat Modeling Agent** (`pkg/threatmodel/`)
   - STRIDE-based analysis
   - Attack surface mapping
   - Risk severity estimation
   - Threat scenario generation

3. **Code Review Agent** (`pkg/codereview/`)
   - Pattern-based vulnerability detection
   - CWE mapping
   - Insecure configuration detection
   - Hardcoded secrets scanning
   - Crypto usage analysis

4. **DAST Agent** (`pkg/dast/`)
   - Runtime application testing
   - HTTP endpoint fuzzing
   - Authentication testing
   - Injection attempt detection

5. **Report Generator** (`pkg/report/`)
   - Markdown report generation
   - JSON output formatting
   - Severity aggregation
   - Remediation recommendations

---

## 📁 Project Structure

```
gosecvibes/
├── cmd/
│   └── gosecvibes/
│       └── main.go              # CLI entry point
├── pkg/
│   ├── assessment/              # Architecture assessment agent
│   │   ├── agent.go
│   │   ├── parser.go
│   │   └── models.go
│   ├── threatmodel/             # Threat modeling agent
│   │   ├── agent.go
│   │   ├── stride.go
│   │   └── models.go
│   ├── codereview/              # Code review agent
│   │   ├── agent.go
│   │   ├── detectors/
│   │   │   ├── go_detector.go
│   │   │   ├── secrets_detector.go
│   │   │   └── common_detector.go
│   │   └── models.go
│   ├── dast/                    # Dynamic testing agent
│   │   ├── agent.go
│   │   └── models.go
│   ├── report/                  # Report generator
│   │   ├── generator.go
│   │   ├── markdown.go
│   │   └── json.go
│   └── scanner/                 # Main orchestrator
│       ├── scanner.go
│       └── config.go
├── internal/
│   └── utils/                   # Shared utilities
│       ├── filewalker.go
│       └── logger.go
├── examples/                    # Example vulnerable code for testing
├── docs/                        # Additional documentation
├── go.mod
├── go.sum
└── README.md
```

---

## 🔍 Supported Vulnerabilities

### Static Analysis (Code Review Agent)

- **Injection Vulnerabilities**
  - SQL Injection
  - Command Injection
  - Path Traversal
  - LDAP Injection

- **Authentication & Authorization**
  - Missing authentication
  - Weak password policies
  - Insecure session management
  - Privilege escalation

- **Cryptography**
  - Weak algorithms (MD5, SHA1)
  - Hardcoded secrets
  - Insecure random number generation
  - Missing encryption

- **Configuration**
  - Insecure defaults
  - Debug mode in production
  - Exposed sensitive endpoints
  - CORS misconfiguration

- **Data Exposure**
  - Sensitive data in logs
  - Information disclosure
  - Missing input validation
  - Insecure error handling

### Dynamic Analysis (DAST Agent)

- Open/unauthenticated endpoints
- CSRF vulnerabilities
- XSS (Cross-Site Scripting)
- Security header analysis
- SSL/TLS configuration

---

## 📊 Report Example

### Markdown Output

```markdown
# Security Audit Report

**Project:** /path/to/project  
**Scan Date:** 2025-11-26 17:00:39  
**Scanner Version:** 1.0.0

## Executive Summary

- **Total Issues:** 15
- **Critical:** 2
- **High:** 5
- **Medium:** 6
- **Low:** 2

## Findings

### [CRITICAL] SQL Injection Vulnerability
**File:** `handlers/user.go:45`  
**CWE:** CWE-89  
**Description:** Potential SQL injection in user query construction

**Code:**
```go
query := "SELECT * FROM users WHERE id = " + userId
```

**Recommendation:** Use parameterized queries or prepared statements
```

### JSON Output

```json
{
  "scan_metadata": {
    "project_path": "/path/to/project",
    "scan_date": "2025-11-26T17:00:39+07:00",
    "scanner_version": "1.0.0"
  },
  "summary": {
    "total_issues": 15,
    "critical": 2,
    "high": 5,
    "medium": 6,
    "low": 2
  },
  "findings": [...]
}
```

---

## 🔧 Configuration

Create a `.gosecvibes.yaml` in your project root for custom configuration:

```yaml
# Directories to exclude from scanning
exclude:
  - vendor
  - node_modules
  - bin
  - dist
  - .git

# Minimum severity to report
min_severity: low

# Enable/disable agents
agents:
  assessment: true
  threat_modeling: true
  code_review: true
  dast: false

# Custom rules for code review
custom_rules:
  - pattern: "password.*=.*['\"].*['\"]"
    severity: high
    message: "Hardcoded password detected"
    cwe: "CWE-798"
```

---

## 🚀 CI/CD Integration

### GitHub Actions Example

```yaml
name: Security Scan

on: [push, pull_request]

jobs:
  security:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Install Go-SecureVibes
        run: go install github.com/yourusername/gosecvibes/cmd/gosecvibes@latest
      
      - name: Run Security Scan
        run: gosecvibes scan . --fail-on high --format both
      
      - name: Upload Report
        uses: actions/upload-artifact@v3
        with:
          name: security-report
          path: SECURITY_AUDIT.*
```

---

## 🤝 Contributing

Contributions are welcome! Areas for contribution:

- **Language Support:** Add parsers for Python, JavaScript, Java, etc.
- **Vulnerability Rules:** Expand detection patterns
- **DAST Tests:** Add more dynamic testing scenarios
- **Plugins:** Create plugin system for custom detectors
- **UI:** Build web dashboard for report visualization

---

## 📝 License

MIT License - See LICENSE file for details

---

## 👤 Author

Built with ❤️ by a Security Engineer passionate about DevSecOps and automated security testing.

**Background:** Pentester | Security Engineer | Compliance Specialist (ISO 27001) | DevOps Enthusiast

---

## 🙏 Acknowledgments

- STRIDE Threat Modeling Framework
- OWASP Top 10
- CWE/SANS Top 25
- Go Security Best Practices

---

## 📚 Additional Resources

- [STRIDE Threat Modeling](https://docs.microsoft.com/en-us/azure/security/develop/threat-modeling-tool-threats)
- [OWASP Secure Coding Practices](https://owasp.org/www-project-secure-coding-practices-quick-reference-guide/)
- [CWE Database](https://cwe.mitre.org/)
- [Go Security Guidelines](https://github.com/OWASP/Go-SCP)
