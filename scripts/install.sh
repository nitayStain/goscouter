#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "🚀 GoScouter Installer"
echo "====================="
echo

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed${NC}"
    echo "Please install Go from https://golang.org/dl/"
    exit 1
fi

# Check if Node.js is installed
if ! command -v node &> /dev/null; then
    echo -e "${RED}Error: Node.js is not installed${NC}"
    echo "Please install Node.js from https://nodejs.org/"
    exit 1
fi

# Check if npm is installed
if ! command -v npm &> /dev/null; then
    echo -e "${RED}Error: npm is not installed${NC}"
    echo "Please install npm (comes with Node.js)"
    exit 1
fi

echo "✓ Go version: $(go version | awk '{print $3}')"
echo "✓ Node.js version: $(node --version)"
echo "✓ npm version: $(npm --version)"
echo

# Build frontend
echo "📦 Building frontend..."
cd frontend
npm install
npm run build
cd ..
echo -e "${GREEN}✓ Frontend built successfully${NC}"
echo

# Build Go binary
echo "🔨 Building Go binary..."
go build -o goscouter .
echo -e "${GREEN}✓ Binary built successfully${NC}"
echo

# Determine install location
INSTALL_DIR="/usr/local/bin"
if [ ! -w "$INSTALL_DIR" ]; then
    echo -e "${YELLOW}⚠ Need sudo access to install to $INSTALL_DIR${NC}"
    USE_SUDO="sudo"
else
    USE_SUDO=""
fi

# Install binary
echo "📥 Installing goscouter to $INSTALL_DIR..."
$USE_SUDO cp goscouter "$INSTALL_DIR/goscouter"
$USE_SUDO chmod +x "$INSTALL_DIR/goscouter"
echo -e "${GREEN}✓ Installed successfully${NC}"
echo

# Create data directory for frontend files
DATA_DIR="$HOME/.goscouter"
echo "📁 Setting up data directory at $DATA_DIR..."
mkdir -p "$DATA_DIR"
cp -r frontend/out "$DATA_DIR/"
echo -e "${GREEN}✓ Data directory created${NC}"
echo

# Update the binary to look for frontend in the data directory
echo "🔧 Configuring paths..."
echo "Frontend files: $DATA_DIR/out"
echo

echo -e "${GREEN}✨ Installation complete!${NC}"
echo
echo "You can now run goscouter from anywhere:"
echo "  $ goscouter run"
echo
echo "The web interface will be available at:"
echo "  http://localhost:8080"
echo
