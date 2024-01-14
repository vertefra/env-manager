#!/bin/bash

APP_PATH=./dist/env-manager

# Compile for Linux
GOOS=linux GOARCH=amd64 go build -o "$APP_PATH"
# Compile for macOS
GOOS=darwin GOARCH=amd64 go build -o "$APP_PATH"

