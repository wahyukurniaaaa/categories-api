#!/bin/bash

# Script untuk push Docker image ke registry
# Usage: ./push-docker.sh [tag]

set -e

# Warna untuk output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default values
IMAGE_NAME="category-api-golang"
TAG="${1:-latest}"
REGISTRY="${DOCKER_REGISTRY:-ghcr.io}"
USERNAME="${DOCKER_USERNAME:-wahyukurniaaaa}"

FULL_IMAGE_NAME="${REGISTRY}/${USERNAME}/${IMAGE_NAME}:${TAG}"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Pushing Docker Image${NC}"
echo -e "${BLUE}========================================${NC}"
echo -e "Image: ${GREEN}${FULL_IMAGE_NAME}${NC}"
echo ""

# Check if logged in
echo -e "${BLUE}Checking Docker registry login...${NC}"
if ! docker info | grep -q "Username"; then
    echo -e "${YELLOW}Not logged in to Docker registry.${NC}"
    echo -e "${BLUE}Logging in to ${REGISTRY}...${NC}"
    echo -e "${YELLOW}You may need to use a Personal Access Token (PAT)${NC}"
    docker login ${REGISTRY}
fi

# Push image
echo -e "${BLUE}Pushing image to registry...${NC}"
docker push "${FULL_IMAGE_NAME}"

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Push successful!${NC}"
else
    echo -e "${RED}✗ Push failed!${NC}"
    exit 1
fi

# Push latest tag if not already
if [ "${TAG}" != "latest" ]; then
    echo -e "${BLUE}Pushing latest tag...${NC}"
    docker push "${REGISTRY}/${USERNAME}/${IMAGE_NAME}:latest"
fi

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Push Complete!${NC}"
echo -e "${GREEN}========================================${NC}"
echo -e "Image available at: ${FULL_IMAGE_NAME}"
echo ""
echo -e "${BLUE}To pull this image:${NC}"
echo "  docker pull ${FULL_IMAGE_NAME}"
echo ""
