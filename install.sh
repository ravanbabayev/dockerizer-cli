#!/bin/bash

# Detect OS and architecture
OS="$(uname -s)"
ARCH="$(uname -m)"

# Convert architecture
case "$ARCH" in
    x86_64)
        ARCH="amd64"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

# Set binary name based on OS
case "$OS" in
    Linux)
        BINARY="dockerizer-linux-${ARCH}"
        ;;
    Darwin)
        BINARY="dockerizer-darwin-${ARCH}"
        ;;
    *)
        echo "Unsupported operating system: $OS"
        exit 1
        ;;
esac

# Create installation directory
INSTALL_DIR="/usr/local/bin"
SUPPORT_DIR="/usr/local/share/dockerizer"

# Download latest release
echo "Downloading latest release..."
LATEST_RELEASE=$(curl -s https://api.github.com/repos/ravanbabayev/dockerizer-cli/releases/latest | grep "tag_name" | cut -d '"' -f 4)
DOWNLOAD_URL="https://github.com/ravanbabayev/dockerizer-cli/releases/download/${LATEST_RELEASE}/${BINARY}.tar.gz"

# Download and extract
curl -L "$DOWNLOAD_URL" | tar xz

# Create directories
sudo mkdir -p "$SUPPORT_DIR"
sudo mv supported "$SUPPORT_DIR/"

# Install binary
sudo mv "$BINARY" "$INSTALL_DIR/dockerizer"
sudo chmod +x "$INSTALL_DIR/dockerizer"

echo "Installation complete! You can now use 'dockerizer' command." 