#!/bin/bash

# Deploy follooow-be API
# This script is intended to be placed in /root and executed with appropriate privileges.

# Configuration
TARGET_DIR="/root/apps/api.follooow.com"
ENV_FILE="$TARGET_DIR/.env"
APP_NAME="follooow-be"

# Ensure target directory is clean and clone repository
if [ -d "$TARGET_DIR" ]; then
  echo "Removing existing directory $TARGET_DIR..."
  rm -rf "$TARGET_DIR" || { echo "Failed to remove $TARGET_DIR"; exit 1; }
fi
echo "Cloning repository into $TARGET_DIR..."
git clone --depth 1 --single-branch git@github.com:yna-studio/follooow-be-go.git "$TARGET_DIR" || { echo "Git clone failed"; exit 1; }

cd "$TARGET_DIR" || { echo "Cannot cd to $TARGET_DIR"; exit 1; }

# Stop existing pm2 process if running
if pm2 describe "$APP_NAME" > /dev/null 2>&1; then
  echo "Stopping existing pm2 process $APP_NAME..."
  pm2 stop "$APP_NAME"
else
  echo "pm2 process $APP_NAME not found; skipping stop."
fi

# Write .env file
cat > "$ENV_FILE" <<EOF
TELEGRAM_FOLLOOOW_CHANNEL=follooow_channel
TELEGRAM_FOLLOOOW_TOKEN=6300896180:AAHxybkFLqQQ-Uc7oPBuKIdUvvhsmlhrm0o
MONGO_URI=mongodb://127.0.0.1:27017/follooow
MONGO_DB=follooow
EOF

# Build the Go binary
echo "Building the Go application..."
if go build -o "$APP_NAME" main.go; then
  echo "Build succeeded."
else
  echo "Go build failed. Exiting."
  exit 1
fi

# Start the application with pm2
echo "Starting pm2 process $APP_NAME..."
pm2 start "$APP_NAME" || pm2 start "$APP_NAME"

echo "Deployment completed."
