#!/usr/bin/env bash

BIN_NAME="linen"
BIN_PATH="./$BIN_NAME"
INSTALL_DIR="/usr/local/bin"

if [[ ! -f "$BIN_PATH" ]]; then 
  echo "ERROR: Binary file not found"
  exit 1
fi
  

echo "Installing linen ..."
sudo cp "$BIN_PATH" "$INSTALL_DIR/$BIN_NAME"

sudo chmod +x "$INSTALL_DIR/$BIN_NAME"

if command -v "$BIN_NAME" &> /dev/null; then 
  echo "Installation Successful"
  exit 0
else
  echo "Installation Unsuccessful, linen command not found"
  exit 2
fi 
