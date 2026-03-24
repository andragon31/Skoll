#!/bin/bash
set -e

REPO="andragon31/Skoll"
INSTALL_DIR="/usr/local/bin"

if [[ "$OSTYPE" == "darwin"* ]]; then
    ARCH=$(uname -m)
    if [ "$ARCH" = "arm64" ]; then
        BIN="skoll-darwin-arm64"
    else
        BIN="skoll-darwin-amd64"
    fi
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    ARCH=$(uname -m)
    if [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
        BIN="skoll-linux-arm64"
    else
        BIN="skoll-linux-amd64"
    fi
else
    echo "Unsupported OS: $OSTYPE"
    exit 1
fi

TMP=$(mktemp)
URL="https://github.com/${REPO}/releases/latest/download/${BIN}"

echo "Downloading Skoll..."
curl -fsSL "$URL" -o "$TMP"
chmod +x "$TMP"

if [ -w "$INSTALL_DIR" ]; then
    mv "$TMP" "$INSTALL_DIR/skoll"
    echo "Installed to $INSTALL_DIR/skoll"
else
    echo "Installing to $INSTALL_DIR requires sudo..."
    sudo mv "$TMP" "$INSTALL_DIR/skoll"
    echo "Installed to $INSTALL_DIR/skoll"
fi

echo ""
echo "Skoll installed! Run:"
echo "  skoll version          # Verify"
echo "  skoll init            # Initialize"
echo ""
