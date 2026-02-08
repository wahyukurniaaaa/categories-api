#!/bin/bash

# Script untuk build Docker image untuk category-api-golang
# Usage: ./build-docker.sh [tag]

set -e

# Warna untuk output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Default values
IMAGE_NAME="category-api-golang"
TAG="${1:-latest}"
REGISTRY="${DOCKER_REGISTRY:-ghcr.io}"
USERNAME="${DOCKER_USERNAME:-wahyukurniaaaa}"

FULL_IMAGE_NAME="${REGISTRY}/${USERNAME}/${IMAGE_NAME}:${TAG}"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Building Docker Image${NC}"
echo -e "${BLUE}========================================${NC}"
echo -e "Image: ${GREEN}${FULL_IMAGE_NAME}${NC}"
echo ""

# Build image
echo -e "${BLUE}Step 1: Building Docker image...${NC}"
docker build -t "${FULL_IMAGE_NAME}" .

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Build successful!${NC}"
else
    echo -e "${RED}✗ Build failed!${NC}"
    exit 1
fi

# Tag as latest if not already
if [ "${TAG}" != "latest" ]; then
    echo -e "${BLUE}Step 2: Tagging as latest...${NC}"
    docker tag "${FULL_IMAGE_NAME}" "${REGISTRY}/${USERNAME}/${IMAGE_NAME}:latest"
fi

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Build Complete!${NC}"
echo -e "${GREEN}========================================${NC}"
echo -e "Image: ${FULL_IMAGE_NAME}"
echo ""
echo -e "${BLUE}To run the container:${NC}"
echo "  docker run -p 8080:8080 --env-file .env ${FULL_IMAGE_NAME}"
echo ""
echo -e "${BLUE}To push to registry:${NC}"
echo "  docker push ${FULL_IMAGE_NAME}"
echo ""
