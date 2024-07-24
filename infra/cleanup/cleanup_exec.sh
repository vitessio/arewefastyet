#!/bin/sh

# Path to the exec directory
EXEC_DIR="/exec"

# Ensure the directory exists and is not empty
if [ -d "$EXEC_DIR" ] && [ "$(ls -A $EXEC_DIR)" ]; then
  # Find and delete directories that are 7-days old or older
  find "$EXEC_DIR" -mindepth 1 -type d -mtime +6 -exec rm -rf {} +
fi
