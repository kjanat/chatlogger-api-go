#!/bin/bash
set -e

# Stop and disable the service if systemd is available
if command -v systemctl >/dev/null; then
    echo "Stopping chatlogger-api service..."
    systemctl stop chatlogger-api.service || true
    systemctl disable chatlogger-api.service || true
    systemctl daemon-reload
fi

# Don't remove user/group or data directories during package removal
# This allows for reinstallation without data loss
echo "Note: User data in /var/lib/chatlogger is preserved for potential reinstallation."
echo "To completely remove all data, manually run: rm -rf /var/lib/chatlogger"

echo "ChatLogger API removal completed!"
exit 0
