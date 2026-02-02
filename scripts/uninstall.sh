#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "🗑️  GoScouter Uninstaller"
echo "======================="
echo

# Check if goscouter is installed
INSTALL_DIR="/usr/local/bin"
DATA_DIR="$HOME/.goscouter"

if [ ! -f "$INSTALL_DIR/goscouter" ] && [ ! -d "$DATA_DIR" ]; then
    echo -e "${YELLOW}⚠ GoScouter is not installed${NC}"
    echo ""
    echo "Nothing to uninstall."
    exit 0
fi

# Remove binary
if [ -f "$INSTALL_DIR/goscouter" ]; then
    if [ ! -w "$INSTALL_DIR" ]; then
        echo -e "${YELLOW}⚠ Need sudo access to remove from $INSTALL_DIR${NC}"
        USE_SUDO="sudo"
    else
        USE_SUDO=""
    fi

    echo "Removing goscouter binary..."
    $USE_SUDO rm -f "$INSTALL_DIR/goscouter"
    echo -e "${GREEN}✓ Binary removed${NC}"
else
    echo "Binary not found in $INSTALL_DIR"
fi

# Remove data directory
if [ -d "$DATA_DIR" ]; then
    echo "Removing data directory..."
    rm -rf "$DATA_DIR"
    echo -e "${GREEN}✓ Data directory removed${NC}"
else
    echo "Data directory not found"
fi

echo
echo -e "${GREEN}✨ Uninstallation complete!${NC}"
echo
