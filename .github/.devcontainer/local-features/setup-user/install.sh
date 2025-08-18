#!/usr/bin/env bash
#-------------------------------------------------------------------------------------------------------------
# Setup user configuration for Bucketeer development environment
#-------------------------------------------------------------------------------------------------------------

USERNAME=${USERNAME:-"codespace"}

set -eux

if [ "$(id -u)" -ne 0 ]; then
  echo -e 'Script must be run as root. Use sudo, su, or add "USER root" to your Dockerfile before running this script.'
  exit 1
fi

# Ensure that login shells get the correct path
rm -f /etc/profile.d/00-restore-env.sh
echo "export PATH=/usr/local/go/bin:/go/bin:/usr/local/bin:/home/${USERNAME}/.local/bin:\$PATH" >/etc/profile.d/00-restore-env.sh
echo "export GOPATH=/go" >>/etc/profile.d/00-restore-env.sh
echo "export GOROOT=/usr/local/go" >>/etc/profile.d/00-restore-env.sh
chmod +x /etc/profile.d/00-restore-env.sh

export DEBIAN_FRONTEND=noninteractive

# Set up user home directory permissions
HOME_DIR="/home/${USERNAME}/"
if [ -d "${HOME_DIR}" ]; then
  chown -R ${USERNAME}:${USERNAME} ${HOME_DIR}
  chmod -R g+r+w "${HOME_DIR}"
  find "${HOME_DIR}" -type d | xargs -n 1 chmod g+s
fi

# Create npm global directory for user
NPM_GLOBAL_DIR="/home/${USERNAME}/.npm-global"
mkdir -p "${NPM_GLOBAL_DIR}"
chown -R ${USERNAME}:${USERNAME} "${NPM_GLOBAL_DIR}"

# Configure sudo PATH
echo "Defaults secure_path=\"/usr/local/go/bin:/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/home/${USERNAME}/.local/bin\"" >>/etc/sudoers.d/$USERNAME

echo "User setup complete!"
