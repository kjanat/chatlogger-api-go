#!/usr/bin/env bash
go build \
    -ldflags \
    -X "chatlogger-api-go/internal/version.Version=$(grep -oP 'Version = "\K[^"]+' internal/version/version.go)" \
    -X "chatlogger-api-go/internal/version.BuildTime=$(date -u +"%Y-%m-%dT%H:%M:%SZ")" \
    -X "chatlogger-api-go/internal/version.GitCommit=$(git rev-parse HEAD)" \
    -o chatlogger-api \
    ./cmd/server
