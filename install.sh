#!/bin/sh
# Crevisto CLI installer — one-liner: curl -fsSL https://crevisto.com/install.sh | sh
set -e

VERSION="latest"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="crevisto"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$OS" in
  linux*)  OS="linux" ;;
  darwin*) OS="darwin" ;;
  *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

case "$ARCH" in
  x86_64|amd64)  ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Determine download URL
if [ "$VERSION" = "latest" ]; then
  URL="https://github.com/javimosch/crevisto/releases/latest/download/crevisto-${OS}-${ARCH}"
else
  URL="https://github.com/javimosch/crevisto/releases/download/${VERSION}/crevisto-${OS}-${ARCH}"
fi

echo "Downloading Crevisto CLI for ${OS}/${ARCH}..."
echo "  $URL"

# Download to temp file
TMP_FILE=$(mktemp)
curl -fsSL "$URL" -o "$TMP_FILE"
chmod +x "$TMP_FILE"

# Move to install directory
if [ -w "$INSTALL_DIR" ]; then
  mv "$TMP_FILE" "$INSTALL_DIR/$BINARY_NAME"
else
  echo "sudo required to install to $INSTALL_DIR"
  sudo mv "$TMP_FILE" "$INSTALL_DIR/$BINARY_NAME"
fi

echo ""
echo "✓ Crevisto CLI installed to $INSTALL_DIR/$BINARY_NAME"
echo ""
echo "Quick start:"
echo "  crevisto trial          # Get 5 free credits"
echo "  crevisto tools          # List all 40 AI tools"
echo "  crevisto generate <slug> --input photo=./img.jpg"
echo "  crevisto whoami         # Check your balance"
echo ""
echo "Docs: https://github.com/javimosch/crevisto"
