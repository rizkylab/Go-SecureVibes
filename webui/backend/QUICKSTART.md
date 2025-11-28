# SecureVibes Dashboard Backend - Quick Start Guide

## 🚀 Quick Start

### 1. Install Dependencies

```bash
cd webui/backend
go mod download
```

### 2. Configuration

You can configure the application using either:

**Option A: Environment Variables**
```bash
cp .env.example .env
# Edit .env with your settings
```

**Option B: YAML Configuration**
```bash
cp config.yaml.example config.yaml
# Edit config.yaml with your settings
```

### 3. Run the Server

```bash
# Development mode
go run cmd/server/main.go

# Or build and run
go build -o securevibes-dashboard cmd/server/main.go
./securevibes-dashboard
```

The server will start on `http://localhost:8080`

### 4. Default Credentials

On first run, a default admin user is created:
- **Username:** `admin`
- **Password:** `admin123`

⚠️ **Important:** Change this password in production!

---

## 📝 API Testing

### 1. Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

Response:
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_at": "2025-11-29T06:53:51Z",
    "user": {
      "id": "...",
      "username": "admin",
      "email": "admin@securevibes.local",
      "role": "admin"
    }
  }
}
```

### 2. Get Dashboard Summary

```bash
TOKEN="your-jwt-token-here"

curl http://localhost:8080/api/v1/dashboard/summary \
  -H "Authorization: Bearer $TOKEN"
```

### 3. Upload Scan Result

```bash
curl -X POST http://localhost:8080/api/v1/scan-result \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d @scan_result.json
```

### 4. List Scans

```bash
curl http://localhost:8080/api/v1/scans?page=1&size=20 \
  -H "Authorization: Bearer $TOKEN"
```

### 5. Get Scan Details

```bash
curl http://localhost:8080/api/v1/scans/{scan_id} \
  -H "Authorization: Bearer $TOKEN"
```

### 6. Get Findings

```bash
curl "http://localhost:8080/api/v1/scans/{scan_id}/findings?severity=high" \
  -H "Authorization: Bearer $TOKEN"
```

### 7. Compare Scans

```bash
curl "http://localhost:8080/api/v1/scans/compare?scan_a={id1}&scan_b={id2}" \
  -H "Authorization: Bearer $TOKEN"
```

### 8. Export Scan

```bash
# Export as JSON
curl "http://localhost:8080/api/v1/scans/{scan_id}/export?format=json" \
  -H "Authorization: Bearer $TOKEN" \
  -o scan_report.json

# Export as Markdown
curl "http://localhost:8080/api/v1/scans/{scan_id}/export?format=markdown" \
  -H "Authorization: Bearer $TOKEN" \
  -o scan_report.md

# Export as SARIF
curl "http://localhost:8080/api/v1/scans/{scan_id}/export?format=sarif" \
  -H "Authorization: Bearer $TOKEN" \
  -o scan_report.sarif
```

---

## 🔑 API Token Management

### Create API Token

```bash
curl -X POST http://localhost:8080/api/v1/tokens \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "CI/CD Pipeline Token",
    "expires_in": 365
  }'
```

### List API Tokens

```bash
curl http://localhost:8080/api/v1/tokens \
  -H "Authorization: Bearer $TOKEN"
```

### Revoke API Token

```bash
curl -X DELETE http://localhost:8080/api/v1/tokens/{token_id} \
  -H "Authorization: Bearer $TOKEN"
```

---

## 📊 Database

The application uses SQLite database stored at `./data/scans.db` by default.

### Schema

- **scans** - Scan metadata
- **findings** - Static analysis findings
- **dast_findings** - DAST findings
- **api_tokens** - API tokens for CI/CD
- **users** - User accounts

---

## 🗂️ File Storage

Scan artifacts are stored in `./output/{scan_id}/`:
- `architecture.json` - Architecture data
- `threat_model.json` - STRIDE threat model
- `screenshots/` - DAST screenshots (if any)

---

## 🔧 Configuration Options

### Server
- `SERVER_HOST` - Server host (default: 0.0.0.0)
- `SERVER_PORT` - Server port (default: 8080)
- `CORS_ORIGINS` - Allowed CORS origins (comma-separated)

### Database
- `DATABASE_PATH` - SQLite database path (default: ./data/scans.db)

### Storage
- `OUTPUT_DIR` - Output directory for scan artifacts (default: ./output)
- `MAX_SCAN_AGE_DAYS` - Auto-delete scans older than N days (default: 90)

### Authentication
- `JWT_SECRET` - JWT signing secret (⚠️ change in production!)
- `TOKEN_EXPIRY` - JWT token expiry duration (default: 24h)

### Features
- `WEBSOCKET_ENABLED` - Enable WebSocket support (default: true)
- `AUTO_CLEANUP` - Enable automatic cleanup of old scans (default: true)

---

## 🐳 Docker

```bash
# Build image
docker build -t securevibes-dashboard .

# Run container
docker run -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/output:/app/output \
  -e JWT_SECRET=your-secret-key \
  securevibes-dashboard
```

---

## 🧪 Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...
```

---

## 📚 Project Structure

```
backend/
├── cmd/
│   └── server/          # Main application
├── config/              # Configuration management
├── database/            # Database layer
├── handlers/            # HTTP handlers
│   ├── auth.go         # Authentication
│   ├── dashboard.go    # Dashboard summary
│   ├── scans.go        # Scan management
│   ├── findings.go     # Findings management
│   ├── analysis.go     # Architecture, threats, DAST
│   ├── comparison.go   # Scan comparison
│   ├── export.go       # Report export
│   └── tokens.go       # API token management
├── middleware/          # Middleware (auth, etc.)
├── models/              # Data models
└── services/            # Business logic (future)
```

---

## 🚨 Security Notes

1. **Change default credentials** immediately in production
2. **Use strong JWT secret** (at least 32 random characters)
3. **Enable HTTPS** in production
4. **Restrict CORS origins** to your frontend domain
5. **Regularly rotate API tokens**
6. **Keep dependencies updated**

---

## 📖 Next Steps

1. ✅ Backend API is complete and functional
2. ⏳ Build frontend with SvelteKit
3. ⏳ Add WebSocket support for real-time updates
4. ⏳ Implement auto-cleanup service
5. ⏳ Add more export formats (PDF, HTML)
6. ⏳ Implement role-based access control

---

**For full documentation, see [../README.md](../README.md)**
