# SecureVibes Web Dashboard

**A full-featured Web UI for visualizing Go-SecureVibes security scanning results**

---

## 🎯 Overview

SecureVibes Web Dashboard is a modern, responsive web application built to visualize and manage security scan results from Go-SecureVibes. It provides comprehensive views of architecture maps, threat models, static analysis findings, DAST results, and scan comparisons.

### Key Features

- 📊 **Dashboard Overview** - Real-time security metrics and trends
- 🏗️ **Architecture Visualization** - Interactive component graphs
- 🎯 **STRIDE Threat Model** - Comprehensive threat analysis
- 🔍 **Findings Explorer** - Advanced filtering and search
- 🚀 **DAST Results** - Dynamic testing evidence
- 📈 **Scan Comparison** - Track security improvements over time
- 🔗 **CI/CD Integration** - Automatic scan uploads from Jenkins/GitLab/GitHub
- 📥 **Report Export** - Download in multiple formats

---

## 🏗️ Architecture

```
┌────────────────────────────────────────────────────────────┐
│                     Web Dashboard                          │
├──────────────────┬──────────────────┬──────────────────────┤
│   Frontend       │    Backend       │     Storage          │
│   (SvelteKit)    │    (Fiber)       │     (SQLite)         │
├──────────────────┼──────────────────┼──────────────────────┤
│ • Dashboard      │ • REST API       │ • Scan Metadata      │
│ • Architecture   │ • WebSocket      │ • Findings Index     │
│ • Threat Model   │ • Auth (JWT)     │ • User Sessions      │
│ • Findings       │ • File Storage   │ • API Tokens         │
│ • DAST Results   │ • Comparison     │                      │
│ • Scan History   │ • Export         │                      │
└──────────────────┴──────────────────┴──────────────────────┘
                            │
                            ▼
                  ┌─────────────────────┐
                  │   CI/CD Systems     │
                  ├─────────────────────┤
                  │ • Jenkins           │
                  │ • GitLab CI         │
                  │ • GitHub Actions    │
                  │ • CircleCI          │
                  └─────────────────────┘
```

---

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Node.js 18+
- SQLite 3

### Backend Setup

```bash
cd webui/backend

# Install dependencies
go mod download

# Run database migrations
go run cmd/migrate/main.go

# Start server
go run cmd/server/main.go
```

Server will start on `http://localhost:8080`

### Frontend Setup

```bash
cd webui/frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

Frontend will start on `http://localhost:5173`

### Docker Setup

```bash
# Build and run with docker-compose
docker-compose -f webui/docker-compose.yml up -d
```

Access at `http://localhost:3000`

---

## 📋 Features

### 1. Dashboard Overview

**URL:** `/`

**Features:**
- Total scans count
- Recent scan summary
- Severity breakdown (pie chart)
- Trend over time (line chart)
- Recent findings list
- Quick actions

**API Endpoint:**
```bash
GET /api/v1/dashboard/summary
```

### 2. Scan List

**URL:** `/scans`

**Features:**
- Paginated scan history
- Filter by date range, status, severity
- Sort by timestamp, duration, findings count
- Quick view summary
- Delete/archive scans

**API Endpoints:**
```bash
GET /api/v1/scans?page=1&size=20&status=completed
GET /api/v1/scans/{id}
DELETE /api/v1/scans/{id}
```

### 3. Architecture Visualization

**URL:** `/scans/{id}/architecture`

**Features:**
- Interactive component graph (Cytoscape.js)
- Component details panel
- Data flow visualization
- Zoom, pan, layout options
- Export as PNG/SVG

**API Endpoint:**
```bash
GET /api/v1/scans/{id}/architecture
```

**Graph Types:**
- HTTP Handlers
- Services
- Databases
- External APIs
- Data Flows

### 4. Threat Model (STRIDE)

**URL:** `/scans/{id}/threats`

**Features:**
- Filterable threat table
- Severity-based coloring
- Component grouping
- Mitigation status tracking
- Export to CSV/JSON

**API Endpoint:**
```bash
GET /api/v1/scans/{id}/threat-model
```

**Threat Categories:**
- **S**poofing
- **T**ampering
- **R**epudiation
- **I**nformation Disclosure
- **D**enial of Service
- **E**levation of Privilege

### 5. Static Analysis Findings

**URL:** `/scans/{id}/findings`

**Features:**
- Advanced filtering (CWE, severity, category, file)
- Code preview with syntax highlighting
- Line-level highlighting
- Remediation suggestions
- Mark as false positive
- Assign to team members (future)

**API Endpoints:**
```bash
GET /api/v1/scans/{id}/findings?severity=high&cwe=CWE-89
GET /api/v1/scans/{id}/findings/{finding_id}
PATCH /api/v1/scans/{id}/findings/{finding_id}/status
```

**Categories:**
- Cryptography
- Authentication
- Injection
- Configuration
- Secrets
- Authorization
- Data Exposure

### 6. DAST Results

**URL:** `/scans/{id}/dast`

**Features:**
- Endpoint vulnerability list
- Request/Response viewer
- Evidence display
- Screenshot gallery
- Exploitation proof
- Remediation steps

**API Endpoint:**
```bash
GET /api/v1/scans/{id}/dast-findings
```

### 7. Scan Comparison

**URL:** `/compare?a={scan_id_1}&b={scan_id_2}`

**Features:**
- Side-by-side comparison
- New vulnerabilities highlight
- Fixed vulnerabilities list
- Severity trend chart
- Delta calculations
- Export comparison report

**API Endpoint:**
```bash
GET /api/v1/scans/compare?scan_a={id1}&scan_b={id2}
```

**Metrics:**
- New findings
- Fixed findings
- Severity changes
- Category distribution
- Time to fix

### 8. Report Export

**Formats:**
- Markdown (`.md`)
- JSON (`.json`)
- SARIF (`.sarif`)
- PDF (future)
- HTML (future)

**API Endpoints:**
```bash
GET /api/v1/scans/{id}/export?format=markdown
GET /api/v1/scans/{id}/export?format=json
GET /api/v1/scans/{id}/export?format=sarif
```

---

## 🔗 CI/CD Integration

### Jenkins Integration

Add to your `Jenkinsfile`:

```groovy
stage('Upload to Dashboard') {
    steps {
        script {
            sh '''
            curl -X POST https://dashboard.example.com/api/v1/scan-result \
              -H "Authorization: Bearer ${SECUREVIBES_TOKEN}" \
              -H "Content-Type: application/json" \
              --data @output/security_report.json
            '''
        }
    }
}
```

### GitLab CI Integration

Add to `.gitlab-ci.yml`:

```yaml
upload_to_dashboard:
  stage: deploy
  script:
    - |
      curl -X POST $DASHBOARD_URL/api/v1/scan-result \
        -H "Authorization: Bearer $SECUREVIBES_TOKEN" \
        -H "Content-Type: application/json" \
        --data @security_report.json
```

### GitHub Actions Integration

Add to workflow:

```yaml
- name: Upload to Dashboard
  run: |
    curl -X POST ${{ secrets.DASHBOARD_URL }}/api/v1/scan-result \
      -H "Authorization: Bearer ${{ secrets.SECUREVIBES_TOKEN }}" \
      -H "Content-Type: application/json" \
      --data @security_report.json
```

---

## 🔐 API Reference

### Authentication

All API requests require authentication using JWT tokens or API keys.

**Get Token:**
```bash
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "your_password"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_at": "2025-11-27T20:00:00Z"
  }
}
```

**Use Token:**
```bash
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

### API Endpoints

#### Scans

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/scans` | List all scans |
| GET | `/api/v1/scans/{id}` | Get scan details |
| POST | `/api/v1/scan-result` | Upload new scan |
| DELETE | `/api/v1/scans/{id}` | Delete scan |

#### Findings

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/scans/{id}/findings` | List findings |
| GET | `/api/v1/scans/{id}/findings/{finding_id}` | Get finding details |
| PATCH | `/api/v1/scans/{id}/findings/{finding_id}` | Update finding status |

#### Architecture

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/scans/{id}/architecture` | Get architecture data |

#### Threat Model

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/scans/{id}/threat-model` | Get threat model |

#### DAST

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/scans/{id}/dast-findings` | Get DAST findings |

#### Comparison

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/scans/compare` | Compare two scans |

#### Export

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/scans/{id}/export` | Export scan report |

---

## 🎨 UI Design

### Color Scheme

```css
/* Primary Colors */
--primary: #1f3b73;        /* Security Blue */
--secondary: #2c5aa0;

/* Severity Colors */
--critical: #8b0000;       /* Dark Red */
--high: #d62828;           /* Red */
--medium: #f77f00;         /* Orange */
--low: #3a86ff;            /* Blue */
--info: #38b000;           /* Green */

/* UI Colors */
--background: #0f172a;     /* Dark Background */
--surface: #1e293b;        /* Card Background */
--text: #f1f5f9;           /* Text */
--text-muted: #94a3b8;     /* Muted Text */
```

### Layout

```
┌─────────────────────────────────────────────────────────┐
│  Navbar: Logo | Search | User Menu                      │
├──────────┬──────────────────────────────────────────────┤
│          │                                              │
│ Sidebar  │          Main Content Area                   │
│          │                                              │
│ • Dashboard                                             │
│ • Scans  │  ┌────────────────────────────────────┐      │
│ • Arch   │  │                                    │      │
│ • Threats│  │         Content Cards              │      │
│ • Findings                                       │      │
│ • DAST   │  └────────────────────────────────────┘      │
│ • Compare│                                              │
│ • Settings                                              │
│          │                                              │
└──────────┴──────────────────────────────────────────────┘
```

---

## 📊 Data Storage

### SQLite Schema

```sql
-- Scans table
CREATE TABLE scans (
    id TEXT PRIMARY KEY,
    timestamp DATETIME NOT NULL,
    project_path TEXT NOT NULL,
    commit_hash TEXT,
    branch TEXT,
    duration INTEGER,
    dast_enabled BOOLEAN DEFAULT 0,
    summary_json TEXT,
    status TEXT DEFAULT 'completed',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Findings table
CREATE TABLE findings (
    id TEXT PRIMARY KEY,
    scan_id TEXT NOT NULL,
    title TEXT NOT NULL,
    severity TEXT NOT NULL,
    cwe TEXT,
    category TEXT,
    file_path TEXT,
    line_number INTEGER,
    status TEXT DEFAULT 'new',
    first_detected DATETIME,
    last_seen DATETIME,
    FOREIGN KEY (scan_id) REFERENCES scans(id) ON DELETE CASCADE
);

-- API Tokens table
CREATE TABLE api_tokens (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    token TEXT UNIQUE NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME,
    last_used DATETIME
);
```

### File Storage

```
output/
  scan_2025-11-26_12-33-00/
    metadata.json          # Scan metadata
    architecture.json      # Architecture data
    threat_model.json      # STRIDE threats
    findings.json          # Static analysis findings
    dast_findings.json     # DAST results
    SECURITY_AUDIT.md      # Markdown report
    screenshots/           # DAST screenshots
      endpoint_1.png
      endpoint_2.png
```

---

## 🔧 Configuration

### Backend Config (`config.yaml`)

```yaml
server:
  host: 0.0.0.0
  port: 8080
  cors_origins:
    - http://localhost:5173
    - https://dashboard.example.com

database:
  path: ./data/scans.db
  
storage:
  output_dir: ./output
  max_scan_age_days: 90  # Auto-delete old scans

auth:
  jwt_secret: your-secret-key
  token_expiry: 24h
  
features:
  websocket_enabled: true
  auto_cleanup: true
```

### Frontend Config (`.env`)

```bash
VITE_API_URL=http://localhost:8080
VITE_WS_URL=ws://localhost:8080/ws
VITE_APP_TITLE=SecureVibes Dashboard
```

---

## 🚢 Deployment

### Docker Deployment

```bash
# Build images
docker-compose -f webui/docker-compose.yml build

# Start services
docker-compose -f webui/docker-compose.yml up -d

# View logs
docker-compose -f webui/docker-compose.yml logs -f
```

### Production Deployment

```bash
# Backend
cd webui/backend
go build -o securevibes-dashboard cmd/server/main.go
./securevibes-dashboard

# Frontend
cd webui/frontend
npm run build
# Serve dist/ with nginx or similar
```

---

## 🧪 Development

### Running Tests

```bash
# Backend tests
cd webui/backend
go test ./...

# Frontend tests
cd webui/frontend
npm test
```

### API Testing

```bash
# Get all scans
curl http://localhost:8080/api/v1/scans

# Upload scan result
curl -X POST http://localhost:8080/api/v1/scan-result \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  --data @test_scan.json
```

---

## 📚 Tech Stack

### Backend
- **Framework**: Fiber (Go)
- **Database**: SQLite
- **Auth**: JWT
- **WebSocket**: Fiber WebSocket
- **Validation**: go-playground/validator

### Frontend
- **Framework**: SvelteKit
- **UI**: Tailwind CSS
- **Charts**: ApexCharts
- **Graphs**: Cytoscape.js
- **Code Viewer**: Monaco Editor
- **Icons**: Heroicons

---

## 🔮 Future Enhancements

- [ ] Real-time scan updates via WebSocket
- [ ] Team collaboration features
- [ ] JIRA/GitHub Issues integration
- [ ] Git blame integration
- [ ] Role-based access control (RBAC)
- [ ] OAuth login (Google/GitHub)
- [ ] AI-powered remediation suggestions
- [ ] PDF report generation
- [ ] Email notifications
- [ ] Slack/Discord webhooks
- [ ] Custom dashboard widgets
- [ ] Multi-tenant support

---

## 📖 Documentation

- [API Documentation](./API.md)
- [Deployment Guide](./DEPLOYMENT.md)
- [Development Guide](./DEVELOPMENT.md)
- [User Guide](./USER_GUIDE.md)

---

## 🤝 Contributing

See [CONTRIBUTING.md](../../CONTRIBUTING.md) for development guidelines.

---

## 📄 License

MIT License - See [LICENSE](../../LICENSE) for details.

---

**Built with ❤️ for secure software development**
