#!/bin/bash

set -e

IMAGE_NAME="surge-protocol"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "Building Docker image: $IMAGE_NAME"
echo "================================"

PLATFORM="linux/amd64"
if [[ "$(uname -m)" == "arm64" || "$(uname -m)" == "aarch64" ]]; then
    PLATFORM="linux/arm64"
fi
echo "Platform: $PLATFORM"

BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
TAG_VERSION=$(git describe --tags --exact-match 2>/dev/null || echo "")

SANITIZED_BRANCH=$(echo "$BRANCH" | sed 's/[^a-zA-Z0-9._-]/-/g')

TAGS="--tag $IMAGE_NAME:latest"
TAGS="$TAGS --tag $IMAGE_NAME:$SANITIZED_BRANCH"
TAGS="$TAGS --tag $IMAGE_NAME:sha-$COMMIT"

if [[ -n "$TAG_VERSION" ]]; then
    TAGS="$TAGS --tag $IMAGE_NAME:$TAG_VERSION"
fi

echo "Tags to be created:"
echo "$TAGS" | tr ' ' '\n' | grep -E "^--tag" | sed 's/--tag /  - /'
echo ""

cd "$SCRIPT_DIR"
docker build \
    --platform "$PLATFORM" \
    --file Dockerfile \
    $TAGS \
    .

echo "✅ Docker build completed successfully!"

