#!/bin/bash
set -euo pipefail

echo "Starting post-create setup..."

# Keep Codespaces-mounted files usable by tools invoked from make docker-compose-up.
# In particular, MySQL ignores world-writable config files and then misses local settings.
run_privileged() {
  if [ "$(id -u)" -eq 0 ]; then
    "$@"
  elif sudo -n true >/dev/null 2>&1; then
    sudo "$@"
  else
    echo "sudo is unavailable; skipping privileged command: $*"
  fi
}

USER_NAME="${USER_NAME:-codespace}"
if id "$USER_NAME" >/dev/null 2>&1; then
  run_privileged mkdir -p /go/pkg "/home/$USER_NAME"
  run_privileged chown -R "$USER_NAME:$USER_NAME" /go/pkg "/home/$USER_NAME"
else
  echo "User '$USER_NAME' not found; skipping ownership setup"
fi

MYSQL_CONF_DIR="${WORKSPACE_DIR:-/workspace/bucketeer}/docker-compose/config/mysql-conf"
if [ -d "$MYSQL_CONF_DIR" ]; then
  run_privileged chown -R "$USER_NAME:$USER_NAME" "$MYSQL_CONF_DIR"
  # MySQL ignores world-writable config files, which can happen with Codespaces workspace mounts.
  find "$MYSQL_CONF_DIR" -type f -name "*.cnf" -exec chmod 0644 {} \;
  find "$MYSQL_CONF_DIR" -type d -exec chmod 0755 {} \;
else
  echo "MySQL config directory not found: $MYSQL_CONF_DIR"
fi

echo "Post-create setup completed"
