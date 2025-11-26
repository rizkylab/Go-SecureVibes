# GitHub Configuration

This directory contains GitHub-specific configuration files for the Go-SecureVibes project.

## 📁 Structure

```
.github/
├── workflows/           # GitHub Actions workflows
│   ├── security-scan.yml    # Security scanning workflow
│   ├── build-test.yml       # Build and test workflow
│   └── release.yml          # Release automation workflow
├── ISSUE_TEMPLATE/      # Issue templates
│   ├── bug_report.md        # Bug report template
│   └── feature_request.md   # Feature request template
└── pull_request_template.md # PR template
```

## 🔄 Workflows

### 1. Security Scan (`security-scan.yml`)

**Triggers:**
- Push to `main` or `develop`
- Pull requests to `main` or `develop`
- Daily schedule (2 AM UTC)
- Manual dispatch

**Features:**
- Automated security scanning with Go-SecureVibes
- PR comments with security summary
- Artifact uploads (30-day retention)
- GitHub Security tab integration
- Exit code-based failure handling

### 2. Build & Test (`build-test.yml`)

**Triggers:**
- Push to `main` or `develop`
- Pull requests to `main` or `develop`

**Features:**
- Multi-OS testing (Ubuntu, macOS, Windows)
- Multi-Go version (1.21, 1.22)
- Code coverage with Codecov
- golangci-lint integration
- Integration tests
- Binary artifact uploads

### 3. Release (`release.yml`)

**Triggers:**
- Git tags matching `v*.*.*`
- Manual workflow dispatch

**Features:**
- Multi-platform binary builds (Linux, macOS, Windows)
- ARM64 and AMD64 architectures
- Automatic changelog generation
- GitHub Release creation
- Docker image publishing to GHCR
- SHA256 checksums

## 📝 Templates

### Issue Templates

**Bug Report** - For reporting bugs and issues
- Environment details
- Reproduction steps
- Expected vs actual behavior
- Error logs and screenshots

**Feature Request** - For suggesting new features
- Problem statement
- Proposed solution
- Use cases
- Implementation ideas

### Pull Request Template

Comprehensive checklist including:
- Type of change
- Related issues
- Testing performed
- Security review
- Documentation updates
- Code quality checks

## 🚀 Usage

### Running Workflows Locally

Use [act](https://github.com/nektos/act) to test workflows locally:

```bash
# Install act
brew install act

# Run security scan workflow
act -j security-scan

# Run build workflow
act -j build
```

### Creating Issues

1. Go to the [Issues](../../issues) tab
2. Click "New Issue"
3. Select appropriate template
4. Fill in all required fields

### Creating Pull Requests

1. Create a feature branch
2. Make your changes
3. Push to your fork
4. Open a PR - template will auto-populate
5. Fill in all sections

## 🔐 Secrets Configuration

Required secrets for workflows:

| Secret | Purpose | Required For |
|--------|---------|--------------|
| `GITHUB_TOKEN` | Automatic | All workflows (auto-provided) |
| `CODECOV_TOKEN` | Code coverage | build-test.yml (optional) |

To add secrets:
1. Go to Settings → Secrets and variables → Actions
2. Click "New repository secret"
3. Add name and value

## 📊 Workflow Status

Check workflow status:
- [Actions Tab](../../actions)
- Status badges in README.md
- Email notifications (configure in Settings)

## 🛠️ Customization

### Modifying Workflows

1. Edit YAML files in `workflows/`
2. Test locally with `act`
3. Commit and push
4. Monitor in Actions tab

### Adding New Workflows

1. Create new `.yml` file in `workflows/`
2. Define triggers and jobs
3. Test thoroughly
4. Document in this README

### Modifying Templates

1. Edit files in `ISSUE_TEMPLATE/` or `pull_request_template.md`
2. Use GitHub's template syntax
3. Test by creating a new issue/PR

## 📚 Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Workflow Syntax](https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions)
- [Issue Templates](https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests)
- [act - Local Testing](https://github.com/nektos/act)

## 🤝 Contributing

When modifying GitHub configuration:

1. Test changes thoroughly
2. Update this README if needed
3. Follow existing patterns
4. Document new features

---

For more information, see the main [CI/CD Integration Guide](../CI-CD_INTEGRATION.md).
