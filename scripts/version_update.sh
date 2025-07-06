#!/usr/bin/env bash
# version_update.sh - A script to update the version in internal/version/version.go

# Check for the version argument
if [ -z "$1" ]; then
  echo "Usage: ./version_update.sh <new_version>"
  echo "Example: ./version_update.sh 1.2.3"
  exit 1
fi

NEW_VERSION=$1

# Validate version format (Semantic Versioning: major.minor.patch)
if ! [[ $NEW_VERSION =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "Error: Version must follow semantic versioning format (e.g., 1.2.3)"
  exit 1
fi

# Get directory of the script
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ROOT_DIR="$DIR/.."

# Update version.go
echo "Updating version in internal/version/version.go..."
sed -i "s/Version = \"[0-9]*\.[0-9]*\.[0-9]*\"/Version = \"$NEW_VERSION\"/" "$ROOT_DIR/internal/version/version.go"

echo "Version updated to $NEW_VERSION"
echo "Remember to commit these changes and push them to your repository"
echo "You can also run the tag-version GitHub Action workflow to create a git tag"
