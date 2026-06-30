#!/bin/bash
set -e

# Create chatlogger user and group if they don't exist
if ! getent group chatlogger >/dev/null; then
    groupadd --system chatlogger
fi

if ! getent passwd chatlogger >/dev/null; then
    useradd --system --gid chatlogger --no-create-home --shell /bin/false chatlogger
fi

# Create necessary directories with proper permissions
mkdir -p /etc/chatlogger
mkdir -p /var/lib/chatlogger/data
mkdir -p /opt/chatlogger

# Set proper permissions
chown -R chatlogger:chatlogger /var/lib/chatlogger
chmod -R 750 /var/lib/chatlogger

# Create default config if it doesn't exist
if [ ! -f /etc/chatlogger/config.yaml ]; then
    cat > /etc/chatlogger/config.yaml << EOF
database:
  driver: postgres
  host: localhost
  port: 5432
  name: chatlogger
  user: chatlogger
  password: change_me_please
  sslmode: disable

server:
  port: 8080
  host: 0.0.0.0
  timeout: 30s

logging:
  level: info
  format: json

auth:
  jwt_secret: change_this_to_a_secure_random_string
  token_expiry: 24h
EOF

    # Secure the config file
    chown root:chatlogger /etc/chatlogger/config.yaml
    chmod 640 /etc/chatlogger/config.yaml

    echo "Created default configuration in /etc/chatlogger/config.yaml"
    echo "IMPORTANT: Please update database credentials and JWT secret before starting the service!"
fi

# Enable and start the service
if command -v systemctl >/dev/null; then
    echo "Enabling chatlogger-api service..."
    systemctl daemon-reload
    systemctl enable chatlogger-api.service
    echo "You can now start the service with: systemctl start chatlogger-api.service"
fi

echo "ChatLogger API installation completed!"
exit 0
