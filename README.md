# GoScouter

A subdomain discovery tool that uses Certificate Transparency logs to find subdomains. Features a modern React frontend and Go backend.

## Installation

### One-Line Install ⚡

```bash
curl -sSf https://raw.githubusercontent.com/nitayStain/goscouter/main/remote-install.sh | sh
```

This will automatically download, build, and install GoScouter system-wide.

### Quick Install (Alternative)

```bash
# Clone the repository
git clone https://github.com/nitayStain/goscouter.git
cd goscouter

# Run the installer
./install.sh
```

After installation:
```bash
goscouter run
```

👉 **For detailed installation options and troubleshooting, see [INSTALL.md](INSTALL.md)**

### Using Makefile

```bash
# Build everything
make build

# Install system-wide
make install

# Uninstall
make uninstall

# Clean build artifacts
make clean

# See all commands
make help
```

### Manual Installation

1. **Prerequisites:**
   - Go 1.23 or later
   - Node.js 18 or later
   - npm

2. **Build:**
   ```bash
   # Build frontend
   cd frontend
   npm install
   npm run build
   cd ..

   # Build backend
   go build -o goscouter .
   ```

3. **Run:**
   ```bash
   ./goscouter run
   ```

### Using Docker

```bash
docker compose up --build
```

The application will be available at http://localhost:8080

## Usage

### CLI Commands

```bash
goscouter run              # Start the web service
goscouter run --debug      # Start with DNS lookup logging
goscouter build            # Build frontend (quiet mode)
goscouter version          # Check current version and updates
goscouter help             # Show help message
```

### Version Checking & Auto-Update

GoScouter automatically checks for updates when you run it. If a newer version is available (which may contain security fixes), you'll be prompted to update:

```
╔═══════════════════════════════════════════════════════════════╗
║  ⚠️  Update Available: v1.0.0 → v1.2.0                        ║
║                                                               ║
║  A newer version is available and may contain important       ║
║  security fixes and improvements.                             ║
╔═══════════════════════════════════════════════════════════════╗

Would you like to update now? [Y/n]:
```

**Auto-update will:**
1. Pull latest changes from GitHub
2. Rebuild and reinstall (or rebuild from source)
3. Automatically restart the app with the new version

**Requirements for auto-update:**
- Must be in a git repository
- Git must be installed
- Make must be installed

To disable automatic checks:
```bash
export GOSCOUTER_SKIP_VERSION_CHECK=1
goscouter run
```

To manually check for updates:
```bash
goscouter version
```

### Web Interface

Open http://localhost:8080 in your browser and enter a domain to scan.

### API Endpoint

```bash
curl "http://localhost:8080/api/subdomains?domain=example.com"
```

## Development

### Frontend (Next.js + React)
```bash
cd frontend
npm install
npm run dev      # Development server on port 3000
npm run build    # Production build
```

### Backend (Go + Gin)
```bash
cd backend
go run cmd/server/main.go
```

## Features

- 🔍 Subdomain discovery via Certificate Transparency logs
- 🎨 Modern React UI with real-time results
- 🚀 Fast Go backend with Gin framework
- 📊 Statistics dashboard (subdomains found, unique IPs, certificate issuers)
- 🌐 RESTful API for programmatic access
- 🐳 Docker support for easy deployment

## Architecture

- **Frontend**: Next.js 16, React, TypeScript, Tailwind CSS
- **Backend**: Go 1.23, Gin framework
- **Data Sources**: crt.sh, certspotter.com, ipinfo.io
