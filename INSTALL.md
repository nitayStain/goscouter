# GoScouter Installation Guide

## Prerequisites

Before installing GoScouter, ensure you have the following installed:

- **Go** 1.23 or later - [Download](https://golang.org/dl/)
- **Node.js** 18 or later - [Download](https://nodejs.org/)
- **npm** (comes with Node.js)

Verify your installations:
```bash
go version
node --version
npm --version
```

## Installation Methods

### Method 1: Quick Install (Recommended)

The easiest way to install GoScouter system-wide:

```bash
# Clone the repository
git clone https://github.com/nitayStain/goscouter.git
cd goscouter

# Run the installer
./install.sh
```

This will:
1. Check for prerequisites (Go, Node.js, npm)
2. Build the frontend (React/Next.js)
3. Build the backend (Go)
4. Install `goscouter` binary to `/usr/local/bin`
5. Copy frontend files to `~/.goscouter`

After installation, you can run from anywhere:
```bash
goscouter run
```

### Method 2: Using Makefile

For developers who want more control:

```bash
# Build frontend and backend
make build

# Install system-wide
make install

# Or combine both
make build install
```

Other useful Makefile commands:
```bash
make help       # Show all available commands
make clean      # Remove build artifacts
make uninstall  # Remove installed files
make run        # Build and run locally
make dev        # Run frontend in dev mode
make test       # Run tests
```

### Method 3: Manual Installation

For advanced users or custom setups:

1. **Build the frontend:**
   ```bash
   cd frontend
   npm install
   npm run build
   cd ..
   ```

2. **Build the backend:**
   ```bash
   go build -o goscouter .
   ```

3. **Install (optional):**
   ```bash
   sudo cp goscouter /usr/local/bin/
   mkdir -p ~/.goscouter
   cp -r frontend/out ~/.goscouter/
   ```

4. **Run:**
   ```bash
   goscouter run
   ```

### Method 4: Using Docker

No local installation required:

```bash
# Build and run
docker compose up --build

# Run in background
docker compose up -d

# Stop
docker compose down
```

## Uninstallation

### Using the uninstall script:
```bash
./uninstall.sh
```

### Using Make:
```bash
make uninstall
```

### Manual uninstall:
```bash
sudo rm /usr/local/bin/goscouter
rm -rf ~/.goscouter
```

## Troubleshooting

### "Command not found" after installation

Make sure `/usr/local/bin` is in your PATH:
```bash
echo $PATH | grep /usr/local/bin
```

If not, add to your shell profile (`~/.bashrc`, `~/.zshrc`, etc.):
```bash
export PATH="/usr/local/bin:$PATH"
```

### Permission denied during installation

The installer will automatically request sudo access when needed. If you prefer not to use sudo, you can:
1. Install to a user directory:
   ```bash
   mkdir -p ~/bin
   cp goscouter ~/bin/
   export PATH="$HOME/bin:$PATH"
   ```

### Frontend not building

Ensure you have npm installed and try:
```bash
cd frontend
rm -rf node_modules .next out
npm install
npm run build
```

### Port 8080 already in use

Find and kill the process using port 8080:
```bash
lsof -ti:8080 | xargs kill -9
```

Or configure a different port in the server code.

## Verifying Installation

After installation, verify goscouter is working:

```bash
# Check version/help
goscouter help

# Start the server
goscouter run

# In another terminal, test the API
curl "http://localhost:8080/api/subdomains?domain=example.com"
```

## Next Steps

- Read the [README](README.md) for usage instructions
- Try scanning a domain: `http://localhost:8080`
- Check out the API documentation
- Report issues on GitHub

## Installation Locations

- **Binary**: `/usr/local/bin/goscouter`
- **Frontend files**: `~/.goscouter/out/`
- **Config** (future): `~/.goscouter/config.yaml`
