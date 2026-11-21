#!/bin/bash

set -e

# The machine to wake: WOL_MAC=00:11:22:33:44:55 ./setup_wolserver.sh
: "${WOL_MAC:?set WOL_MAC to the MAC address of the machine to wake}"

APP_NAME="wolserver"
INSTALL_DIR="/opt/$APP_NAME"
SERVICE_FILE="/etc/systemd/system/${APP_NAME}.service"
USER="pi"
GROUP="pi"

echo "📦 Building Go binary..."
GOOS=linux GOARCH=$(uname -m | grep -q aarch64 && echo arm64 || echo arm) go build -o "$APP_NAME" main.go

echo "📁 Creating install directory at $INSTALL_DIR..."
sudo mkdir -p "$INSTALL_DIR"

echo "📂 Copying files..."
sudo cp "$APP_NAME" "$INSTALL_DIR/"
sudo cp -r static "$INSTALL_DIR/"

echo "📝 Creating systemd service file..."

cat <<EOF | sudo tee "$SERVICE_FILE" > /dev/null
[Unit]
Description=Wake-on-LAN Web Server
After=network.target

[Service]
Environment=WOL_MAC=${WOL_MAC}
ExecStart=${INSTALL_DIR}/${APP_NAME}
WorkingDirectory=${INSTALL_DIR}
Restart=on-failure
User=${USER}
Group=${GROUP}

[Install]
WantedBy=multi-user.target
EOF

echo "🔄 Reloading systemd and enabling service..."
sudo systemctl daemon-reexec
sudo systemctl daemon-reload
sudo systemctl enable "$APP_NAME"
sudo systemctl restart "$APP_NAME"

echo "✅ Done! Service status:"
sudo systemctl status "$APP_NAME" --no-pager

