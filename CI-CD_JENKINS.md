# Jenkins CI/CD Integration for Go-SecureVibes

This document details how to integrate **Go-SecureVibes** into your Jenkins pipelines.

## Overview

Go-SecureVibes supports a dedicated `--ci-mode` that:
1.  Disables interactive output (colors, spinners).
2.  Forces JSON output (by default).
3.  Returns specific exit codes for pipeline logic:
    *   `0`: No issues found.
    *   `1`: Low/Medium issues found.
    *   `2`: High/Critical issues found (typically fails the build).
    *   `3`: Scanner internal error.

## Integration Methods

### 1. Direct Pipeline Integration (Jenkinsfile)

You can run the scanner directly using `sh` steps in your `Jenkinsfile`.

```groovy
stage('Security Scan') {
    steps {
        sh 'go build -o gosecvibes cmd/gosecvibes/main.go'
        
        // Run scan
        // Returns 0, 1, or 2 based on findings
        sh './gosecvibes scan . --ci-mode --output report.json'
    }
}
```

See `internal/ci/jenkins/Jenkinsfile` for a complete example with error handling and artifact archiving.

### 2. Jenkins Shared Library

If you have a Shared Library set up, you can use the provided Groovy script wrapper.

1.  Copy `internal/ci/jenkins/securevibes.groovy` to your Shared Library's `vars/` directory (rename to `secureScan.groovy`).
2.  Call it in your pipeline:

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
                    failBuild: true,
                    includeDast: false
                )
            }
        }
    }
}
```

## Configuration Options

| Flag | Description |
|------|-------------|
| `--ci-mode` | Enables CI mode (required). |
| `--format` | Output format (`json`, `markdown`, `both`). Default `json` in CI mode. |
| `--output` | Path to save the report. |
| `--include-dast` | Enable Dynamic Analysis (requires running app). |
| `--target` | Target URL for DAST (e.g., `http://localhost:8080`). |

## Artifacts

The scanner generates a report file (default `SECURITY_AUDIT.md` or `report.json`). You should use `archiveArtifacts` to save this file for audit purposes.

```groovy
archiveArtifacts artifacts: 'report.json', fingerprint: true
```
