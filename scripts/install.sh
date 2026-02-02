#!/bin/sh
# GoScouter Installer
# Usage: curl -sSf https://raw.githubusercontent.com/nitayStain/goscouter/main/scripts/install.sh | sh

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Clear screen for clean display
clear

# Detect OS and Architecture
OS="$(uname -s)"
ARCH="$(uname -m)"

printf "${BLUE}╔═══════════════════════════════════════════════════════════╗${NC}\n"
printf "${BLUE}║                   GoScouter Installer                     ║${NC}\n"
printf "${BLUE}║          Subdomain Discovery Tool Installation            ║${NC}\n"
printf "${BLUE}╚═══════════════════════════════════════════════════════════╝${NC}\n"
printf "\n"

# Check prerequisites
printf "${YELLOW}→${NC} Checking prerequisites...\n"

command -v git >/dev/null 2>&1 || { printf "${RED}✗ git is required but not installed.${NC}\n" >&2; exit 1; }
command -v go >/dev/null 2>&1 || { printf "${RED}✗ Go 1.23+ is required but not installed. Get it at https://golang.org/dl/${NC}\n" >&2; exit 1; }
command -v node >/dev/null 2>&1 || { printf "${RED}✗ Node.js 18+ is required but not installed. Get it at https://nodejs.org${NC}\n" >&2; exit 1; }
command -v npm >/dev/null 2>&1 || { printf "${RED}✗ npm is required but not installed.${NC}\n" >&2; exit 1; }

printf "${GREEN}✓${NC} All prerequisites found\n"
printf "\n"

# Check if goscouter is already installed
if command -v goscouter >/dev/null 2>&1; then
    CURRENT_VERSION=$(goscouter version 2>&1 | head -1 | grep -o 'v[0-9.]*' || echo "unknown")
    printf "${YELLOW}⚠${NC}  GoScouter is already installed (${CURRENT_VERSION})\n"
    printf "\n"
    printf "Do you want to reinstall? [y/N]: "
    read REPLY
    case "$REPLY" in
        [Yy]|[Yy][Ee][Ss])
            printf "\n"
            ;;
        *)
            printf "Installation cancelled.\n"
            exit 0
            ;;
    esac
fi

# Set install directory
INSTALL_DIR="${HOME}/.goscouter-build"
BRANCH="${GOSCOUTER_BRANCH:-main}"

# Clean up old installation if exists
if [ -d "$INSTALL_DIR" ]; then
    printf "${YELLOW}→${NC} Removing old build directory...\n"
    rm -rf "$INSTALL_DIR"
fi

# Clone repository
printf "${YELLOW}→${NC} Downloading GoScouter from GitHub...\n"
git clone --depth 1 --branch "$BRANCH" https://github.com/nitayStain/goscouter.git "$INSTALL_DIR" >/dev/null 2>&1

cd "$INSTALL_DIR"

# Build frontend
printf "${YELLOW}→${NC} Building frontend (this may take a minute)...\n"
cd frontend
npm install --silent >/dev/null 2>&1
npm run build >/dev/null 2>&1
cd ..

# Build backend
printf "${YELLOW}→${NC} Building backend...\n"
go build -o goscouter ./cmd/goscouter >/dev/null 2>&1

# Install binary
printf "${YELLOW}→${NC} Installing goscouter to /usr/local/bin...\n"
if [ -w "/usr/local/bin" ]; then
    cp goscouter /usr/local/bin/goscouter
else
    sudo cp goscouter /usr/local/bin/goscouter
fi

# Install frontend assets
printf "${YELLOW}→${NC} Installing frontend assets...\n"
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
    printf "\n"
    printf "${GREEN}╔═══════════════════════════════════════════════════════════╗${NC}\n"
    printf "${GREEN}║            ✓ GoScouter installed successfully!            ║${NC}\n"
    printf "${GREEN}╚═══════════════════════════════════════════════════════════╝${NC}\n"
    printf "\n"
    printf "${BLUE}Version:${NC} $VERSION\n"
    printf "${BLUE}Installed to:${NC} /usr/local/bin/goscouter\n"
    printf "\n"
    printf "${YELLOW}Get started:${NC}\n"
    printf "  ${GREEN}goscouter run${NC}       # Start the web service\n"
    printf "  ${GREEN}goscouter version${NC}   # Check version\n"
    printf "  ${GREEN}goscouter help${NC}      # Show help\n"
    printf "\n"
    printf "${BLUE}Web Interface:${NC} http://localhost:8080\n"
    printf "\n"
else
    printf "${RED}✗ Installation failed. Please check errors above.${NC}\n"
    exit 1
fi
