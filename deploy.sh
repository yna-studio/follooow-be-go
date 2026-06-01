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
TELEGRAM_FOLLOOOW_CHANNEL=@follooow_channel
TELEGRAM_FOLLOOOW_TOKEN=6300896180:AAHxybkFLqQQ-Uc7oPBuKIdUvvhsmlhrm0o
MONGO_URI=mongodb+srv://webdev:Rahasia20@follooow-prod.8ewjbku.mongodb.net/?appName=follooow-prod
MONGO_DB=follooow
CLOUDINARY_DIR=/follooow
CLOUDINARY_CLOUD_NAME=dhjkktmal
CLOUDINARY_API_KEY=546653438788785
CLOUDINARY_API_SECRET=pAtLP1NVgyxcSKzG68eCH-RcbWw
EOF

# Build the Go binary
echo "Building the Go application..."
go build main.go

# Start the application with pm2
echo "Starting pm2 process $APP_NAME..."
pm2 start "$APP_NAME" || pm2 start "$APP_NAME"

echo "Deployment completed."
