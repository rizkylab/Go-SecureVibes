// vars/secureScan.groovy

def call(Map config = [:]) {
    def target = config.target ?: '.'
    def threshold = config.threshold ?: 'high'
    def failBuild = config.failBuild ?: true
    def format = config.format ?: 'json'
    
    node {
        stage('SecureVibes Scan') {
            // Ensure binary exists or download it
            if (!fileExists('gosecvibes')) {
                sh 'go build -o gosecvibes cmd/gosecvibes/main.go'
            }

            def cmd = "./gosecvibes scan ${target} --ci-mode --format ${format} --output security_report.${format}"
            
            if (config.includeDast) {
                cmd += " --include-dast"
                if (config.targetUrl) {
                    cmd += " --target ${config.targetUrl}"
                }
            }

            def exitCode = sh(script: cmd, returnStatus: true)

            archiveArtifacts artifacts: "security_report.${format}", fingerprint: true

            if (exitCode == 2 && failBuild) {
                error "Security Scan failed: High/Critical vulnerabilities found."
            } else if (exitCode == 3) {
                error "Security Scanner encountered an internal error."
            }
        }
    }
}
