#!/bin/bash
set -eu

# ====================================================================================
# VARIABLES
# ====================================================================================
TIMEZONE=Asia/Taipei
USERNAME=joker

# Prompt to enter a password for the PostgreSQL joker user.
read -p "Enter password for joker DB user: " DB_PASSWORD

# Force all output to use English locale to avoid locale errors
export LC_ALL=en_US.UTF-8

# ====================================================================================
# SCRIPT LOGIC
# ====================================================================================

# Enable universe repo & update package list
add-apt-repository --yes universe
apt update

# Set timezone and install all locales
timedatectl set-timezone ${TIMEZONE}
apt --yes install locales-all

# Add new user and give sudo access
useradd --create-home --shell "/bin/bash" --groups sudo "${USERNAME}"
passwd --delete "${USERNAME}"
chage --lastday 0 "${USERNAME}"

# Copy root SSH keys to new user
rsync --archive --chown=${USERNAME}:${USERNAME} /root/.ssh /home/${USERNAME}

# Enable UFW (firewall) for SSH, HTTP, HTTPS
ufw allow 22
ufw allow 80/tcp
ufw allow 443/tcp
ufw --force enable

# Install fail2ban
apt --yes install fail2ban

# Install Goose (https://github.com/pressly/goose)
curl -fsSL \
    https://raw.githubusercontent.com/pressly/goose/master/install.sh |\
    sh


# Install PostgreSQL
apt --yes install postgresql

# Create PostgreSQL database and user
sudo -i -u postgres psql -c "CREATE DATABASE joker"
sudo -i -u postgres psql -d joker -c "CREATE ROLE joker WITH LOGIN PASSWORD '${DB_PASSWORD}'"

# Add system-wide environment variable for DSN
echo "JOKER_DB_DSN='postgres://joker:${DB_PASSWORD}@localhost/joker?sslmode=disable'" >> /etc/environment

# Install Caddy for reverse proxy + HTTPS
apt install -y debian-keyring debian-archive-keyring apt-transport-https
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
apt update
apt --yes install caddy

# Final upgrade
apt --yes -o Dpkg::Options::="--force-confnew" upgrade

echo "Setup complete! Rebooting..."
reboot
