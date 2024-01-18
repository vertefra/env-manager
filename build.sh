#!/bin/bash

set -e

# The name of your Go binary
BIN_NAME="env-manager"
BUILD_DIR="dist"

# Build for Windows
echo "Building for Windows..."
GOOS=windows GOARCH=amd64 go build -o "${BUILD_DIR}/${BIN_NAME}-win-amd64.exe" .
echo "Windows build completed."

# Build for Linux
echo "Building for Linux..."
GOOS=linux GOARCH=amd64 go build -o "${BUILD_DIR}/${BIN_NAME}-linux-amd64" .
echo "Linux build completed."

# Build for macOS
echo "Building for macOS..."
GOOS=darwin GOARCH=amd64 go build -o "${BUILD_DIR}/${BIN_NAME}-darwin-amd64" .
echo "macOS build completed."

# Build for macOS (ARM architecture)
echo "Building for macOS (ARM)..."
GOOS=darwin GOARCH=arm64 go build -o "${BIN_NAME}-darwin-arm64" .
echo "macOS (ARM) build completed."

echo "Build process completed."
