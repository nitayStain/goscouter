#!/bin/bash
# GoScouter Remote Installer
# Usage: curl -sSf https://raw.githubusercontent.com/nitayStain/goscouter/main/scripts/remote-install.sh | sh

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Detect OS and Architecture
OS="$(uname -s)"
ARCH="$(uname -m)"

echo -e "${BLUE}╔═══════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║                   GoScouter Installer                     ║${NC}"
echo -e "${BLUE}║          Subdomain Discovery Tool Installation            ║${NC}"
echo -e "${BLUE}╚═══════════════════════════════════════════════════════════╝${NC}"
echo ""

# Check prerequisites
echo -e "${YELLOW}→${NC} Checking prerequisites..."

command -v git >/dev/null 2>&1 || { echo -e "${RED}✗ git is required but not installed.${NC}" >&2; exit 1; }
command -v go >/dev/null 2>&1 || { echo -e "${RED}✗ Go 1.23+ is required but not installed. Get it at https://golang.org/dl/${NC}" >&2; exit 1; }
command -v node >/dev/null 2>&1 || { echo -e "${RED}✗ Node.js 18+ is required but not installed. Get it at https://nodejs.org${NC}" >&2; exit 1; }
command -v npm >/dev/null 2>&1 || { echo -e "${RED}✗ npm is required but not installed.${NC}" >&2; exit 1; }

echo -e "${GREEN}✓${NC} All prerequisites found"
echo ""

# Set install directory
INSTALL_DIR="${HOME}/.goscouter-build"
BRANCH="${GOSCOUTER_BRANCH:-main}"

# Clean up old installation if exists
if [ -d "$INSTALL_DIR" ]; then
    echo -e "${YELLOW}→${NC} Removing old build directory..."
    rm -rf "$INSTALL_DIR"
fi

# Clone repository
echo -e "${YELLOW}→${NC} Downloading GoScouter from GitHub..."
git clone --depth 1 --branch "$BRANCH" https://github.com/nitayStain/goscouter.git "$INSTALL_DIR" >/dev/null 2>&1

cd "$INSTALL_DIR"

# Build frontend
echo -e "${YELLOW}→${NC} Building frontend (this may take a minute)..."
cd frontend
npm install --silent >/dev/null 2>&1
npm run build >/dev/null 2>&1
cd ..

# Build backend
echo -e "${YELLOW}→${NC} Building backend..."
go build -o goscouter . >/dev/null 2>&1

# Install binary
echo -e "${YELLOW}→${NC} Installing goscouter to /usr/local/bin..."
if [ -w "/usr/local/bin" ]; then
    cp goscouter /usr/local/bin/goscouter
else
    sudo cp goscouter /usr/local/bin/goscouter
fi

# Install frontend assets
echo -e "${YELLOW}→${NC} Installing frontend assets..."
mkdir -p "${HOME}/.goscouter"
cp -r frontend/out "${HOME}/.goscouter/"

# Make executable
if [ -w "/usr/local/bin/goscouter" ]; then
    chmod +x /usr/local/bin/goscouter
else
    sudo chmod +x /usr/local/bin/goscouter
fi

# Clean up build directory
cd "$HOME"
rm -rf "$INSTALL_DIR"

# Verify installation
if command -v goscouter >/dev/null 2>&1; then
    VERSION=$(goscouter version 2>&1 | head -1 | grep -o 'v[0-9.]*' || echo "dev")
    echo ""
    echo -e "${GREEN}╔═══════════════════════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║            ✓ GoScouter installed successfully!            ║${NC}"
    echo -e "${GREEN}╚═══════════════════════════════════════════════════════════╝${NC}"
    echo ""
    echo -e "${BLUE}Version:${NC} $VERSION"
    echo -e "${BLUE}Installed to:${NC} /usr/local/bin/goscouter"
    echo ""
    echo -e "${YELLOW}Get started:${NC}"
    echo -e "  ${GREEN}goscouter run${NC}       # Start the web service"
    echo -e "  ${GREEN}goscouter version${NC}   # Check version"
    echo -e "  ${GREEN}goscouter help${NC}      # Show help"
    echo ""
    echo -e "${BLUE}Web Interface:${NC} http://localhost:8080"
    echo ""
else
    echo -e "${RED}✗ Installation failed. Please check errors above.${NC}"
    exit 1
fi
