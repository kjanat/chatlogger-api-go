#!/usr/bin/env bash
set -euo pipefail

# scripts/docs_generate.sh
# This script generates Swagger/OpenAPI documentation from Go annotations using swag

# Parse command line arguments
v31=false
outPath="./docs"

while [ "$#" -gt 0 ]; do
  case "$1" in
    --v31)
      v31=true
      shift 1
      ;;
    # Removed redundant --out flag handling
    --outPath)
      outPath="$2"
      shift 2
      ;;
    *)
      echo "Unknown option: $1"
    echo "Usage: $0 [--v31] [--outPath <path>]"
    # Removed --out flag from usage instructions
      exit 1
      ;;
  esac
done

# Function to prepare environment
prep_env() {
    local outPath="$1"
    local filePath="$2"

    # Ensure we're starting with fresh docs
    if [ -d "$outPath" ]; then
        echo "Cleaning previous documentation..."
        find "$outPath" -type f -name "${filePath}_*" -delete
        # Removed redundant rm command as find already handles file removal
    fi
    # Ensure the output directory exists
    mkdir -p "$outPath"

    # Display current swag version
    swagVersion=$(swag --version)
    echo "Swag version: $swagVersion"
}

# Function to check run result
test_run() {
    local outPath="$1"
    local filePath="$2"
    local exitCode="$3"

    if [ $exitCode -ne 0 ]; then
        echo "Documentation generation failed with code $exitCode"
        exit $exitCode
    fi

    # Verify the documentation files were created
    if ls "${outPath}/${filePath}_"* 1> /dev/null 2>&1; then
        echo "✅ Documentation generated ${filePath} successfully in ${outPath}/"
    else
        echo "❌ Failed to generate documentation files"
        exit 1
    fi
}

# Change directory to project root to ensure proper parsing
cd "$(dirname "${BASH_SOURCE[0]}")/.." || { echo "Failed to change directory to script root"; exit 1; }

# Generate documentation based on flags
if [ "$v31" = true ]; then
    echo "Generating OpenAPI v3 documentation..."
    envPath="OpenAPIv3"
    prep_env "$outPath" "$envPath"
    swag init \
        --output "$outPath" \
        --dir "cmd/server,internal/api,internal/config,internal/domain,internal/handler,internal/hash,internal/jobs,internal/middleware,internal/repository,internal/service,internal/strategy,internal/version" \
        --exclude "$outPath" \
        --instanceName "$envPath" --v3.1 --packageName "docs_v3" \
        --parseInternal --parseFuncBody --requiredByDefault \
        --outputTypes "yaml,json"
    test_run "$outPath" "$envPath" $?
else
    echo "Generating OpenAPI v2 documentation..."
    envPath="OpenAPI"
    prep_env "$outPath" "$envPath"
    swag init \
        --output "$outPath" \
        --dir "cmd/server,internal/api,internal/config,internal/domain,internal/handler,internal/hash,internal/jobs,internal/middleware,internal/repository,internal/service,internal/strategy,internal/version" \
        --exclude "$outPath" \
        --instanceName "$envPath" --packageName "docs" \
        --parseInternal --parseFuncBody --requiredByDefault
    test_run "$outPath" "$envPath" $?
fi

if [ "$v31" = true ]; then
    echo "API documentation generation complete for OpenAPI v3.1."
else
    echo "API documentation generation complete for OpenAPI v2."
fi
