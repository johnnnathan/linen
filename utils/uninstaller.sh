#!/usr/bin/env bash 

BIN_NAME="linen"
INSTALL_DIR="/usr/local/bin"


echo "Uninstalling linen..."

if command -v "$BIN_NAME" &> /dev/null; then
  sudo rm -f "$INSTALL_DIR/$BIN_NAME"
  if [ $? -eq 0 ]; then
    echo "Uninstallation successful"
    exit 0
  else
    echo "Uninstallation unsuccessful, remove command returned non-zero value"
    exit 2
  fi
else
  echo "Binary does not exist"
  exit 1
fi
