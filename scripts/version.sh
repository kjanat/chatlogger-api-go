#!/usr/bin/env bash
# Semantic version management

VERSION_FILE="internal/version/version.go"
VERSION_REGEX='Version = "([0-9]+)\.([0-9]+)\.([0-9]+)"'

# Get current version
if [[ $(grep -E "$VERSION_REGEX" "$VERSION_FILE") =~ $VERSION_REGEX ]]; then
  MAJOR="${BASH_REMATCH[1]}"
  MINOR="${BASH_REMATCH[2]}"
  PATCH="${BASH_REMATCH[3]}"
  CURRENT="$MAJOR.$MINOR.$PATCH"
else
  echo "Could not extract version information"
  exit 1
fi

echo "Current version: $CURRENT"

# Determine new version based on bump type
case "$1" in
  major)
    NEW_VERSION="$((MAJOR + 1)).0.0"
    ;;
  minor)
    NEW_VERSION="$MAJOR.$((MINOR + 1)).0"
    ;;
  patch)
    NEW_VERSION="$MAJOR.$MINOR.$((PATCH + 1))"
    ;;
  *)
    echo "Usage: $0 [major|minor|patch]"
    exit 1
    ;;
esac

echo "New version: $NEW_VERSION"

# Update version in file
sed -i "s/Version = \"$CURRENT\"/Version = \"$NEW_VERSION\"/" "$VERSION_FILE"
echo "Updated version to $NEW_VERSION in $VERSION_FILE"
