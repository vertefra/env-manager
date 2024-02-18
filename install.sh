#!/bin/bash

# Set the GitHub repository
REPO="vertefra/env-manager"

# Get the OS and architecture for the system
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [ "$ARCH" == "x86_64" ]; then
  ARCH="amd64"
fi

# Set the binary name
BINARY_NAME="env-manager-$OS-$ARCH"

# GitHub's API URL for the latest release
API_URL="https://api.github.com/repos/$REPO/releases/latest"

# Use curl to fetch the latest release data and extract the asset download URL for your binary
DOWNLOAD_URL=$(curl -s $API_URL | grep "browser_download_url.*$BINARY_NAME" | cut -d '"' -f 4)

if [ -z "$DOWNLOAD_URL" ]; then
  echo "Error: Unable to find the download URL for $BINARY_NAME in the latest release."
  exit 1
fi

# Download the binary
curl -L -o $BINARY_NAME $DOWNLOAD_URL

# Make the binary executable
chmod +x $BINARY_NAME

echo "Installation completed. The binary has been saved as ./$BINARY_NAME"
