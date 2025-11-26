# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Web Dashboard Feature** (Foundation)
  - Comprehensive documentation for Web UI (400+ lines)
  - Complete data models for scans, findings, threats, and DAST
  - Main server structure with Fiber framework
  - API endpoint definitions for all features
  - CI/CD integration examples for Jenkins, GitLab, GitHub Actions
  - Implementation guide and roadmap
- Comprehensive CI/CD integration for GitHub Actions, GitLab CI, Jenkins, and CircleCI
- Docker support with multi-stage builds
- Pre-commit and pre-push Git hooks
- Makefile with common development tasks
- golangci-lint configuration
- Contributing guidelines
- Environment variable templates

### Changed
- Updated .gitignore to exclude build artifacts and reports
- Enhanced CI/CD documentation

### Planned (Web Dashboard)
- Backend API implementation (handlers, services, database)
- SvelteKit frontend with modern UI
- Architecture visualization with Cytoscape.js
- STRIDE threat model viewer
- Findings explorer with advanced filtering
- DAST results viewer
- Scan comparison functionality
- Docker deployment setup

### Fixed
- N/A

### Security
- N/A

## [1.0.0] - 2025-11-26

### Added
- Initial release of Go-SecureVibes
- Multi-agent architecture (Architecture, Threat Modeling, Static Analysis, DAST)
- STRIDE-based threat modeling
- Static code analysis for Go
- Optional dynamic application security testing
- Markdown and JSON report generation
- CI/CD mode with exit codes
- Comprehensive documentation

### Features
- 🏗️ Architecture assessment
- 🎯 STRIDE threat modeling
- 🔍 Static code analysis
- 🚀 Dynamic testing (optional)
- 📊 Comprehensive reporting
- 🔧 CI/CD integration
- 🌐 Multi-language support foundation

---

## Release Notes Format

### Version [X.Y.Z] - YYYY-MM-DD

#### Added
- New features

#### Changed
- Changes in existing functionality

#### Deprecated
- Soon-to-be removed features

#### Removed
- Removed features

#### Fixed
- Bug fixes

#### Security
- Security improvements and vulnerability fixes

---

[Unreleased]: https://github.com/rizkylab/Go-SecureVibes/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/rizkylab/Go-SecureVibes/releases/tag/v1.0.0
